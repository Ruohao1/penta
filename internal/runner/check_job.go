package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/checks"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/ui"
)

type CheckJob struct {
	StageName string
	HostKey   string // for per-host gate

	Checker checks.Checker
	Input   any

	Sink ui.EventSink
}

func (j CheckJob) Key() string { return j.HostKey }

func (j CheckJob) Run(ctx context.Context) error {
	emit := func(x any) {
		switch v := x.(type) {
		case model.Event:
			switch t := v.Type; t {
			case model.EventFinding:
				j.Sink.Emit(ctx, model.Event{
					EmittedAt: time.Now().UTC(),
					Type:      model.EventFinding,
					Stage:     j.StageName,
					Finding:   v.Finding,
				})
			}
		case model.HostStateEvent:
			j.Sink.Emit(ctx, model.Event{
				EmittedAt: time.Now().UTC(),
				Type:      model.EventHostState,
				Stage:     j.StageName,
				HostState: &v,
			})
		default:
			j.Sink.Emit(ctx, model.Event{
				EmittedAt: time.Now().UTC(),
				Type:      model.EventFinding,
				Stage:     j.StageName,
				Err:       fmt.Sprintf("unknown event type: %T", x),
			})
		}
	}
	return j.Checker.Check()(ctx, j.Input, emit)
}
