//go:generate go run github.com/vektah/dataloaden PackageLoaderNoPaging ParamNoPaging []*github.com/kyma-incubator/compass/components/director/pkg/graphql.Package

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const loadersKeyNoPaging = "dataloadersNoPaging"

type LoadersNoPaging struct {
	PkgByIdNoPaging PackageLoaderNoPaging
}

type ParamNoPaging struct {
	ID  string
	Ctx context.Context
}

func HandlerNoPaging(fetchFunc func(keys []ParamNoPaging) ([][]*graphql.Package, []error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKeyNoPaging, &LoadersNoPaging{
				PkgByIdNoPaging: PackageLoaderNoPaging{
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

func ForNoPaging(ctx context.Context) *LoadersNoPaging {
	return ctx.Value(loadersKeyNoPaging).(*LoadersNoPaging)
}
