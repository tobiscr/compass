//go:generate go run github.com/vektah/dataloaden EventDefLoaderNoPaging ParamEventDefNoPaging []*github.com/kyma-incubator/compass/components/director/pkg/graphql.EventDefinition

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

const loadersKeyEventDefNoPaging = "dataloadersEventDefNoPaging"

type EventDefLoadersNoPaging struct {
	EventDefByIdNoPaging EventDefLoaderNoPaging
}

type ParamEventDefNoPaging struct {
	ID  string
	Ctx context.Context
}

func HandlerEventDefNoPaging(fetchFunc func(keys []ParamEventDefNoPaging) ([][]*graphql.EventDefinition, []error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKeyEventDefNoPaging, &EventDefLoadersNoPaging{
				EventDefByIdNoPaging: EventDefLoaderNoPaging{
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

func EventDefForNoPaging(ctx context.Context) *EventDefLoadersNoPaging {
	return ctx.Value(loadersKeyEventDefNoPaging).(*EventDefLoadersNoPaging)
}
