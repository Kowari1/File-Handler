package worker

import (
	"context"
	"sync"
)

type JobHandler interface {
	Handle(ctx context.Context, job string) error
}

type Pool struct {
	workerCount int
	handler     JobHandler
	jobs        <-chan string
	wg          sync.WaitGroup
}

func New(
	workerCount int,
	handler JobHandler,
	jobs <-chan string,
) *Pool {
	return &Pool{
		workerCount: workerCount,
		handler:     handler,
		jobs:        jobs,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-p.jobs:
			if !ok {
				return
			}

			if err := p.handler.Handle(ctx, job); err != nil {
				continue
			}
		}
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
