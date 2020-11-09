package viewer

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"

	"github.com/kyma-incubator/compass/components/director/internal/consumer"
	"github.com/pkg/errors"
)

type Resolver struct {
}

func NewViewerResolver() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Viewer(ctx context.Context) (*externalschema.Viewer, error) {
	cons, err := consumer.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "while getting viewer from context")
	}

	return ToViewer(cons)
}
