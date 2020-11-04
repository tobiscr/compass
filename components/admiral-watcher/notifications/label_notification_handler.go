package notifications

import (
	"context"
	"encoding/json"
	"github.com/kyma-incubator/compass/components/admiral-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	"strings"
)

type Label struct {
	ID        string   `db:"id" json:"id"`
	TenantID  string   `db:"tenant_id" json:"tenant_id"`
	Key       string   `db:"key" json:"key"`
	AppID     string   `db:"app_id" json:"app_id"`
	RuntimeID string   `db:"runtime_id" json:"runtime_id"`
	Value     []string `db:"value" json:"value"`
}

type LabelNotificationHandler struct {
	Handlers map[resource.Type]NotificationLabelHandler
}

func (l *LabelNotificationHandler) HandleCreate(ctx context.Context, data []byte) error {
	label := Label{}
	if err := json.Unmarshal(data, &label); err != nil {
		return errors.Errorf("could not unmarshal label: %s", err)
	}

	if !strings.Contains(strings.ToLower(label.Key), "scenario") {
		log.C(ctx).Warnf("handling events for creation of labels with key %s is noop", label.Key)
		return nil
	}

	if len(label.Value) == 1 && label.Value[0] == "DEFAULT" {
		log.C(ctx).Warnf("handling events for creation of labels with key %s and single value DEFAULT is noop", label.Key)
		return nil
	}

	if len(label.AppID) != 0 {
		handler, ok := l.Handlers[resource.Application]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleCreate(ctx, label); err != nil {
			return err
		}
	} else if len(label.RuntimeID) != 0 {
		handler, ok := l.Handlers[resource.Runtime]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleCreate(ctx, label); err != nil {
			return err
		}
	} else {
		log.C(ctx).Infof("label %v does not belong to runtimes or applications", label)
	}

	log.C(ctx).Infof("Successfully handled create event for label %v", label)
	return nil
}

func (l *LabelNotificationHandler) HandleUpdate(ctx context.Context, data []byte) error {
	label := Label{}
	if err := json.Unmarshal(data, &label); err != nil {
		return errors.Errorf("could not unmarshal label: %s", err)
	}

	if !strings.Contains(strings.ToLower(label.Key), "scenario") {
		log.C(ctx).Warnf("handling events for creation of labels with key %s is noop", label.Key)
		return nil
	}

	if len(label.AppID) != 0 {
		handler, ok := l.Handlers[resource.Application]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleUpdate(ctx, label); err != nil {
			return err
		}
	} else if len(label.RuntimeID) != 0 {
		handler, ok := l.Handlers[resource.Runtime]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleUpdate(ctx, label); err != nil {
			return err
		}
	} else {
		log.C(ctx).Infof("label %v does not belong to runtimes or applications", label)
	}

	log.C(ctx).Infof("Successfully handled update event for label %v", label)
	return nil
}

func (l *LabelNotificationHandler) HandleDelete(ctx context.Context, data []byte) error {
	label := Label{}
	if err := json.Unmarshal(data, &label); err != nil {
		return errors.Errorf("could not unmarshal label: %s", err)
	}

	if !strings.Contains(strings.ToLower(label.Key), "scenario") {
		log.C(ctx).Warnf("handling events for creation of labels with key %s is noop", label.Key)
		return nil
	}

	if len(label.Value) == 1 && label.Value[0] == "DEFAULT" {
		log.C(ctx).Warnf("handling events for creation of labels with key %s and single value DEFAULT is noop", label.Key)
		return nil
	}

	if len(label.AppID) != 0 {
		handler, ok := l.Handlers[resource.Application]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleDelete(ctx, label); err != nil {
			return err
		}
	} else if len(label.RuntimeID) != 0 {
		handler, ok := l.Handlers[resource.Runtime]
		if !ok {
			return errors.New("handler for applications label creation not found")
		}

		if err := handler.HandleDelete(ctx, label); err != nil {
			return err
		}
	} else {
		log.C(ctx).Infof("label %v does not belong to runtimes or applications", label)
	}

	log.C(ctx).Infof("Successfully handled delete event for label %v", label)
	return nil
}
