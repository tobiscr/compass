package packageinstanceauth

import (
	"database/sql"
	"encoding/json"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/repo"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/pkg/errors"
)

//go:generate mockery -name=AuthConverter -output=automock -outpkg=automock -case=underscore
type AuthConverter interface {
	ToGraphQL(in *model.Auth) (*externalschema.Auth, error)
	InputFromGraphQL(in *externalschema.AuthInput) (*model.AuthInput, error)
}

type converter struct {
	authConverter AuthConverter
}

func NewConverter(authConverter AuthConverter) *converter {
	return &converter{
		authConverter: authConverter,
	}
}

func (c *converter) ToGraphQL(in *model.PackageInstanceAuth) (*externalschema.PackageInstanceAuth, error) {
	if in == nil {
		return nil, nil
	}

	auth, err := c.authConverter.ToGraphQL(in.Auth)
	if err != nil {
		return nil, errors.Wrap(err, "while converting Auth to GraphQL")
	}

	return &externalschema.PackageInstanceAuth{
		ID:          in.ID,
		Context:     c.strPtrToJSONPtr(in.Context),
		InputParams: c.strPtrToJSONPtr(in.InputParams),
		Auth:        auth,
		Status:      c.statusToGraphQL(in.Status),
	}, nil
}

func (c *converter) MultipleToGraphQL(in []*model.PackageInstanceAuth) ([]*externalschema.PackageInstanceAuth, error) {
	var packageInstanceAuths []*externalschema.PackageInstanceAuth
	for _, r := range in {
		if r == nil {
			continue
		}
		pia, err := c.ToGraphQL(r)
		if err != nil {
			return nil, err
		}
		packageInstanceAuths = append(packageInstanceAuths, pia)
	}

	return packageInstanceAuths, nil
}

func (c *converter) RequestInputFromGraphQL(in externalschema.PackageInstanceAuthRequestInput) model.PackageInstanceAuthRequestInput {
	return model.PackageInstanceAuthRequestInput{
		Context:     c.jsonPtrToStrPtr(in.Context),
		InputParams: c.jsonPtrToStrPtr(in.InputParams),
	}
}

func (c *converter) SetInputFromGraphQL(in externalschema.PackageInstanceAuthSetInput) (model.PackageInstanceAuthSetInput, error) {
	auth, err := c.authConverter.InputFromGraphQL(in.Auth)
	if err != nil {
		return model.PackageInstanceAuthSetInput{}, errors.Wrap(err, "while converting Auth")
	}

	out := model.PackageInstanceAuthSetInput{
		Auth: auth,
	}

	if in.Status != nil {
		out.Status = &model.PackageInstanceAuthStatusInput{
			Condition: model.PackageInstanceAuthSetStatusConditionInput(in.Status.Condition),
			Message:   in.Status.Message,
			Reason:    in.Status.Reason,
		}
	}

	return out, nil
}

func (c *converter) ToEntity(in model.PackageInstanceAuth) (Entity, error) {
	out := Entity{
		ID:          in.ID,
		PackageID:   in.PackageID,
		TenantID:    in.Tenant,
		Context:     repo.NewNullableString(in.Context),
		InputParams: repo.NewNullableString(in.InputParams),
	}
	authValue, err := c.nullStringFromAuthPtr(in.Auth)
	if err != nil {
		return Entity{}, err
	}
	out.AuthValue = authValue

	if in.Status != nil {
		out.StatusCondition = string(in.Status.Condition)
		out.StatusTimestamp = in.Status.Timestamp
		out.StatusMessage = in.Status.Message
		out.StatusReason = in.Status.Reason
	}

	return out, nil
}

func (c *converter) FromEntity(in Entity) (model.PackageInstanceAuth, error) {
	auth, err := c.authPtrFromNullString(in.AuthValue)
	if err != nil {
		return model.PackageInstanceAuth{}, err
	}

	return model.PackageInstanceAuth{
		ID:          in.ID,
		PackageID:   in.PackageID,
		Tenant:      in.TenantID,
		Context:     repo.StringPtrFromNullableString(in.Context),
		InputParams: repo.StringPtrFromNullableString(in.InputParams),
		Auth:        auth,
		Status: &model.PackageInstanceAuthStatus{
			Condition: model.PackageInstanceAuthStatusCondition(in.StatusCondition),
			Timestamp: in.StatusTimestamp,
			Message:   in.StatusMessage,
			Reason:    in.StatusReason,
		},
	}, nil
}

func (c *converter) statusToGraphQL(in *model.PackageInstanceAuthStatus) *externalschema.PackageInstanceAuthStatus {
	if in == nil {
		return nil
	}

	return &externalschema.PackageInstanceAuthStatus{
		Condition: externalschema.PackageInstanceAuthStatusCondition(in.Condition),
		Timestamp: externalschema.Timestamp(in.Timestamp),
		Message:   in.Message,
		Reason:    in.Reason,
	}
}

func (c *converter) strPtrToJSONPtr(in *string) *externalschema.JSON {
	if in == nil {
		return nil
	}
	out := externalschema.JSON(*in)
	return &out
}

func (c *converter) jsonPtrToStrPtr(in *externalschema.JSON) *string {
	if in == nil {
		return nil
	}
	out := string(*in)
	return &out
}

func (c *converter) nullStringFromAuthPtr(in *model.Auth) (sql.NullString, error) {
	if in == nil {
		return sql.NullString{}, nil
	}
	valueMarshalled, err := json.Marshal(*in)
	if err != nil {
		return sql.NullString{}, errors.Wrap(err, "while marshalling Auth")
	}
	return sql.NullString{
		String: string(valueMarshalled),
		Valid:  true,
	}, nil
}

func (c *converter) authPtrFromNullString(in sql.NullString) (*model.Auth, error) {
	if !in.Valid {
		return nil, nil
	}
	var auth model.Auth
	err := json.Unmarshal([]byte(in.String), &auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
