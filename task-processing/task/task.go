package task

import (
	"context"
	"sync"
)

type Worker interface {
	Work(ctx context.Context)
}

type Task struct {
	ctx context.Context
	work chan Worker
	wg sync.WaitGroup
}

func New(ctx context.Context, maxThreads int) *Task {
	t := Task{
		work: make(chan Worker),
		ctx: ctx,
	}
	t.wg.Add(maxThreads)
	for i := 0; i < maxThreads; i++ {
		go func() {
			for w := range t.work {
				w.Work(ctx)
			}
			t.wg.Done()
		}()
	}
	return &t
}

func (t *Task) Shutdown() {
	close(t.work)
	t.wg.Wait()
}

func (t *Task) Do(w Worker) {
	t.work <- w
}