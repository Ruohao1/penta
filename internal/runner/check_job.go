package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/checks"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/sinks"
)

type CheckJob struct {
	StageName string
	HostKey   string

	Checker checks.Checker
	Input   any

	Sink sinks.Sink
}

func (j CheckJob) Key() string { return j.HostKey }

func (j CheckJob) Run(ctx context.Context) error {
	emit := func(x any) {
		switch v := x.(type) {
		case model.Finding:
			ev := model.NewFindingEvent(&v)
			ev.Stage = j.StageName
			j.Sink.Emit(ctx, ev)

		default:
			j.Sink.Emit(ctx, model.Event{
				EmittedAt: time.Now().UTC(),
				Type:      model.EventUnknown,
				Stage:     j.StageName,
				Err:       fmt.Sprintf("unknown event type: %T", x),
			})
		}
	}
	return j.Checker.Check()(ctx, j.Input, emit)
}
