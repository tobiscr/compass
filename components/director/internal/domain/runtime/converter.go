package runtime

import (
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

type converter struct{}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGraphQL(in *model.Runtime) *externalschema.Runtime {
	if in == nil {
		return nil
	}

	return &externalschema.Runtime{
		ID:          in.ID,
		Status:      c.statusToGraphQL(in.Status),
		Name:        in.Name,
		Description: in.Description,
		Metadata:    c.metadataToGraphQL(in.CreationTimestamp),
	}
}

func (c *converter) MultipleToGraphQL(in []*model.Runtime) []*externalschema.Runtime {
	var runtimes []*externalschema.Runtime
	for _, r := range in {
		if r == nil {
			continue
		}

		runtimes = append(runtimes, c.ToGraphQL(r))
	}

	return runtimes
}

func (c *converter) InputFromGraphQL(in externalschema.RuntimeInput) model.RuntimeInput {
	var labels map[string]interface{}
	if in.Labels != nil {
		labels = *in.Labels
	}

	return model.RuntimeInput{
		Name:            in.Name,
		Description:     in.Description,
		Labels:          labels,
		StatusCondition: c.statusConditionToModel(in.StatusCondition),
	}
}

func (c *converter) statusToGraphQL(in *model.RuntimeStatus) *externalschema.RuntimeStatus {
	if in == nil {
		return &externalschema.RuntimeStatus{
			Condition: externalschema.RuntimeStatusConditionInitial,
		}
	}

	var condition externalschema.RuntimeStatusCondition

	switch in.Condition {
	case model.RuntimeStatusConditionInitial:
		condition = externalschema.RuntimeStatusConditionInitial
	case model.RuntimeStatusConditionProvisioning:
		condition = externalschema.RuntimeStatusConditionProvisioning
	case model.RuntimeStatusConditionFailed:
		condition = externalschema.RuntimeStatusConditionFailed
	case model.RuntimeStatusConditionConnected:
		condition = externalschema.RuntimeStatusConditionConnected
	default:
		condition = externalschema.RuntimeStatusConditionInitial
	}

	return &externalschema.RuntimeStatus{
		Condition: condition,
		Timestamp: externalschema.Timestamp(in.Timestamp),
	}
}

func (c *converter) metadataToGraphQL(creationTimestamp time.Time) *externalschema.RuntimeMetadata {
	return &externalschema.RuntimeMetadata{
		CreationTimestamp: externalschema.Timestamp(creationTimestamp),
	}
}

func (c *converter) statusConditionToModel(in *externalschema.RuntimeStatusCondition) *model.RuntimeStatusCondition {
	if in == nil {
		return nil
	}

	var condition model.RuntimeStatusCondition
	switch *in {
	case externalschema.RuntimeStatusConditionConnected:
		condition = model.RuntimeStatusConditionConnected
	case externalschema.RuntimeStatusConditionFailed:
		condition = model.RuntimeStatusConditionFailed
	case externalschema.RuntimeStatusConditionProvisioning:
		condition = model.RuntimeStatusConditionProvisioning
	case externalschema.RuntimeStatusConditionInitial:
		fallthrough
	default:
		condition = model.RuntimeStatusConditionInitial
	}

	return &condition
}
