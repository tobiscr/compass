package scheduler

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

type Worker struct {
	wg             *sync.WaitGroup
	maxConcurrency int
	errChan        chan error
	ctx            context.Context
	cancelFunc     context.CancelFunc
	work           chan func(context.Context) error
}

func NewWorker(ctx context.Context, maxConcurrency int, work chan func(context.Context) error) *Worker {
	childContext, cancel := context.WithCancel(ctx)

	return &Worker{wg: &sync.WaitGroup{},
		maxConcurrency: maxConcurrency,
		ctx:            childContext,
		cancelFunc:     cancel,
		errChan:        make(chan error),
		work:           work,
	}
}

func (s Worker) Start() {
	for i := 0; i < s.maxConcurrency; i++ {
		s.wg.Add(1)
		go func() {
			defer func() {
				s.wg.Done()
			}()
			for w := range s.work {
				err := w(s.ctx)
				if err != nil {
					select {
					case <-s.ctx.Done():
						return
					case s.errChan <- err:
						return
					}
				}
			}
		}()
	}
}

func (s Worker) Wait() error {
	success := make(chan interface{})
	go func() {
		s.wg.Wait()
		close(success)
	}()

	select {
	case <-success:
		return nil
	case err := <-s.errChan:
		s.cancelFunc()
		return errors.Wrap(err, "while working on job")
	}
}
