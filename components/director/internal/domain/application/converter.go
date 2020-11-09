package application

import (
	"encoding/json"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	log "github.com/sirupsen/logrus"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"
)

type converter struct {
	webhook WebhookConverter

	pkg PackageConverter
}

func NewConverter(webhook WebhookConverter, pkgConverter PackageConverter) *converter {
	return &converter{webhook: webhook, pkg: pkgConverter}
}

func (c *converter) ToEntity(in *model.Application) (*Entity, error) {
	if in == nil {
		return nil, nil
	}

	if in.Status == nil {
		return nil, apperrors.NewInternalError("invalid input model")
	}

	return &Entity{
		ID:                  in.ID,
		TenantID:            in.Tenant,
		Name:                in.Name,
		Description:         repo.NewNullableString(in.Description),
		StatusCondition:     string(in.Status.Condition),
		StatusTimestamp:     in.Status.Timestamp,
		HealthCheckURL:      repo.NewNullableString(in.HealthCheckURL),
		IntegrationSystemID: repo.NewNullableString(in.IntegrationSystemID),
		ProviderName:        repo.NewNullableString(in.ProviderName),
	}, nil
}

func (c *converter) FromEntity(entity *Entity) *model.Application {
	if entity == nil {
		return nil
	}

	return &model.Application{
		ID:          entity.ID,
		Tenant:      entity.TenantID,
		Name:        entity.Name,
		Description: repo.StringPtrFromNullableString(entity.Description),
		Status: &model.ApplicationStatus{
			Condition: model.ApplicationStatusCondition(entity.StatusCondition),
			Timestamp: entity.StatusTimestamp,
		},
		IntegrationSystemID: repo.StringPtrFromNullableString(entity.IntegrationSystemID),
		HealthCheckURL:      repo.StringPtrFromNullableString(entity.HealthCheckURL),
		ProviderName:        repo.StringPtrFromNullableString(entity.ProviderName),
	}
}

func (c *converter) ToGraphQL(in *model.Application) *externalschema.Application {
	if in == nil {
		return nil
	}

	return &externalschema.Application{
		ID:                  in.ID,
		Status:              c.statusToGraphQL(in.Status),
		Name:                in.Name,
		Description:         in.Description,
		HealthCheckURL:      in.HealthCheckURL,
		IntegrationSystemID: in.IntegrationSystemID,
		ProviderName:        in.ProviderName,
	}
}

func (c *converter) MultipleToGraphQL(in []*model.Application) []*externalschema.Application {
	var runtimes []*externalschema.Application
	for _, r := range in {
		if r == nil {
			continue
		}

		runtimes = append(runtimes, c.ToGraphQL(r))
	}

	return runtimes
}

func (c *converter) CreateInputFromGraphQL(in externalschema.ApplicationRegisterInput) (model.ApplicationRegisterInput, error) {
	var labels map[string]interface{}
	if in.Labels != nil {
		labels = *in.Labels
	}

	log.Debugf("Converting Webhooks from Application registration GraphQL input with name %s", in.Name)
	webhooks, err := c.webhook.MultipleInputFromGraphQL(in.Webhooks)
	if err != nil {
		return model.ApplicationRegisterInput{}, errors.Wrap(err, "while converting Webhooks")
	}

	log.Debugf("Converting Packages from Application registration GraphQL input with name %s", in.Name)
	packages, err := c.pkg.MultipleCreateInputFromGraphQL(in.Packages)
	if err != nil {
		return model.ApplicationRegisterInput{}, errors.Wrap(err, "while converting Packages")
	}

	return model.ApplicationRegisterInput{
		Name:                in.Name,
		Description:         in.Description,
		Labels:              labels,
		HealthCheckURL:      in.HealthCheckURL,
		IntegrationSystemID: in.IntegrationSystemID,
		StatusCondition:     c.statusConditionToModel(in.StatusCondition),
		ProviderName:        in.ProviderName,
		Webhooks:            webhooks,
		Packages:            packages,
	}, nil
}

func (c *converter) UpdateInputFromGraphQL(in externalschema.ApplicationUpdateInput) model.ApplicationUpdateInput {
	return model.ApplicationUpdateInput{
		Description:         in.Description,
		HealthCheckURL:      in.HealthCheckURL,
		IntegrationSystemID: in.IntegrationSystemID,
		ProviderName:        in.ProviderName,
		StatusCondition:     c.statusConditionToModel(in.StatusCondition),
	}
}

func (c *converter) CreateInputJSONToGQL(in string) (externalschema.ApplicationRegisterInput, error) {
	var appInput externalschema.ApplicationRegisterInput
	err := json.Unmarshal([]byte(in), &appInput)
	if err != nil {
		return externalschema.ApplicationRegisterInput{}, errors.Wrap(err, "while unmarshalling string to ApplicationRegisterInput")
	}

	return appInput, nil
}

func (c *converter) CreateInputGQLToJSON(in *externalschema.ApplicationRegisterInput) (string, error) {
	appInput, err := json.Marshal(in)
	if err != nil {
		return "", errors.Wrap(err, "while marshaling application input")
	}

	return string(appInput), nil
}

func (c *converter) GraphQLToModel(obj *externalschema.Application, tenantID string) *model.Application {
	if obj == nil {
		return nil
	}

	return &model.Application{
		ID:                  obj.ID,
		ProviderName:        obj.ProviderName,
		Tenant:              tenantID,
		Name:                obj.Name,
		Description:         obj.Description,
		Status:              c.statusToModel(obj.Status),
		HealthCheckURL:      obj.HealthCheckURL,
		IntegrationSystemID: obj.IntegrationSystemID,
	}
}

func (c *converter) statusToGraphQL(in *model.ApplicationStatus) *externalschema.ApplicationStatus {
	if in == nil {
		return &externalschema.ApplicationStatus{Condition: externalschema.ApplicationStatusConditionInitial}
	}

	var condition externalschema.ApplicationStatusCondition

	switch in.Condition {
	case model.ApplicationStatusConditionInitial:
		condition = externalschema.ApplicationStatusConditionInitial
	case model.ApplicationStatusConditionFailed:
		condition = externalschema.ApplicationStatusConditionFailed
	case model.ApplicationStatusConditionConnected:
		condition = externalschema.ApplicationStatusConditionConnected
	default:
		condition = externalschema.ApplicationStatusConditionInitial
	}

	return &externalschema.ApplicationStatus{
		Condition: condition,
		Timestamp: externalschema.Timestamp(in.Timestamp),
	}
}

func (c *converter) statusToModel(in *externalschema.ApplicationStatus) *model.ApplicationStatus {
	if in == nil {
		return &model.ApplicationStatus{Condition: model.ApplicationStatusConditionInitial}
	}

	var condition model.ApplicationStatusCondition

	switch in.Condition {
	case externalschema.ApplicationStatusConditionInitial:
		condition = model.ApplicationStatusConditionInitial
	case externalschema.ApplicationStatusConditionFailed:
		condition = model.ApplicationStatusConditionFailed
	case externalschema.ApplicationStatusConditionConnected:
		condition = model.ApplicationStatusConditionConnected
	default:
		condition = model.ApplicationStatusConditionInitial
	}
	return &model.ApplicationStatus{
		Condition: condition,
		Timestamp: time.Time(in.Timestamp),
	}
}

func (c *converter) statusConditionToModel(in *externalschema.ApplicationStatusCondition) *model.ApplicationStatusCondition {
	if in == nil {
		return nil
	}

	var condition model.ApplicationStatusCondition
	switch *in {
	case externalschema.ApplicationStatusConditionConnected:
		condition = model.ApplicationStatusConditionConnected
	case externalschema.ApplicationStatusConditionFailed:
		condition = model.ApplicationStatusConditionFailed
	case externalschema.ApplicationStatusConditionInitial:
		fallthrough
	default:
		condition = model.ApplicationStatusConditionInitial
	}

	return &condition
}
