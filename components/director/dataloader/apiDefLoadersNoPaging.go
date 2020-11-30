//go:generate go run github.com/vektah/dataloaden ApiDefLoaderNoPaging ParamApiDefNoPaging []*github.com/kyma-incubator/compass/components/director/pkg/graphql.APIDefinition

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const loadersKeyApiDefNoPaging = "dataloadersApiDefNoPaging"

type ApiDefLoadersNoPaging struct {
	ApiDefByIdNoPaging ApiDefLoaderNoPaging
}

type ParamApiDefNoPaging struct {
	ID  string
	Ctx context.Context
}

func HandlerApiDefNoPaging(fetchFunc func(keys []ParamApiDefNoPaging) ([][]*graphql.APIDefinition, []error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKeyApiDefNoPaging, &ApiDefLoadersNoPaging{
				ApiDefByIdNoPaging: ApiDefLoaderNoPaging{
					maxBatch: 100,
					wait:     1 * time.Millisecond,
					fetch:    fetchFunc,
				},
			})
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ApiDefForNoPaging(ctx context.Context) *ApiDefLoadersNoPaging {
	return ctx.Value(loadersKeyApiDefNoPaging).(*ApiDefLoadersNoPaging)
}
