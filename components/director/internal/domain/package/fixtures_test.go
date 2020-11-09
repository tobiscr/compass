package mp_package_test

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	mp_package "github.com/kyma-incubator/compass/components/director/internal/domain/package"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/pagination"
)

func fixModelAPIDefinition(id string, pkgID string, name, description string, group string) *model.APIDefinition {
	return &model.APIDefinition{
		ID:          id,
		PackageID:   pkgID,
		Name:        name,
		Description: &description,
		Group:       &group,
	}
}

func fixGQLAPIDefinition(id string, pkgID string, name, description string, group string) *externalschema.APIDefinition {
	return &externalschema.APIDefinition{
		ID:          id,
		PackageID:   pkgID,
		Name:        name,
		Description: &description,
		Group:       &group,
	}
}

func fixAPIDefinitionPage(apiDefinitions []*model.APIDefinition) *model.APIDefinitionPage {
	return &model.APIDefinitionPage{
		Data: apiDefinitions,
		PageInfo: &pagination.Page{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(apiDefinitions),
	}
}

func fixGQLAPIDefinitionPage(apiDefinitions []*externalschema.APIDefinition) *externalschema.APIDefinitionPage {
	return &externalschema.APIDefinitionPage{
		Data: apiDefinitions,
		PageInfo: &externalschema.PageInfo{
			StartCursor: "start",
			EndCursor:   "end",
			HasNextPage: false,
		},
		TotalCount: len(apiDefinitions),
	}
}

func fixModelEventAPIDefinition(id string, packageID string, name, description string, group string) *model.EventDefinition {
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
func fixGQLEventDefinition(id string, packageID string, name, description string, group string) *externalschema.EventDefinition {
	return &externalschema.EventDefinition{
		ID:          id,
		PackageID:   packageID,
		Name:        name,
		Description: &description,
		Group:       &group,
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

var (
	docKind  = "fookind"
	docTitle = "footitle"
	docData  = "foodata"
	docCLOB  = externalschema.CLOB(docData)
	desc     = "Lorem Ipsum"
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

const (
	packageID        = "ddddddddd-dddd-dddd-dddd-dddddddddddd"
	appID            = "aaaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	tenantID         = "ttttttttt-tttt-tttt-tttt-tttttttttttt"
	externalTenantID = "eeeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
)

func fixPackageModel(t *testing.T, name, desc string) *model.Package {
	return &model.Package{
		ID:                             packageID,
		TenantID:                       tenantID,
		ApplicationID:                  appID,
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            fixModelAuth(),
	}
}

func fixGQLPackage(id, name, desc string) *externalschema.Package {
	schema := externalschema.JSONSchema(`{"$id":"https://example.com/person.schema.json","$schema":"http://json-schema.org/draft-07/schema#","properties":{"age":{"description":"Age in years which must be equal to or greater than zero.","minimum":0,"type":"integer"},"firstName":{"description":"The person's first name.","type":"string"},"lastName":{"description":"The person's last name.","type":"string"}},"title":"Person","type":"object"}`)
	return &externalschema.Package{
		ID:                             id,
		Name:                           name,
		Description:                    &desc,
		InstanceAuthRequestInputSchema: &schema,
		DefaultInstanceAuth:            fixGQLAuth(),
	}
}

func fixGQLPackageCreateInput(name, description string) externalschema.PackageCreateInput {
	basicCredentialDataInput := externalschema.BasicCredentialDataInput{
		Username: "test",
		Password: "pwd",
	}

	credentialDataInput := externalschema.CredentialDataInput{Basic: &basicCredentialDataInput}
	defaultAuth := externalschema.AuthInput{
		Credential: &credentialDataInput,
	}

	return externalschema.PackageCreateInput{
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: fixBasicInputSchema(),
		DefaultInstanceAuth:            &defaultAuth,
		APIDefinitions: []*externalschema.APIDefinitionInput{
			{Name: "api1", TargetURL: "foo.bar"},
			{Name: "api2", TargetURL: "foo.bar2"},
		},
		EventDefinitions: []*externalschema.EventDefinitionInput{
			{Name: "event1", Description: &desc},
			{Name: "event2", Description: &desc},
		},
		Documents: []*externalschema.DocumentInput{
			{DisplayName: "doc1", Kind: &docKind},
			{DisplayName: "doc2", Kind: &docKind},
		},
	}
}

func fixModelPackageCreateInput(name, description string) model.PackageCreateInput {
	basicCredentialDataInput := model.BasicCredentialDataInput{
		Username: "test",
		Password: "pwd",
	}
	authInput := model.AuthInput{
		Credential: &model.CredentialDataInput{Basic: &basicCredentialDataInput},
	}

	return model.PackageCreateInput{
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &authInput,
		APIDefinitions: []*model.APIDefinitionInput{
			{Name: "api1", TargetURL: "foo.bar"},
			{Name: "api2", TargetURL: "foo.bar2"},
		},
		EventDefinitions: []*model.EventDefinitionInput{
			{Name: "event1", Description: &desc},
			{Name: "event2", Description: &desc},
		},
		Documents: []*model.DocumentInput{
			{DisplayName: "doc1", Kind: &docKind},
			{DisplayName: "doc2", Kind: &docKind},
		},
	}
}

func fixGQLPackageUpdateInput(name, description string) externalschema.PackageUpdateInput {
	basicCredentialDataInput := externalschema.BasicCredentialDataInput{
		Username: "test",
		Password: "pwd",
	}

	credentialDataInput := externalschema.CredentialDataInput{Basic: &basicCredentialDataInput}
	defaultAuth := externalschema.AuthInput{
		Credential: &credentialDataInput,
	}

	return externalschema.PackageUpdateInput{
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: fixBasicInputSchema(),
		DefaultInstanceAuth:            &defaultAuth,
	}
}

func fixModelPackageUpdateInput(t *testing.T, name, description string) model.PackageUpdateInput {
	basicCredentialDataInput := model.BasicCredentialDataInput{
		Username: "test",
		Password: "pwd",
	}
	authInput := model.AuthInput{
		Credential: &model.CredentialDataInput{Basic: &basicCredentialDataInput},
	}

	return model.PackageUpdateInput{
		Name:                           name,
		Description:                    &description,
		InstanceAuthRequestInputSchema: fixBasicSchema(),
		DefaultInstanceAuth:            &authInput,
	}
}

func fixModelAuthInput(headers map[string][]string) *model.AuthInput {
	return &model.AuthInput{
		AdditionalHeaders: headers,
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

func fixEntityPackage(id, name, desc string) *mp_package.Entity {
	descSQL := sql.NullString{desc, true}
	schemaSQL := sql.NullString{
		String: inputSchemaString(),
		Valid:  true,
	}
	authSQL := sql.NullString{
		String: `{"Credential":{"Basic":{"Username":"foo","Password":"bar"},"Oauth":null},"AdditionalHeaders":{"test":["foo","bar"]},"AdditionalQueryParams":{"test":["foo","bar"]},"RequestAuth":{"Csrf":{"TokenEndpointURL":"foo.url","Credential":{"Basic":{"Username":"boo","Password":"far"},"Oauth":null},"AdditionalHeaders":{"test":["foo","bar"]},"AdditionalQueryParams":{"test":["foo","bar"]}}}}`,
		Valid:  true,
	}

	return &mp_package.Entity{
		ID:                            id,
		TenantID:                      tenantID,
		ApplicationID:                 appID,
		Name:                          name,
		Description:                   descSQL,
		InstanceAuthRequestJSONSchema: schemaSQL,
		DefaultInstanceAuth:           authSQL,
	}
}

func fixPackageColumns() []string {
	return []string{"id", "tenant_id", "app_id", "name", "description", "instance_auth_request_json_schema", "default_instance_auth"}
}

func fixPackageRow(id, placeholder string) []driver.Value {
	return []driver.Value{id, tenantID, appID, "foo", "bar", fixSchema(), fixDefaultAuth()}
}

func fixPackageCreateArgs(defAuth, schema string, pkg *model.Package) []driver.Value {
	return []driver.Value{packageID, tenantID, appID, pkg.Name, pkg.Description, schema, defAuth}
}

func fixDefaultAuth() string {
	return `{"Credential":{"Basic":{"Username":"foo","Password":"bar"},"Oauth":null},"AdditionalHeaders":{"test":["foo","bar"]},"AdditionalQueryParams":{"test":["foo","bar"]},"RequestAuth":{"Csrf":{"TokenEndpointURL":"foo.url","Credential":{"Basic":{"Username":"boo","Password":"far"},"Oauth":null},"AdditionalHeaders":{"test":["foo","bar"]},"AdditionalQueryParams":{"test":["foo","bar"]}}}}`
}

func inputSchemaString() string {
	return `{"$id":"https://example.com/person.schema.json","$schema":"http://json-schema.org/draft-07/schema#","properties":{"age":{"description":"Age in years which must be equal to or greater than zero.","minimum":0,"type":"integer"},"firstName":{"description":"The person's first name.","type":"string"},"lastName":{"description":"The person's last name.","type":"string"}},"title":"Person","type":"object"}`
}

func fixBasicInputSchema() *externalschema.JSONSchema {
	sch := inputSchemaString()
	jsonSchema := externalschema.JSONSchema(sch)
	return &jsonSchema
}

func fixBasicSchema() *string {
	sch := inputSchemaString()
	return &sch
}

func fixSchema() string {
	return `{"$id":"https://example.com/person.schema.json","$schema":"http://json-schema.org/draft-07/schema#","properties":{"age":{"description":"Age in years which must be equal to or greater than zero.","minimum":0,"type":"integer"},"firstName":{"description":"The person's first name.","type":"string"},"lastName":{"description":"The person's last name.","type":"string"}},"title":"Person","type":"object"}`
}

func fixModelPackageInstanceAuth(id string) *model.PackageInstanceAuth {
	status := model.PackageInstanceAuthStatus{
		Condition: model.PackageInstanceAuthStatusConditionPending,
		Timestamp: time.Time{},
		Message:   "test-message",
		Reason:    "test-reason",
	}

	context := "ctx"
	params := "test-param"
	return &model.PackageInstanceAuth{
		ID:          id,
		PackageID:   packageID,
		Tenant:      tenantID,
		Context:     &context,
		InputParams: &params,
		Auth:        fixModelAuth(),
		Status:      &status,
	}
}

func fixGQLPackageInstanceAuth(id string) *externalschema.PackageInstanceAuth {
	msg := "test-message"
	reason := "test-reason"
	status := externalschema.PackageInstanceAuthStatus{
		Condition: externalschema.PackageInstanceAuthStatusConditionPending,
		Timestamp: externalschema.Timestamp{},
		Message:   msg,
		Reason:    reason,
	}

	params := externalschema.JSON("test-param")
	ctx := externalschema.JSON("ctx")
	return &externalschema.PackageInstanceAuth{
		ID:          id,
		Context:     &ctx,
		InputParams: &params,
		Auth:        fixGQLAuth(),
		Status:      &status,
	}
}

func fixFetchRequest(url string, objectType model.FetchRequestReferenceObjectType, timestamp time.Time) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     "foo",
		Tenant: tenantID,
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

func fixFetchRequestWithCondition(url string, objectType model.FetchRequestReferenceObjectType, timestamp time.Time, condition model.FetchRequestStatusCondition) *model.FetchRequest {
	return &model.FetchRequest{
		ID:     "foo",
		Tenant: tenantID,
		URL:    url,
		Auth:   nil,
		Mode:   "SINGLE",
		Filter: nil,
		Status: &model.FetchRequestStatus{
			Condition: condition,
			Timestamp: timestamp,
		},
		ObjectType: objectType,
		ObjectID:   "foo",
	}
}
