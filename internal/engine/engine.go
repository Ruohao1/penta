package engine

import (
	"context"
	"fmt"

	"github.com/Ruohao1/penta/internal/engine/external"
	"github.com/Ruohao1/penta/internal/hosts"
	"github.com/Ruohao1/penta/internal/model"
)

type Engine struct {
	Opts *RunOptions
}

func New(opts *RunOptions) *Engine {
	return &Engine{Opts: opts}
}

func (eng Engine) Run(ctx context.Context, req model.Request) <-chan model.Event {
	bus := NewBus(4096)

	bus.Go(func() {
		bus.Emit(ctx, model.NewEvent(model.EventEngineStart))
		defer bus.Emit(ctx, model.NewEvent(model.EventEngineDone))

		var err error
		switch req.Mode {
		case model.ModeHosts:
			err = eng.runHosts(ctx, bus.Emit, req)
		case model.ModePorts:
			err = eng.runPorts(ctx, bus.Emit, req)
		case model.ModeWeb:
			err = eng.runWeb(ctx, bus.Emit, req)
		default:
			err = fmt.Errorf("unknown mode %q", req.Mode)
		}

		if err != nil {
			e := model.NewEvent(model.EventError)
			e.Err = err.Error()
			bus.Emit(ctx, e)
		}
	})

	go bus.WaitClose()
	return bus.Out()
}

func (eng Engine) runHosts(ctx context.Context, emit func(context.Context, model.Event), req model.Request) error {
	switch req.Backend {
	case model.BackendInternal:
		return hosts.Run(ctx, req, emit)
	case model.BackendNmap:
		return external.RunNmap(ctx, req, emit)
	default:
		return fmt.Errorf("unknown backend %q", req.Backend)
	}
}

func (eng Engine) runPorts(ctx context.Context, emit func(context.Context, model.Event), req model.Request) error {
	return nil
}

func (eng Engine) runWeb(ctx context.Context, emit func(context.Context, model.Event), req model.Request) error {
	return nil
}
