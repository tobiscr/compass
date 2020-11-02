package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
)

type RuntimeLister interface {
	List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error)
}

type ApplicationLister interface {
	ListByRuntimeID(ctx context.Context, runtimeID uuid.UUID, pageSize int, cursor string) (*model.ApplicationPage, error)
}

type RuntimeGetter interface {
	GetByID(ctx context.Context, tenant, id string) (*model.Runtime, error)
}
