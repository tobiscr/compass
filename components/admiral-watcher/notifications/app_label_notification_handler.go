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
		if runtime.Name != "runtime-poc" {
			log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", "runtime-poc", runtime.Name)
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

		if len(appNames) == 0 {
			if err := a.ScriptRunner.DeleteDependency(ctx, "dep-rt-"+runtime.ID, "admiral.yaml"); err != nil {
				return err
			}
		} else {
			dep := types.Dependency{
				TypeMeta: types.TypeMeta{
					Kind:       "Dependency",
					APIVersion: "admiral.io/v1alpha1",
				},
				ObjectMeta: types.ObjectMeta{
					Name:      "dep-rt-" + runtime.ID,
					Namespace: "admiral",
				},
				Spec: types.MDependency{
					//Source:        "webapp-rt-" + runtime.ID,
					Source:        "webapp",
					IdentityLabel: "identity",
					Destinations:  appNames,
				},
			}
			if err := a.ScriptRunner.ApplyDependency(ctx, dep, "admiral.yaml"); err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
