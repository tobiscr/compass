package viewer

import (
	"github.com/kyma-incubator/compass/components/director/internal/consumer"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
)

func ToViewer(cons consumer.Consumer) (*externalschema.Viewer, error) {
	switch cons.ConsumerType {
	case consumer.Runtime:
		return &externalschema.Viewer{ID: cons.ConsumerID, Type: externalschema.ViewerTypeRuntime}, nil
	case consumer.Application:
		return &externalschema.Viewer{ID: cons.ConsumerID, Type: externalschema.ViewerTypeApplication}, nil
	case consumer.IntegrationSystem:
		return &externalschema.Viewer{ID: cons.ConsumerID, Type: externalschema.ViewerTypeIntegrationSystem}, nil
	case consumer.User:
		return &externalschema.Viewer{ID: cons.ConsumerID, Type: externalschema.ViewerTypeUser}, nil
	}

	return nil, apperrors.NewInternalError("viewer does not exist")

}
