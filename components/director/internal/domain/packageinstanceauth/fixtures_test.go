package packageinstanceauth_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kyma-incubator/compass/components/director/internal/domain/packageinstanceauth"

	"github.com/stretchr/testify/require"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

var (
	testID             = "foo"
	testPackageID      = "bar"
	testTenant         = "baz"
	testExternalTenant = "foobaz"
	testContext        = `{"foo": "bar"}`
	testInputParams    = `{"bar": "baz"}`
	testError          = errors.New("test")
	testTime           = time.Now()
	testTableColumns   = []string{"id", "tenant_id", "package_id", "context", "input_params", "auth_value", "status_condition", "status_timestamp", "status_message", "status_reason"}
)

func fixModelPackageInstanceAuth(id, packageID, tenant string, auth *model.Auth, status *model.PackageInstanceAuthStatus) *model.PackageInstanceAuth {
	pia := fixModelPackageInstanceAuthWithoutContextAndInputParams(id, packageID, tenant, auth, status)
	pia.Context = &testContext
	pia.InputParams = &testInputParams

	return pia
}
func fixModelPackageInstanceAuthWithoutContextAndInputParams(id, packageID, tenant string, auth *model.Auth, status *model.PackageInstanceAuthStatus) *model.PackageInstanceAuth {
	return &model.PackageInstanceAuth{
		ID:        id,
		PackageID: packageID,
		Tenant:    tenant,
		Auth:      auth,
		Status:    status,
	}
}

func fixGQLPackageInstanceAuth(id string, auth *externalschema.Auth, status *externalschema.PackageInstanceAuthStatus) *externalschema.PackageInstanceAuth {
	context := externalschema.JSON(testContext)
	inputParams := externalschema.JSON(testInputParams)

	out := fixGQLPackageInstanceAuthWithoutContextAndInputParams(id, auth, status)
	out.Context = &context
	out.InputParams = &inputParams

	return out
}

func fixGQLPackageInstanceAuthWithoutContextAndInputParams(id string, auth *externalschema.Auth, status *externalschema.PackageInstanceAuthStatus) *externalschema.PackageInstanceAuth {
	return &externalschema.PackageInstanceAuth{
		ID:     id,
		Auth:   auth,
		Status: status,
	}
}

func fixModelStatusSucceeded() *model.PackageInstanceAuthStatus {
	return &model.PackageInstanceAuthStatus{
		Condition: model.PackageInstanceAuthStatusConditionSucceeded,
		Timestamp: testTime,
		Message:   "Credentials were provided.",
		Reason:    "CredentialsProvided",
	}
}

func fixModelStatusPending() *model.PackageInstanceAuthStatus {
	return &model.PackageInstanceAuthStatus{
		Condition: model.PackageInstanceAuthStatusConditionPending,
		Timestamp: testTime,
		Message:   "Credentials were not yet provided.",
		Reason:    "CredentialsNotProvided",
	}
}

func fixGQLStatusSucceeded() *externalschema.PackageInstanceAuthStatus {
	return &externalschema.PackageInstanceAuthStatus{
		Condition: externalschema.PackageInstanceAuthStatusConditionSucceeded,
		Timestamp: externalschema.Timestamp(testTime),
		Message:   "Credentials were provided.",
		Reason:    "CredentialsProvided",
	}
}

func fixGQLStatusPending() *externalschema.PackageInstanceAuthStatus {
	return &externalschema.PackageInstanceAuthStatus{
		Condition: externalschema.PackageInstanceAuthStatusConditionPending,
		Timestamp: externalschema.Timestamp(testTime),
		Message:   "Credentials were not yet provided.",
		Reason:    "CredentialsNotProvided",
	}
}

func fixModelStatusInput(condition model.PackageInstanceAuthSetStatusConditionInput, message, reason string) *model.PackageInstanceAuthStatusInput {
	return &model.PackageInstanceAuthStatusInput{
		Condition: condition,
		Message:   message,
		Reason:    reason,
	}
}

func fixGQLStatusInput(condition externalschema.PackageInstanceAuthSetStatusConditionInput, message, reason string) *externalschema.PackageInstanceAuthStatusInput {
	return &externalschema.PackageInstanceAuthStatusInput{
		Condition: condition,
		Message:   message,
		Reason:    reason,
	}
}

func fixModelRequestInput() *model.PackageInstanceAuthRequestInput {
	return &model.PackageInstanceAuthRequestInput{
		Context:     &testContext,
		InputParams: &testInputParams,
	}
}

func fixGQLRequestInput() *externalschema.PackageInstanceAuthRequestInput {
	context := externalschema.JSON(testContext)
	inputParams := externalschema.JSON(testInputParams)

	return &externalschema.PackageInstanceAuthRequestInput{
		Context:     &context,
		InputParams: &inputParams,
	}
}

func fixModelSetInput() *model.PackageInstanceAuthSetInput {
	return &model.PackageInstanceAuthSetInput{
		Auth:   fixModelAuthInput(),
		Status: fixModelStatusInput(model.PackageInstanceAuthSetStatusConditionInputSucceeded, "foo", "bar"),
	}
}

func fixGQLSetInput() *externalschema.PackageInstanceAuthSetInput {
	return &externalschema.PackageInstanceAuthSetInput{
		Auth:   fixGQLAuthInput(),
		Status: fixGQLStatusInput(externalschema.PackageInstanceAuthSetStatusConditionInputSucceeded, "foo", "bar"),
	}
}

func fixEntityPackageInstanceAuth(t *testing.T, id, packageID, tenant string, auth *model.Auth, status *model.PackageInstanceAuthStatus) *packageinstanceauth.Entity {
	out := fixEntityPackageInstanceAuthWithoutContextAndInputParams(t, id, packageID, tenant, auth, status)
	out.Context = sql.NullString{Valid: true, String: testContext}
	out.InputParams = sql.NullString{Valid: true, String: testInputParams}

	return out
}

func fixEntityPackageInstanceAuthWithoutContextAndInputParams(t *testing.T, id, packageID, tenant string, auth *model.Auth, status *model.PackageInstanceAuthStatus) *packageinstanceauth.Entity {
	out := packageinstanceauth.Entity{
		ID:        id,
		PackageID: packageID,
		TenantID:  tenant,
	}

	if auth != nil {
		marshalled, err := json.Marshal(auth)
		require.NoError(t, err)
		out.AuthValue = sql.NullString{
			String: string(marshalled),
			Valid:  true,
		}
	}

	if status != nil {
		out.StatusCondition = string(status.Condition)
		out.StatusTimestamp = status.Timestamp
		out.StatusMessage = status.Message
		out.StatusReason = status.Reason
	}

	return &out
}

func fixModelAuth() *model.Auth {
	return &model.Auth{
		Credential: model.CredentialData{
			Basic: &model.BasicCredentialData{
				Username: "foo",
				Password: "bar",
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
	}
}

func fixModelAuthInput() *model.AuthInput {
	return &model.AuthInput{
		Credential: &model.CredentialDataInput{
			Basic: &model.BasicCredentialDataInput{
				Username: "foo",
				Password: "bar",
			},
		},
	}
}

func fixGQLAuthInput() *externalschema.AuthInput {
	return &externalschema.AuthInput{
		Credential: &externalschema.CredentialDataInput{
			Basic: &externalschema.BasicCredentialDataInput{
				Username: "foo",
				Password: "bar",
			},
		},
	}
}

type sqlRow struct {
	id              string
	tenantID        string
	packageID       string
	context         sql.NullString
	inputParams     sql.NullString
	authValue       sql.NullString
	statusCondition string
	statusTimestamp time.Time
	statusMessage   string
	statusReason    string
}

func fixSQLRows(rows []sqlRow) *sqlmock.Rows {
	out := sqlmock.NewRows(testTableColumns)
	for _, row := range rows {
		out.AddRow(row.id, row.tenantID, row.packageID, row.context, row.inputParams, row.authValue, row.statusCondition, row.statusTimestamp, row.statusMessage, row.statusReason)
	}
	return out
}

func fixSQLRowFromEntity(entity packageinstanceauth.Entity) sqlRow {
	return sqlRow{
		id:              entity.ID,
		tenantID:        entity.TenantID,
		packageID:       entity.PackageID,
		context:         entity.Context,
		inputParams:     entity.InputParams,
		authValue:       entity.AuthValue,
		statusCondition: entity.StatusCondition,
		statusTimestamp: entity.StatusTimestamp,
		statusMessage:   entity.StatusMessage,
		statusReason:    entity.StatusReason,
	}
}

func fixCreateArgs(ent packageinstanceauth.Entity) []driver.Value {
	return []driver.Value{ent.ID, ent.TenantID, ent.PackageID, ent.Context, ent.InputParams, ent.AuthValue, ent.StatusCondition, ent.StatusTimestamp, ent.StatusMessage, ent.StatusReason}
}

func fixSimpleModelPackageInstanceAuth(id string) *model.PackageInstanceAuth {
	return &model.PackageInstanceAuth{
		ID: id,
	}
}

func fixSimpleGQLPackageInstanceAuth(id string) *externalschema.PackageInstanceAuth {
	return &externalschema.PackageInstanceAuth{
		ID: id,
	}
}

func fixModelPackage(id string, requestInputSchema *string, defaultAuth *model.Auth) *model.Package {
	return &model.Package{
		ID:                             id,
		TenantID:                       testTenant,
		ApplicationID:                  "foo",
		Name:                           "test-package",
		InstanceAuthRequestInputSchema: requestInputSchema,
		DefaultInstanceAuth:            defaultAuth,
	}
}
