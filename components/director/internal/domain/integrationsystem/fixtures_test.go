package integrationsystem_test

import (
	"database/sql/driver"
	"errors"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kyma-incubator/compass/components/director/internal/domain/integrationsystem"
	"github.com/kyma-incubator/compass/components/director/internal/model"
)

const (
	testTenant         = "tnt"
	testExternalTenant = "external-tnt"
	testID             = "foo"
	testName           = "bar"
	testPageSize       = 3
	testCursor         = ""
)

var (
	testError        = errors.New("test error")
	testDescription  = "bazz"
	testTableColumns = []string{"id", "name", "description"}
)

func fixModelIntegrationSystem(id, name string) *model.IntegrationSystem {
	return &model.IntegrationSystem{
		ID:          id,
		Name:        name,
		Description: &testDescription,
	}
}

func fixGQLIntegrationSystem(id, name string) *externalschema.IntegrationSystem {
	return &externalschema.IntegrationSystem{
		ID:          id,
		Name:        name,
		Description: &testDescription,
	}
}

func fixModelIntegrationSystemInput(name string) model.IntegrationSystemInput {
	return model.IntegrationSystemInput{
		Name:        name,
		Description: &testDescription,
	}
}

func fixGQLIntegrationSystemInput(name string) externalschema.IntegrationSystemInput {
	return externalschema.IntegrationSystemInput{
		Name:        name,
		Description: &testDescription,
	}
}

func fixEntityIntegrationSystem(id, name string) *integrationsystem.Entity {
	return &integrationsystem.Entity{
		ID:          id,
		Name:        name,
		Description: &testDescription,
	}
}

type sqlRow struct {
	id          string
	name        string
	description *string
}

func fixSQLRows(rows []sqlRow) *sqlmock.Rows {
	out := sqlmock.NewRows(testTableColumns)
	for _, row := range rows {
		out.AddRow(row.id, row.name, row.description)
	}
	return out
}

func fixIntegrationSystemCreateArgs(ent integrationsystem.Entity) []driver.Value {
	return []driver.Value{ent.ID, ent.Name, ent.Description}
}

func fixModelIntegrationSystemPage(intSystems []*model.IntegrationSystem) model.IntegrationSystemPage {
	return model.IntegrationSystemPage{
		Data: intSystems,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(intSystems),
	}
}

func fixGQLIntegrationSystemPage(intSystems []*externalschema.IntegrationSystem) externalschema.IntegrationSystemPage {
	return externalschema.IntegrationSystemPage{
		Data: intSystems,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(intSystems),
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

func fixModelSystemAuth(id, intSysID string, auth *model.Auth) model.SystemAuth {
	return model.SystemAuth{
		ID:                  id,
		TenantID:            nil,
		IntegrationSystemID: &intSysID,
		Value:               auth,
	}
}

func fixGQLSystemAuth(id string, auth *externalschema.Auth) *externalschema.SystemAuth {
	return &externalschema.SystemAuth{
		ID:   id,
		Auth: auth,
	}
}
