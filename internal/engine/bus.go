package engine

import (
	"context"
	"sync"
	"time"

	"github.com/Ruohao1/penta/internal/model"
)

type Bus struct {
	out chan model.Event
	wg  sync.WaitGroup
}

func NewBus(buffer int) *Bus {
	out := make(chan model.Event, buffer)
	return &Bus{out: out}
}

func (b *Bus) Out() <-chan model.Event { return b.out }

func (b *Bus) Go(fn func()) {
	b.wg.Go(fn)
}

func (b *Bus) Emit(ctx context.Context, e model.Event) {
	if e.Time.IsZero() {
		e.Time = time.Now().UTC()
	}
	select {
	case <-ctx.Done():
	case b.out <- e:
	}
}

func (b *Bus) WaitClose() {
	b.wg.Wait()
	close(b.out)
}
