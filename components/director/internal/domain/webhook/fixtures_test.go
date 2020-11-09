package webhook_test

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

func fixModelWebhook(id, appID, tenant, url string) *model.Webhook {
	return &model.Webhook{
		ID:            id,
		ApplicationID: appID,
		Tenant:        tenant,
		Type:          model.WebhookTypeConfigurationChanged,
		URL:           url,
		Auth:          &model.Auth{},
	}
}

func fixGQLWebhook(id, appID, url string) *externalschema.Webhook {
	return &externalschema.Webhook{
		ID:            id,
		ApplicationID: appID,
		Type:          externalschema.ApplicationWebhookTypeConfigurationChanged,
		URL:           url,
		Auth:          &externalschema.Auth{},
	}
}

func fixModelWebhookInput(url string) *model.WebhookInput {
	return &model.WebhookInput{
		Type: model.WebhookTypeConfigurationChanged,
		URL:  url,
		Auth: &model.AuthInput{},
	}
}

func fixGQLWebhookInput(url string) *externalschema.WebhookInput {
	return &externalschema.WebhookInput{
		Type: externalschema.ApplicationWebhookTypeConfigurationChanged,
		URL:  url,
		Auth: &externalschema.AuthInput{},
	}
}
