package application_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/internal/repo"

	"github.com/kyma-incubator/compass/components/director/internal/domain/application"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/stretchr/testify/require"
)

var (
	testURL      = "https://foo.bar"
	intSysID     = "iiiiiiiii-iiii-iiii-iiii-iiiiiiiiiiii"
	providerName = "provider name"
)

func fixApplicationPage(applications []*model.Application) *model.ApplicationPage {
	return &model.ApplicationPage{
		Data: applications,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(applications),
	}
}

func fixGQLApplicationPage(applications []*externalschema.Application) *externalschema.ApplicationPage {
	return &externalschema.ApplicationPage{
		Data: applications,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(applications),
	}
}

func fixModelApplication(id, tenant, name, description string) *model.Application {
	return &model.Application{
		ID:     id,
		Tenant: tenant,
		Status: &model.ApplicationStatus{
			Condition: model.ApplicationStatusConditionInitial,
		},
		Name:        name,
		Description: &description,
	}
}

func fixModelApplicationWithAllUpdatableFields(id, tenant, name, description, url string, conditionStatus model.ApplicationStatusCondition, conditionTimestamp time.Time) *model.Application {
	return &model.Application{
		ID:     id,
		Tenant: tenant,
		Status: &model.ApplicationStatus{
			Condition: conditionStatus,
			Timestamp: conditionTimestamp,
		},
		IntegrationSystemID: &intSysID,
		Name:                name,
		Description:         &description,
		HealthCheckURL:      &url,
		ProviderName:        &providerName,
	}
}

func fixGQLApplication(id, name, description string) *externalschema.Application {
	return &externalschema.Application{
		ID: id,
		Status: &externalschema.ApplicationStatus{
			Condition: externalschema.ApplicationStatusConditionInitial,
		},
		Name:        name,
		Description: &description,
	}
}

func fixDetailedModelApplication(t *testing.T, id, tenant, name, description string) *model.Application {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &model.Application{
		ID: id,
		Status: &model.ApplicationStatus{
			Condition: model.ApplicationStatusConditionInitial,
			Timestamp: time,
		},
		Name:                name,
		Description:         &description,
		Tenant:              tenant,
		HealthCheckURL:      &testURL,
		IntegrationSystemID: &intSysID,
		ProviderName:        &providerName,
	}
}

func fixDetailedGQLApplication(t *testing.T, id, name, description string) *externalschema.Application {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &externalschema.Application{
		ID: id,
		Status: &externalschema.ApplicationStatus{
			Condition: externalschema.ApplicationStatusConditionInitial,
			Timestamp: externalschema.Timestamp(time),
		},
		Name:                name,
		Description:         &description,
		HealthCheckURL:      &testURL,
		IntegrationSystemID: &intSysID,
		ProviderName:        str.Ptr("provider name"),
	}
}

func fixDetailedEntityApplication(t *testing.T, id, tenant, name, description string) *application.Entity {
	ts, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &application.Entity{
		ID:                  id,
		TenantID:            tenant,
		Name:                name,
		Description:         repo.NewValidNullableString(description),
		StatusCondition:     string(model.ApplicationStatusConditionInitial),
		StatusTimestamp:     ts,
		HealthCheckURL:      repo.NewValidNullableString(testURL),
		IntegrationSystemID: repo.NewNullableString(&intSysID),
		ProviderName:        repo.NewNullableString(&providerName),
	}
}

func fixModelApplicationRegisterInput(name, description string) model.ApplicationRegisterInput {
	desc := "Sample"
	kind := "test"
	return model.ApplicationRegisterInput{
		Name:        name,
		Description: &description,
		Labels: map[string]interface{}{
			"test": []string{"val", "val2"},
		},
		HealthCheckURL:      &testURL,
		IntegrationSystemID: &intSysID,
		ProviderName:        &providerName,
		Webhooks: []*model.WebhookInput{
			{URL: "webhook1.foo.bar"},
			{URL: "webhook2.foo.bar"},
		},
		Packages: []*model.PackageCreateInput{
			{
				Name: "foo",
				APIDefinitions: []*model.APIDefinitionInput{
					{Name: "api1", TargetURL: "foo.bar"},
					{Name: "api2", TargetURL: "foo.bar2"},
				},
				EventDefinitions: []*model.EventDefinitionInput{
					{Name: "event1", Description: &desc},
					{Name: "event2", Description: &desc},
				},
				Documents: []*model.DocumentInput{
					{DisplayName: "doc1", Kind: &kind},
					{DisplayName: "doc2", Kind: &kind},
				},
			},
		},
	}
}

func fixModelApplicationUpdateInput(name, description, url string, statusCondition model.ApplicationStatusCondition) model.ApplicationUpdateInput {
	return model.ApplicationUpdateInput{
		Description:         &description,
		HealthCheckURL:      &url,
		IntegrationSystemID: &intSysID,
		ProviderName:        &providerName,
		StatusCondition:     &statusCondition,
	}
}

func fixModelApplicationUpdateInputStatus(statusCondition model.ApplicationStatusCondition) model.ApplicationUpdateInput {
	return model.ApplicationUpdateInput{
		StatusCondition: &statusCondition,
	}
}

func fixGQLApplicationRegisterInput(name, description string) externalschema.ApplicationRegisterInput {
	labels := externalschema.Labels{
		"test": []string{"val", "val2"},
	}
	kind := "test"
	desc := "Sample"
	return externalschema.ApplicationRegisterInput{
		Name:                name,
		Description:         &description,
		Labels:              &labels,
		HealthCheckURL:      &testURL,
		IntegrationSystemID: &intSysID,
		ProviderName:        &providerName,
		Webhooks: []*externalschema.WebhookInput{
			{URL: "webhook1.foo.bar"},
			{URL: "webhook2.foo.bar"},
		},
		Packages: []*externalschema.PackageCreateInput{
			{
				Name: "foo",
				APIDefinitions: []*externalschema.APIDefinitionInput{
					{Name: "api1", TargetURL: "foo.bar"},
					{Name: "api2", TargetURL: "foo.bar2"},
				},
				EventDefinitions: []*externalschema.EventDefinitionInput{
					{Name: "event1", Description: &desc},
					{Name: "event2", Description: &desc},
				},
				Documents: []*externalschema.DocumentInput{
					{DisplayName: "doc1", Kind: &kind},
					{DisplayName: "doc2", Kind: &kind},
				},
			},
		},
	}
}

func fixGQLApplicationUpdateInput(name, description, url string, statusCondition externalschema.ApplicationStatusCondition) externalschema.ApplicationUpdateInput {
	return externalschema.ApplicationUpdateInput{
		Description:         &description,
		HealthCheckURL:      &url,
		IntegrationSystemID: &intSysID,
		ProviderName:        &providerName,
		StatusCondition:     &statusCondition,
	}
}

var (
	docKind  = "fookind"
	docTitle = "footitle"
	docData  = "foodata"
	docCLOB  = externalschema.CLOB(docData)
)

func fixModelDocument(packageID, id string) *model.Document {
	return &model.Document{
		PackageID: packageID,
		ID:        id,
		Title:     docTitle,
		Format:    model.DocumentFormatMarkdown,
		Kind:      &docKind,
		Data:      &docData,
	}
}

func fixModelDocumentPage(documents []*model.Document) *model.DocumentPage {
	return &model.DocumentPage{
		Data: documents,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(documents),
	}
}

func fixGQLDocument(id string) *externalschema.Document {
	return &externalschema.Document{
		ID:     id,
		Title:  docTitle,
		Format: externalschema.DocumentFormatMarkdown,
		Kind:   &docKind,
		Data:   &docCLOB,
	}
}

func fixGQLDocumentPage(documents []*externalschema.Document) *externalschema.DocumentPage {
	return &externalschema.DocumentPage{
		Data: documents,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(documents),
	}
}

func fixModelWebhook(appID, id string) *model.Webhook {
	return &model.Webhook{
		ApplicationID: appID,
		ID:            id,
		Type:          model.WebhookTypeConfigurationChanged,
		URL:           "foourl",
		Auth:          &model.Auth{},
	}
}

func fixGQLWebhook(id string) *externalschema.Webhook {
	return &externalschema.Webhook{
		ID:   id,
		Type: externalschema.ApplicationWebhookTypeConfigurationChanged,
		URL:  "foourl",
		Auth: &externalschema.Auth{},
	}
}

func fixEventAPIDefinitionPage(eventAPIDefinitions []*model.EventDefinition) *model.EventDefinitionPage {
	return &model.EventDefinitionPage{
		Data: eventAPIDefinitions,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(eventAPIDefinitions),
	}
}

func fixGQLEventDefinitionPage(eventAPIDefinitions []*externalschema.EventDefinition) *externalschema.EventDefinitionPage {
	return &externalschema.EventDefinitionPage{
		Data: eventAPIDefinitions,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(eventAPIDefinitions),
	}
}

func fixModelEventAPIDefinition(id string, appId, packageID string, name, description string, group string) *model.EventDefinition {
	return &model.EventDefinition{
		ID:          id,
		PackageID:   packageID,
		Name:        name,
		Description: &description,
		Group:       &group,
	}
}
func fixMinModelEventAPIDefinition(id, placeholder string) *model.EventDefinition {
	return &model.EventDefinition{ID: id, Tenant: "ttttttttt-tttt-tttt-tttt-tttttttttttt",
		PackageID: "ppppppppp-pppp-pppp-pppp-pppppppppppp", Name: placeholder}
}
func fixGQLEventDefinition(id string, appId, packageID string, name, description string, group string) *externalschema.EventDefinition {
	return &externalschema.EventDefinition{
		ID:          id,
		PackageID:   packageID,
		Name:        name,
		Description: &description,
		Group:       &group,
	}
}

func fixFetchRequest(url string, objectType model.FetchRequestReferenceObjectType, timestamp time.Time) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     "foo",
		Tenant: "tenant",
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: model.FetchRequestStatusConditionInitial,
			Timestamp: timestamp,
		},
		ObjectType: objectType,
		ObjectID:   "foo",
	}
}

func fixLabelInput(key string, value string, objectID string, objectType model.LabelableObject) *model.LabelInput {
	return &model.LabelInput{
		Key:        key,
		Value:      value,
		ObjectID:   objectID,
		ObjectType: objectType,
	}
}

func fixModelApplicationEventingConfiguration(t *testing.T, rawURL string) *model.ApplicationEventingConfiguration {
	validURL, err := url.Parse(rawURL)
	require.NoError(t, err)
	require.NotNil(t, validURL)
	return &model.ApplicationEventingConfiguration{
		EventingConfiguration: model.EventingConfiguration{
			DefaultURL: *validURL,
		},
	}
}

func fixGQLApplicationEventingConfiguration(url string) *externalschema.ApplicationEventingConfiguration {
	return &externalschema.ApplicationEventingConfiguration{
		DefaultURL: url,
	}
}

func fixModelPackage(id, tenantID, appId, name, description string) *model.Package {
	return &model.Package{
		ID:                             id,
		TenantID:                       tenantID,
		ApplicationID:                  appId,
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: nil,
		DefaultInstanceAuth:            nil,
	}
}

func fixGQLPackage(id, appId, name, description string) *externalschema.Package {
	return &externalschema.Package{
		ID:                             id,
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: nil,
		DefaultInstanceAuth:            nil,
	}
}

func fixGQLPackagePage(packages []*externalschema.Package) *externalschema.PackagePage {
	return &externalschema.PackagePage{
		Data: packages,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(packages),
	}
}

func fixPackagePage(packages []*model.Package) *model.PackagePage {
	return &model.PackagePage{
		Data: packages,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(packages),
	}
}
