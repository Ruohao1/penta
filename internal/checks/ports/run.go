package ports

import (
	"context"
	"time"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/runner"
	"golang.org/x/time/rate"
)

type PortJob struct {
	Host *model.Host
	Port int
}

func Run(ctx context.Context, opts model.RunOptions, emit func(context.Context, model.Event)) error {
	jobs := make([]PortJob, 0, len(opts.Targets))
	for _, t := range opts.Targets {
		if t.Host == nil {
			continue
		}
		for _, p := range t.Host.Ports {
			jobs = append(jobs, PortJob{Host: t.Host, Port: p.Number})
		}
	}

	emit(ctx, model.NewEventWithProgress(model.EventScanStart, len(jobs)))

	// Global rate limiter
	var limiter *rate.Limiter
	if opts.MaxRate > 0 {
		burst := opts.MaxInFlight * 2
		if burst < 1 {
			burst = 1
		}
		if burst > opts.MaxRate {
			burst = opts.MaxRate
		}
		limiter = rate.NewLimiter(rate.Limit(opts.MaxRate), burst)
	}

	gate := runner.NewHostGate(opts.MaxInFlightPerHost)

	scanner := &TCPScanner{} // should use shared netprobe underneath

	err := runner.RunPool(ctx, jobs, opts.MaxInFlight, limiter, func(ctx context.Context, job PortJob) error {
		hostKey := job.Host.Addr.String()
		release := gate.Acquire(hostKey)
		defer release()

		f, err := scanner.ScanPort(ctx, job.Host, job.Port, opts)
		if err != nil {
			ev := model.NewEvent(model.EventFinding)
			ev.Finding = &model.Finding{
				Time:   time.Now().UTC(),
				Check:  "port_scan_error",
				Host:   job,
				Port:   job.Port,
				Proto:  "tcp",
				Meta:   map[string]any{"error": err.Error()},
				Reason: "scan_error",
			}
			emit(ctx, ev)
			return nil // keep scanning
		}

		ev := model.NewEvent(model.EventFinding)
		ev.Finding = &f
		emit(ctx, ev)
		return nil
	})
	if err != nil {
		return err
	}

	emit(ctx, model.NewEvent(model.EventScanDone))
	return nil
}
