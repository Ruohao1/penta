package engine

import (
	"context"
	"fmt"

	"github.com/Ruohao1/penta/internal/engine/external"
	"github.com/Ruohao1/penta/internal/hosts"
	"github.com/Ruohao1/penta/internal/model"
)

type Engine struct {
	Opts *model.RunOptions
}

func New(opts *model.RunOptions) *Engine {
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
			err = eng.runHosts(ctx, req, bus.Emit)
		case model.ModePorts:
			err = eng.runPorts(ctx, req, bus.Emit)
		case model.ModeWeb:
			err = eng.runWeb(ctx, req, bus.Emit)
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

func (eng Engine) runHosts(ctx context.Context, req model.Request, emit func(context.Context, model.Event)) error {
	eng.Opts.Compile(req)
	switch req.Backend {
	case model.BackendInternal:
		return hosts.Run(ctx, *eng.Opts, emit)
	case model.BackendNmap:
		return external.RunNmap(ctx, req, emit)
	default:
		return fmt.Errorf("unknown backend %q", req.Backend)
	}
}

func (eng Engine) runPorts(ctx context.Context, req model.Request, emit func(context.Context, model.Event)) error {
	return nil
}

func (eng Engine) runWeb(ctx context.Context, req model.Request, emit func(context.Context, model.Event)) error {
	return nil
}
