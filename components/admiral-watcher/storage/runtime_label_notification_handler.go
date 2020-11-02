package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/kyma-incubator/compass/components/admiral-watcher/types"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
)

type RuntimeLabelNotificationHandler struct {
	RuntimeGetter RuntimeGetter
	AppLister     ApplicationLister
	scriptRunner  script.Runner
}

func (a *RuntimeLabelNotificationHandler) HandleCreate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *RuntimeLabelNotificationHandler) HandleUpdate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *RuntimeLabelNotificationHandler) HandleDelete(ctx context.Context, label Label) error {
	if label.Key != model.ScenariosKey {
		log.C(ctx).Infof("label %v is not scenarios", label)
		return nil
	}

	if len(label.RuntimeID) == 0 {
		log.C(ctx).Infof("label %v is not runtimes", label)
		return nil
	}

	ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	runtime, err := a.RuntimeGetter.GetByID(ctx, label.TenantID, label.RuntimeID)
	if err != nil {
		return err
	}

	if err := a.scriptRunner.DeleteDependency(ctx, "dep-"+runtime.Name, "runtime.yaml"); err != nil {
		return err
	}

	return nil
}

func (a *RuntimeLabelNotificationHandler) handle(ctx context.Context, label Label) error {
	if label.Key != model.ScenariosKey {
		log.C(ctx).Infof("label %v is not scenarios", label)
		return nil
	}

	if len(label.RuntimeID) == 0 {
		log.C(ctx).Infof("label %v is not runtimes", label)
		return nil
	}

	ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	runtime, err := a.RuntimeGetter.GetByID(ctx, label.TenantID, label.RuntimeID)
	if err != nil {
		return err
	}
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

	return nil
}
