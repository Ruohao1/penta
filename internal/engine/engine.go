// Package engine provides a generic engine for running checks
package engine

import (
	"context"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/runner"
	"github.com/Ruohao1/penta/internal/sinks"
	"github.com/Ruohao1/penta/internal/stages"
	"golang.org/x/time/rate"
)

type Engine struct {
	Stages []stages.Stage
	Pool   func(opts model.RunOptions) runner.Pool
	Sink   sinks.Sink
}

func (e Engine) Run(ctx context.Context, task model.Task, opts model.RunOptions) error {
	pool := e.Pool(opts)
	for _, st := range e.Stages {
		jobs, err := st.Build(ctx, task, opts, e.Sink)
		if err != nil {
			return err
		}
		// 3) sink.Emit(stage_start)
		// 4) pool.Run(jobs)
		if err = pool.Run(ctx, jobs); err != nil {
			return err
		}
		if err = st.After(ctx, task, opts, e.Sink); err != nil {
			return err
		}
		// 5) sink.Emit(stage_done)

	}
	return nil
}

func DefaultPool(opts model.RunOptions) runner.Pool {
	var lim *rate.Limiter
	if opts.Limits.MaxRate > 0 {
		burst := opts.Limits.MaxRate
		lim = rate.NewLimiter(rate.Limit(opts.Limits.MaxRate), burst)
	}
	var gate runner.HostGate
	if opts.Limits.MaxInFlightPerHost > 0 {
		gate = &runner.PerHostGate{N: opts.Limits.MaxInFlightPerHost}
	}
	return runner.Pool{
		MaxInFlight: opts.Limits.MaxInFlight,
		Limiter:     lim,
		Gate:        gate,
	}
}
