package runner

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type WorkFn[T any] func(context.Context, T) error

func RunPool[T any](ctx context.Context, items []T, workers int, limiter *rate.Limiter, fn WorkFn[T]) error {
	if workers < 1 {
		workers = 1
	}

	ch := make(chan T)
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for it := range ch {
			if limiter != nil {
				if err := limiter.Wait(ctx); err != nil {
					select {
					case errCh <- err:
					default:
					}
					return
				}
			}
			if err := fn(ctx, it); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		defer close(ch)
		for _, it := range items {
			select {
			case <-ctx.Done():
				return
			case ch <- it:
			}
		}
	}()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-done:
		return nil
	case err := <-errCh:
		return err
	}
}
