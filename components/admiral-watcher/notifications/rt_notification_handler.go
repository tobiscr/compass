package notifications

import (
	"context"
	"encoding/json"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/pkg/errors"
)

type Runtime struct {
	ID                  string `db:"id" json:"id"`
	TenantID            string `db:"tenant_id" json:"tenant_id"`
	Name                string `db:"name" json:"name"`
	ProviderName        string `db:"provider_name" json:"provider_name"`
	Description         string `db:"description" json:"description"`
	StatusCondition     string `db:"status_condition" json:"status_condition"`
	HealthCheckURL      string `db:"healthcheck_url" json:"healthcheck_url"`
	IntegrationSystemID string `db:"integration_system_id" json:"integration_system_id"`
}

type RtNotificationHandler struct {
	ScriptRunner script.Runner
}

func (l *RtNotificationHandler) HandleCreate(ctx context.Context, data []byte) error {
	entity := Runtime{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal runtime: %s", err)
	}

	if entity.Name != "runtime-poc" {
		log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", "runtime-poc", entity.Name)
		return nil
	}

	if err := l.ScriptRunner.RegisterRuntime(ctx, "runtime.yaml"); err != nil {
		return err
	}

	log.C(ctx).Infof("Successfully handled create event for runtime %v", entity)
	return nil
}

func (l *RtNotificationHandler) HandleUpdate(ctx context.Context, data []byte) error {
	entity := Runtime{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal runtime: %s", err)
	}

	if entity.Name != "runtime-poc" {
		log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", "runtime-poc", entity.Name)
		return nil
	}

	log.C(ctx).Infof("Successfully handled update event for runtime %v", entity)
	return nil
}

func (l *RtNotificationHandler) HandleDelete(ctx context.Context, data []byte) error {
	entity := Runtime{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal runtime: %s", err)
	}

	if entity.Name != "runtime-poc" {
		log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", "runtime-poc", entity.Name)
		return nil
	}

	if err := l.ScriptRunner.DeleteRuntime(ctx, "runtime.yaml"); err != nil {
		return err
	}

	log.C(ctx).Infof("Successfully handled delete event for runtime %v", entity)
	return nil
}
