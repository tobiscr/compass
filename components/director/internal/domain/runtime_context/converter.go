package runtime_context

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

type converter struct{}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGraphQL(in *model.RuntimeContext) *externalschema.RuntimeContext {
	if in == nil {
		return nil
	}

	return &externalschema.RuntimeContext{
		ID:    in.ID,
		Key:   in.Key,
		Value: in.Value,
	}
}

func (c *converter) MultipleToGraphQL(in []*model.RuntimeContext) []*externalschema.RuntimeContext {
	var runtimeContexts []*externalschema.RuntimeContext
	for _, r := range in {
		if r == nil {
			continue
		}

		runtimeContexts = append(runtimeContexts, c.ToGraphQL(r))
	}

	return runtimeContexts
}

func (c *converter) InputFromGraphQL(in externalschema.RuntimeContextInput, runtimeID string) model.RuntimeContextInput {
	var labels map[string]interface{}
	if in.Labels != nil {
		labels = *in.Labels
	}

	return model.RuntimeContextInput{
		Key:       in.Key,
		Value:     in.Value,
		RuntimeID: runtimeID,
		Labels:    labels,
	}
}
