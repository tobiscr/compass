package eventing

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

func RuntimeEventingConfigurationToGraphQL(in *model.RuntimeEventingConfiguration) *externalschema.RuntimeEventingConfiguration {
	if in == nil {
		return nil
	}

	return &externalschema.RuntimeEventingConfiguration{
		DefaultURL: in.DefaultURL.String(),
	}
}

func ApplicationEventingConfigurationToGraphQL(in *model.ApplicationEventingConfiguration) *externalschema.ApplicationEventingConfiguration {
	if in == nil {
		return nil
	}

	return &externalschema.ApplicationEventingConfiguration{
		DefaultURL: in.DefaultURL.String(),
	}
}
