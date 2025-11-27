package hosts

import (
	"context"
	"net/netip"
	"slices"
	"sync"

	"github.com/Ruohao1/penta/internal/scan"
	"github.com/Ruohao1/penta/internal/utils"
	"golang.org/x/time/rate"
)

func Discover(ctx context.Context, addrs []netip.Addr, opts scan.HostsOptions, emit func(scan.Result) error) error {
	logger := utils.LoggerFrom(ctx)

	rateLimit := opts.Rate
	if rateLimit <= 0 {
		rateLimit = 200
	}
	limiter := rate.NewLimiter(rate.Limit(rateLimit), rateLimit)

	const workers = 64
	addrCh := make(chan netip.Addr)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for ip := range addrCh {
			logger.Debug().Str("ip", ip.String()).Msg("host probe started")

			if err := limiter.Wait(ctx); err != nil {
				errCh <- err
				return
			}
			res, err := probeOne(ctx, ip, opts)
			if err != nil {
				// log but still emit best-effort result
				logger.Error().Err(err).Str("ip", ip.String()).Msg("host probe failed")
			}

			if err := emit(res); err != nil {
				errCh <- err
				return
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		defer close(addrCh)
		for _, ip := range addrs {
			select {
			case <-ctx.Done():
				return
			case addrCh <- ip:
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
		return nil
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func probeOne(ctx context.Context, ip netip.Addr, opts scan.HostsOptions) (scan.Result, error) {
	// TODO: Check if it has privileged capabilities
	// To do ARP, or ICMP probing first
	tcpProber := tcpProber{}
	res, err := tcpProber.Probe(ctx, ip, opts)
	if err != nil {
		return res, err
	}

	// If TCP says "up", we're done.
	if res.Status == scan.StatusUp {
		return res, nil
	}

	if slices.Contains(opts.Methods, scan.MethodARP) {
		arpProber := arpProber{}
		arpRes, err := arpProber.Probe(ctx, ip, opts)
		if err != nil {
			// fall back to the TCP result, but mark error in meta
			if res.Meta == nil {
				res.Meta = map[string]any{}
			}
			res.Meta["arp_err"] = err.Error()
			if res.Status == scan.StatusUnknown {
				res.Status = scan.StatusDown
				res.Meta["signal"] = "no_response"
			}
			return res, nil
		}

		if arpRes.Status == scan.StatusUp {
			return arpRes, nil
		}
	}

	// TCP unknown/down + ARP down/unknown -> mark down
	if res.Status == scan.StatusUnknown {
		res.Status = scan.StatusDown
		if res.Meta == nil {
			res.Meta = map[string]any{}
		}
		if _, ok := res.Meta["signal"]; !ok {
			res.Meta["signal"] = "no_response"
		}
	}

	return res, nil
}
