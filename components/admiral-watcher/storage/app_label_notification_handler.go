package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/kyma-incubator/compass/components/admiral-watcher/types"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
)

type AppLabelNotificationHandler struct {
	RuntimeLister RuntimeLister
	AppLister     ApplicationLister
	scriptRunner  script.Runner
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
		log.C(ctx).Infof("label %v is not apps", label)
		return nil
	}

	ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	runtimesList, err := a.RuntimeLister.List(ctx, []*labelfilter.LabelFilter{
		labelfilter.NewForKeyWithQuery(model.ScenariosKey, label.Value),
	}, 100, "")
	if err != nil {
		return err
	}
	for _, runtime := range runtimesList.Data {
		parsedID, err := uuid.Parse(runtime.ID)
		if err != nil {
			return err
		}

		appsList, err := a.AppLister.ListByRuntimeID(ctx, parsedID, 100, "")
		if err != nil {
			return err
		}

		appNames := make([]string, 0, appsList.TotalCount)
		for _, app := range appsList.Data {
			appNames = append(appNames, app.Name)
		}

		dep := types.Dependency{
			TypeMeta: types.TypeMeta{
				Kind:       "Dependency",
				APIVersion: "admiral.io/v1alpha1",
			},
			ObjectMeta: types.ObjectMeta{
				Name:      "dep-" + runtime.Name,
				Namespace: "admiral",
			},
			Spec: types.MDependency{
				Source:        "webapp-" + runtime.Name,
				IdentityLabel: "identity",
				Destinations:  appNames,
			},
		}
		if err := a.scriptRunner.ApplyDependency(ctx, dep, "runtime.yaml"); err != nil {
			return err
		}
	}

	return nil
}
