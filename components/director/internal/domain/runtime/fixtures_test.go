package runtime_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func fixRuntimePage(runtimes []*model.Runtime) *model.RuntimePage {
	return &model.RuntimePage{
		Data: runtimes,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(runtimes),
	}
}

func fixGQLRuntimePage(runtimes []*externalschema.Runtime) *externalschema.RuntimePage {
	return &externalschema.RuntimePage{
		Data: runtimes,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(runtimes),
	}
}

func fixModelRuntime(t *testing.T, id, tenant, name, description string) *model.Runtime {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &model.Runtime{
		ID:     id,
		Tenant: tenant,
		Status: &model.RuntimeStatus{
			Condition: model.RuntimeStatusConditionInitial,
		},
		Name:              name,
		Description:       &description,
		CreationTimestamp: time,
	}
}

func fixGQLRuntime(t *testing.T, id, name, description string) *externalschema.Runtime {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &externalschema.Runtime{
		ID: id,
		Status: &externalschema.RuntimeStatus{
			Condition: externalschema.RuntimeStatusConditionInitial,
		},
		Name:        name,
		Description: &description,
		Metadata: &externalschema.RuntimeMetadata{
			CreationTimestamp: externalschema.Timestamp(time),
		},
	}
}

func fixDetailedModelRuntime(t *testing.T, id, name, description string) *model.Runtime {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &model.Runtime{
		ID: id,
		Status: &model.RuntimeStatus{
			Condition: model.RuntimeStatusConditionInitial,
			Timestamp: time,
		},
		Name:              name,
		Description:       &description,
		Tenant:            "tenant",
		CreationTimestamp: time,
	}
}

func fixDetailedGQLRuntime(t *testing.T, id, name, description string) *externalschema.Runtime {
	time, err := time.Parse(time.RFC3339, "2002-10-02T10:00:00-05:00")
	require.NoError(t, err)

	return &externalschema.Runtime{
		ID: id,
		Status: &externalschema.RuntimeStatus{
			Condition: externalschema.RuntimeStatusConditionInitial,
			Timestamp: externalschema.Timestamp(time),
		},
		Name:        name,
		Description: &description,
		Metadata: &externalschema.RuntimeMetadata{
			CreationTimestamp: externalschema.Timestamp(time),
		},
	}
}

func fixModelRuntimeInput(name, description string) model.RuntimeInput {
	return model.RuntimeInput{
		Name:        name,
		Description: &description,
		Labels: map[string]interface{}{
			"test": []string{"val", "val2"},
		},
	}
}

func fixGQLRuntimeInput(name, description string) externalschema.RuntimeInput {
	labels := externalschema.Labels{
		"test": []string{"val", "val2"},
	}

	return externalschema.RuntimeInput{
		Name:        name,
		Description: &description,
		Labels:      &labels,
	}
}

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

func fixModelApplication(id, name, description string) *model.Application {
	return &model.Application{
		ID: id,
		Status: &model.ApplicationStatus{
			Condition: model.ApplicationStatusConditionInitial,
		},
		Name:        name,
		Description: &description,
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

func fixModelAuth() *model.Auth {
	return &model.Auth{
		Credential: model.CredentialData{
			Basic: &model.BasicCredentialData{
				Username: "foo",
				Password: "bar",
			},
		},
		AdditionalHeaders:     map[string][]string{"test": {"foo", "bar"}},
		AdditionalQueryParams: map[string][]string{"test": {"foo", "bar"}},
		RequestAuth: &model.CredentialRequestAuth{
			Csrf: &model.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: "foo.url",
				Credential: model.CredentialData{
					Basic: &model.BasicCredentialData{
						Username: "boo",
						Password: "far",
					},
				},
				AdditionalHeaders:     map[string][]string{"test": {"foo", "bar"}},
				AdditionalQueryParams: map[string][]string{"test": {"foo", "bar"}},
			},
		},
	}
}

func fixGQLAuth() *externalschema.Auth {
	return &externalschema.Auth{
		Credential: &externalschema.BasicCredentialData{
			Username: "foo",
			Password: "bar",
		},
		AdditionalHeaders:     &externalschema.HttpHeaders{"test": {"foo", "bar"}},
		AdditionalQueryParams: &externalschema.QueryParams{"test": {"foo", "bar"}},
		RequestAuth: &externalschema.CredentialRequestAuth{
			Csrf: &externalschema.CSRFTokenCredentialRequestAuth{
				TokenEndpointURL: "foo.url",
				Credential: &externalschema.BasicCredentialData{
					Username: "boo",
					Password: "far",
				},
				AdditionalHeaders:     &externalschema.HttpHeaders{"test": {"foo", "bar"}},
				AdditionalQueryParams: &externalschema.QueryParams{"test": {"foo", "bar"}},
			},
		},
	}
}

func fixModelSystemAuth(id, tenant, runtimeID string, auth *model.Auth) model.SystemAuth {
	return model.SystemAuth{
		ID:        id,
		TenantID:  &tenant,
		RuntimeID: &runtimeID,
		Value:     auth,
	}
}

func fixGQLSystemAuth(id string, auth *externalschema.Auth) *externalschema.SystemAuth {
	return &externalschema.SystemAuth{
		ID:   id,
		Auth: auth,
	}
}

func fixModelRuntimeEventingConfiguration(t *testing.T, rawURL string) *model.RuntimeEventingConfiguration {
	validURL := fixValidURL(t, rawURL)
	return &model.RuntimeEventingConfiguration{
		EventingConfiguration: model.EventingConfiguration{
			DefaultURL: validURL,
		},
	}
}

func fixGQLRuntimeEventingConfiguration(url string) *externalschema.RuntimeEventingConfiguration {
	return &externalschema.RuntimeEventingConfiguration{
		DefaultURL: url,
	}
}

func fixValidURL(t *testing.T, rawURL string) url.URL {
	eventingURL, err := url.Parse(rawURL)
	require.NoError(t, err)
	require.NotNil(t, eventingURL)
	return *eventingURL
}
