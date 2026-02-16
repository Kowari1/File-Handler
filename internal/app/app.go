package app

import (
	"context"
	"sync"
)

type Scanner interface {
	Scan(ctx context.Context, out chan<- string) error
}

type WorkerPool interface {
	Start(ctx context.Context)
	Wait()
}

type App struct {
	scanner    Scanner
	workerPool WorkerPool
}

func New(scanner Scanner, workerPool WorkerPool) *App {
	return &App{
		scanner:    scanner,
		workerPool: workerPool,
	}
}

func (a *App) Run(ctx context.Context) error {
	jobs := make(chan string)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.workerPool.Start(ctx)
	}()

	if err := a.scanner.Scan(ctx, jobs); err != nil {
		return err
	}

	close(jobs)

	wg.Wait()

	return nil
}
