package notifications

import (
	"context"
	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
	"github.com/lib/pq"
)

type RuntimeLister interface {
	List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error)
}

type ApplicationLister interface {
	ListByRuntimeID(ctx context.Context, runtimeID uuid.UUID, pageSize int, cursor string) (*model.ApplicationPage, error)
}

type ApplicationLabelGetter interface {
	GetLabel(ctx context.Context, applicationID string, key string) (*model.Label, error)
}

type RuntimeLabelGetter interface {
	GetLabel(ctx context.Context, runtimeID string, key string) (*model.Label, error)
}

type RuntimeGetter interface {
	Get(ctx context.Context, id string) (*model.Runtime, error)
}

type NotificationLabelHandler interface {
	HandleCreate(ctx context.Context, label Label) error
	HandleUpdate(ctx context.Context, label Label) error
	HandleDelete(ctx context.Context, label Label) error
}

type NotificationHandler interface {
	HandleCreate(ctx context.Context, data []byte) error
	HandleUpdate(ctx context.Context, data []byte) error
	HandleDelete(ctx context.Context, data []byte) error
}

type NotificationListener interface {
	Listen(channel string) error
	Ping() error
	Close() error
	NotificationChannel() <-chan *pq.Notification
}
