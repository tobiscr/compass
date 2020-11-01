package storage

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
)

type RuntimeLister interface {
	List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (*model.RuntimePage, error)
}
type AppLabelNotificationHandler struct {
}

func (a *AppLabelNotificationHandler) HandleCreate(ctx context.Context, label Label) error {
	panic("implement me")
}

func (a *AppLabelNotificationHandler) HandleUpdate(ctx context.Context, label Label) error {
	panic("implement me")
}

func (a *AppLabelNotificationHandler) HandleDelete(ctx context.Context, label Label) error {
	panic("implement me")
}
