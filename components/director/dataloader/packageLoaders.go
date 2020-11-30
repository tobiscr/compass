//go:generate go run github.com/vektah/dataloaden PackageLoader Param *github.com/kyma-incubator/compass/components/director/pkg/graphql.PackagePage

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const loadersKey = "dataloaders"

type Loaders struct {
	PkgById PackageLoader
}

type Param struct {
	ID    string
	First *int
	After *graphql.PageCursor
	Ctx   context.Context
}

func Handler(fetchFunc func(keys []Param) ([]*graphql.PackagePage, []error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
				PkgById: PackageLoader{
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

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
