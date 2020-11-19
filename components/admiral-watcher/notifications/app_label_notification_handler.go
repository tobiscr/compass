package notifications

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/kyma-incubator/compass/components/admiral-watcher/types"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"strings"
)

const runtimeName = "runtime-poc"
const commerceMockName = "commerce-mock"

type AppLabelNotificationHandler struct {
	RuntimeLister      RuntimeLister
	AppLister          ApplicationLister
	AppLabelGetter     ApplicationLabelGetter
	RuntimeLabelGetter RuntimeLabelGetter
	Transact           persistence.Transactioner
	ScriptRunner       script.Runner
}

func (a *AppLabelNotificationHandler) HandleCreate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) HandleUpdate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) HandleDelete(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) handle(ctx context.Context, label Label) error {
	if label.Key != model.ScenariosKey {
		log.C(ctx).Infof("label %v is not scenarios", label)
		return nil
	}

	if len(label.AppID) == 0 {
		log.C(ctx).Infof("label %v is not for apps", label)
		return nil
	}

	tx, err := a.Transact.Begin()
	if err != nil {
		return err
	}
	defer a.Transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)
	ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	query := `$[*] ? ( `
	queryEnd := ` )`
	queries := make([]string, 0, len(label.Value))
	for _, val := range label.Value {
		queries = append(queries, fmt.Sprintf("@ == \"%s\"", val))
	}
	query = query + strings.Join(queries, "||") + queryEnd
	runtimesList, err := a.RuntimeLister.List(ctx, []*labelfilter.LabelFilter{
		labelfilter.NewForKeyWithQuery(model.ScenariosKey, query),
	}, 100, "")
	if err != nil {
		return err
	}
	for _, runtime := range runtimesList.Data {
		if runtime.Name != runtimeName {
			log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", runtimeName, runtime.Name)
			continue
		}

		scenarioLabel, err := a.RuntimeLabelGetter.GetLabel(ctx, runtime.ID, "scenarios")
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Warnf("runtime with id %s does not have scenarios label, skipping", runtime.ID)
				continue
			}
			return err
		}
		scenarioLabelSlice := scenarioLabel.Value.([]interface{})
		if len(scenarioLabelSlice) == 1 && scenarioLabelSlice[0] == "DEFAULT" {
			log.C(ctx).Warnf("app with id %s is only in the DEFAULT scenario, skipping", runtime.ID)
			continue
		}

		parsedID, err := uuid.Parse(runtime.ID)
		if err != nil {
			return err
		}

		appsList, err := a.AppLister.ListByRuntimeID(ctx, parsedID, 100, "")
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Warnf("app with id %s not found during handling of label event", label.AppID)
				err = tx.Commit()
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}

		appNames := make([]string, 0, appsList.TotalCount)
		for _, app := range appsList.Data {
			if app.Status.Condition != model.ApplicationStatusConditionConnected {
				log.C(ctx).Infof("app %s is not connected but is in status %s", app.Name, app.Status.Condition)
				continue
			}
			scenarioLabel, err := a.AppLabelGetter.GetLabel(ctx, app.ID, "scenarios")
			if err != nil {
				if apperrors.IsNotFoundError(err) {
					log.C(ctx).Warnf("app with id %s does not have scenarios label, skipping", label.AppID)
					continue
				}
				return err
			}
			scenarioLabelSlice := scenarioLabel.Value.([]interface{})
			if len(scenarioLabelSlice) == 1 && scenarioLabelSlice[0] == "DEFAULT" {
				log.C(ctx).Warnf("app with id %s is only in the DEFAULT scenario, skipping", label.AppID)
				continue
			}

			appNames = append(appNames, app.Name)
		}

		log.C(ctx).Infof("Number of applications in scenario with test runtime: %d", len(appNames))
		dependencyName := "dep-rt-" + runtime.ID
		if len(appNames) == 0 {
			if err := a.ScriptRunner.DeleteDependency(ctx, dependencyName, "admiral.yaml", "runtime.yaml"); err != nil {
				return err
			}
		} else {
			exists, err := a.ScriptRunner.DependencyExists(ctx, dependencyName)
			if err != nil {
				log.C(ctx).Errorf("unable to determine existance of %q dependency: %s", dependencyName, err)
				return err
			}

			shouldCreateDep := stringsAnyEquals(appNames, commerceMockName)

			if shouldCreateDep && !exists {
				dep := types.Dependency{
					TypeMeta: types.TypeMeta{
						Kind:       "Dependency",
						APIVersion: "admiral.io/v1alpha1",
					},
					ObjectMeta: types.ObjectMeta{
						Name:      dependencyName,
						Namespace: "admiral",
					},
					Spec: types.MDependency{
						//Source:        "webapp-rt-" + runtime.ID,
						Source:        "webapp",
						IdentityLabel: "identity",
						Destinations:  []string{commerceMockName},
					},
				}

				if err := a.ScriptRunner.ApplyDependency(ctx, dep, "admiral.yaml"); err != nil {
					return err
				}
			} else if !shouldCreateDep && exists {
				if err := a.ScriptRunner.DeleteDependency(ctx, dependencyName, "admiral.yaml", "runtime.yaml"); err != nil {
					return err
				}
			}
		}

		if err := syncServiceEntries(ctx, a.ScriptRunner, appNames); err != nil {
			log.C(ctx).Errorf("unable to sync service entries for applications as part of application label event: %s", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func syncServiceEntries(ctx context.Context, scriptRunner script.Runner, expectedAppNames []string) error {
	existingAppNames, err := scriptRunner.GetExistingServices(ctx)
	if err != nil {
		return err
	}

	appsToDelete := make([]string, 0)
	for _, app := range existingAppNames {
		if stringsAnyEquals(expectedAppNames, app) {
			continue
		} else {
			appsToDelete = append(appsToDelete, app)
		}
	}

	appsToCreate := make([]string, 0)
	for _, app := range expectedAppNames {
		if stringsAnyEquals(existingAppNames, app) {
			continue
		} else {
			appsToCreate = append(appsToCreate, app)
		}
	}

	// delete all resources for systems removed from scenario
	for _, appName := range appsToDelete {
		if appName == commerceMockName {
			log.C(ctx).Infof("service resources won't be edited for the %q application", appName)
			continue
		}

		if err := scriptRunner.DeleteResource(ctx, fmt.Sprintf("service-entries/%s.yaml", appName)); err != nil {
			return err
		}
	}

	// apply resources for systems added to scenario
	for _, appName := range appsToCreate {
		if appName == commerceMockName {
			log.C(ctx).Infof("service resources won't be edited for the %q application", appName)
			continue
		}

		if err := scriptRunner.ApplyResource(ctx, fmt.Sprintf("service-entries/%s.yaml", appName)); err != nil {
			return err
		}
	}

	return nil
}

func cleanupServiceEntries(ctx context.Context, scriptRunner script.Runner) error {
	return scriptRunner.DeleteResource(ctx, "service-entries/")
}

// stringsAnyEquals returns true if any of the strings in the slice equal the given string.
func stringsAnyEquals(stringSlice []string, str string) bool {
	for _, v := range stringSlice {
		if v == str {
			return true
		}
	}
	return false
}