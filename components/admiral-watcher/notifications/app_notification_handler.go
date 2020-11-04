package notifications

import (
	"context"
	"encoding/json"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/admiral-watcher/script"
	"github.com/pkg/errors"
	"strings"
)

type Entity struct {
	ID                  string `db:"id" json:"id"`
	TenantID            string `db:"tenant_id" json:"tenant_id"`
	Name                string `db:"name" json:"name"`
	ProviderName        string `db:"provider_name" json:"provider_name"`
	Description         string `db:"description" json:"description"`
	StatusCondition     string `db:"status_condition" json:"status_condition"`
	HealthCheckURL      string `db:"healthcheck_url" json:"healthcheck_url"`
	IntegrationSystemID string `db:"integration_system_id" json:"integration_system_id"`
}

type AppNotificationHandler struct {
	ScriptRunner script.Runner
}

func (l *AppNotificationHandler) HandleCreate(ctx context.Context, data []byte) error {
	entity := Entity{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal app: %s", err)
	}

	if entity.Name != "commerce-mock" {
		log.C(ctx).Infof("event is not for the test application %s but for %s, skipping", "commerce-mock", entity.Name)
		return nil
	}

	if strings.ToLower(entity.StatusCondition) != "connected" {
		log.C(ctx).Infof("create event for app %v; status is not Connected; ignoring event")
		return nil
	}

	if err := l.ScriptRunner.RegisterApplication(ctx, "commerce.yaml"); err != nil {
		return err
	}

	log.C(ctx).Infof("Successfully handled create event for app %v", entity)
	return nil
}

func (l *AppNotificationHandler) HandleUpdate(ctx context.Context, data []byte) error {
	entity := Entity{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal app: %s", err)
	}

	if entity.Name != "commerce-mock" {
		log.C(ctx).Infof("event is not for the test application %s but for %s, skipping", "commerce-mock", entity.Name)
		return nil
	}

	if strings.ToLower(entity.StatusCondition) != "connected" {
		log.C(ctx).Infof("create event for app %v; status is not Connected; ignoring event")
		return nil
	}

	if err := l.ScriptRunner.RegisterApplication(ctx, "commerce.yaml"); err != nil {
		return err
	}

	log.C(ctx).Infof("Successfully handled update event for app %v", entity)

	//fetch runtimes in scenario with app
	// create dependencies
	return nil
}

func (l *AppNotificationHandler) HandleDelete(ctx context.Context, data []byte) error {
	entity := Entity{}
	if err := json.Unmarshal(data, &entity); err != nil {
		return errors.Errorf("could not unmarshal app: %s", err)
	}

	if entity.Name != "commerce-mock" {
		log.C(ctx).Infof("event is not for the test application %s but for %s, skipping", "commerce-mock", entity.Name)
		return nil
	}

	if err := l.ScriptRunner.DeleteApplication(ctx, "commerce.yaml"); err != nil {
		return err
	}

	log.C(ctx).Infof("Successfully handled delete event for app %v", entity)
	return nil
}
