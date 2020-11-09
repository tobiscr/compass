package appdetails

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql/externalschema"
	"github.com/pkg/errors"
)

type AppDetailsContextKey struct{}

var NoAppDetailsError = errors.New("cannot read Application details from context")
var NilContextError = errors.New("context is empty")

func LoadFromContext(ctx context.Context) (externalschema.ApplicationExt, error) {
	if ctx == nil {
		return externalschema.ApplicationExt{}, NilContextError
	}

	value := ctx.Value(AppDetailsContextKey{})

	appDetails, ok := value.(externalschema.ApplicationExt)

	if !ok {
		return externalschema.ApplicationExt{}, NoAppDetailsError
	}

	return appDetails, nil
}

func SaveToContext(ctx context.Context, appDetails externalschema.ApplicationExt) context.Context {
	return context.WithValue(ctx, AppDetailsContextKey{}, appDetails)
}
