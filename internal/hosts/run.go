package hosts

import (
	"context"
	"sync"
	"time"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/utils"
	"golang.org/x/time/rate"
)

func Run(ctx context.Context, opts model.RunOptions, emit func(context.Context, model.Event)) error {
	var ev model.Event

	logger := utils.LoggerFrom(ctx)
	logger.Info().Int("count", len(opts.Targets)).Msg("host probe started")

	ev = model.NewEventWithProgress(model.EventProbeStart, len(opts.Targets))
	emit(ctx, ev)

	var limiter *rate.Limiter
	if opts.MaxRate > 0 {
		// burst: allow short spikes but keep it bounded
		burst := opts.Concurrency * 2
		if burst < 1 {
			burst = 1
		}
		if burst > opts.MaxRate {
			burst = opts.MaxRate
		} // don't exceed max-rate
		limiter = rate.NewLimiter(rate.Limit(opts.MaxRate), burst)
	}

	workers := opts.Concurrency

	targetCh := make(chan model.Target)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for target := range targetCh {
			if limiter != nil {
				if err := limiter.Wait(ctx); err != nil {
					errCh <- err
					return
				}
			}

			finding, err := probeOne(ctx, target, opts)
			if err != nil {
				ev := model.NewEvent(model.EventFinding)
				host := target.MakeHost()
				ev.Finding = &model.Finding{
					Time:   time.Now().UTC(),
					Check:  "probe_error",
					Host:   &host,
					Meta:   map[string]any{"error": err.Error()},
					Reason: "probe_error",
				}
				emit(ctx, ev)
				continue
			}

			ev := model.NewEvent(model.EventFinding)
			ev.Finding = &finding

			emit(ctx, ev)
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		defer close(targetCh)
		for _, target := range opts.Targets {
			select {
			case <-ctx.Done():
				return
			case targetCh <- target:
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:

		ev := model.NewEvent(model.EventProbeDone)
		emit(ctx, ev)
	case err := <-errCh:
		return err
	}

	return nil
}

func probeOne(ctx context.Context, target model.Target, opts model.RunOptions) (model.Finding, error) {
	var prober Prober
	// TODO : swicth betwwen high privilege or not, currently just unpriviled probing
	prober = &tcpProber{}
	result, err := prober.Probe(ctx, target, opts)
	if err != nil {
		return result, err
	}
	if result.Host.State == model.HostStateUp {
		return result, nil
	}
	prober = &arpProber{}
	return prober.Probe(ctx, target, opts)

	// if opts.ARP {
	// 	prober = &arpProber{}
	// 	// } else if opts.ICMP {
	// 	// 	prober = &icmpProber{}
	// } else {
	// 	prober = &tcpProber{}
	// }
	//
	// prober = &tcpProber{}
	// return prober.Probe(ctx, target, opts)
}
