package hosts

import (
	"context"
	"errors"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/utils"
)

type tcpProber struct{}

func (p *tcpProber) Name() string { return "tcp" }

func (p *tcpProber) Probe(ctx context.Context, target model.Target, opts model.RunOptions) (model.Finding, error) {
	logger := utils.LoggerFrom(ctx)
	logger.Debug().Str("target", target.Addr.String()).Msg("probing")

	ports := []model.Port{
		{Port: 80},
		{Port: 443},
	}

	host := target.MakeHost()
	host.State = model.HostStateUnknown

	finding := model.Finding{
		Check: "tcp_probe",
		Proto: model.ProtocolTCP,
		Host:  &host,
		Meta:  map[string]any{},
	}

	timeout := opts.Timeout
	dialer := net.Dialer{Timeout: timeout}

	// Track best signals across all probes
	var (
		sawTimeout bool
		lastErrStr string
	)

	touch := func(now time.Time) {
		// Finding timestamp (your output shows this exists)
		if finding.Time.IsZero() {
			finding.Time = now
		} else {
			finding.Time = now
		}

		// Host timestamps
		if host.FirstSeen.IsZero() {
			host.FirstSeen = now
		}
		host.LastSeen = now
	}

	for _, port := range ports {
		select {
		case <-ctx.Done():
			return finding, ctx.Err()
		default:
		}

		address := target.Address(port)
		start := time.Now()
		conn, err := dialer.DialContext(ctx, "tcp", address)
		elapsed := time.Since(start)
		now := time.Now()
		finding.RTTMs = float64(elapsed) / float64(time.Millisecond)

		// Success => definitive UP
		if err == nil {
			_ = conn.Close()

			touch(now)
			host.State = model.HostStateUp
			host.Ports = append(host.Ports, port)
			finding.Reason = "tcp_connect_success"
			return finding, nil
		}

		// Context cancellation should bubble out
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return finding, err
		}

		kind, reason := classifyDialErr(err)
		lastErrStr = err.Error()

		switch kind {
		case dialUp:
			// Definitive UP (RST/refused/reset etc.)
			touch(now)
			host.State = model.HostStateUp
			host.Ports = append(host.Ports, port)
			finding.Reason = reason
			return finding, nil

		case dialHardDown:
			// Definitive unreachable (routing/host/network unreachable)
			host.State = model.HostStateDown
			finding.Reason = reason
			finding.Meta["error"] = lastErrStr
			return finding, nil

		case dialTimeout:
			// Not definitive; keep trying other ports
			sawTimeout = true
			// Keep the most useful "reason" only if nothing else exists
			if finding.Reason == "" {
				finding.Reason = "timeout"
			}
			continue

		case dialOther:
			// Not definitive; keep trying other ports, but record something
			if _, ok := finding.Meta["error"]; !ok {
				finding.Meta["error"] = lastErrStr
			}
			continue
		}
	}

	// If we got here: no definitive UP and no definitive unreachable.
	// For discovery, treat as DOWN (nmap would not mark it up).
	host.State = model.HostStateDown

	if sawTimeout && finding.Reason == "" {
		finding.Reason = "timeout"
	}
	// Keep an error string for debugging, but don’t pretend it’s definitive.
	if lastErrStr != "" {
		if _, ok := finding.Meta["error"]; !ok {
			finding.Meta["error"] = lastErrStr
		}
	}
	if finding.Reason == "" {
		finding.Reason = "no_response"
	}

	return finding, nil
}

type dialKind int

const (
	dialOther dialKind = iota
	dialTimeout
	dialUp
	dialHardDown
)

func classifyDialErr(err error) (dialKind, string) {
	// Timeout detection (covers i/o timeout, context deadline at net layer, etc.)
	var nerr net.Error
	if errors.As(err, &nerr) && nerr.Timeout() {
		return dialTimeout, "timeout"
	}

	// Unwrap syscall errno when possible
	if errno, ok := unwrapErrno(err); ok {
		switch errno {
		case syscall.ECONNREFUSED:
			// RST => host is alive; port closed
			return dialUp, "tcp_connect_refused"
		case syscall.ECONNRESET:
			return dialUp, "tcp_connect_reset"
		case syscall.ENETUNREACH:
			return dialHardDown, "net_unreachable"
		case syscall.EHOSTUNREACH:
			return dialHardDown, "host_unreachable"
		}
	}

	// Some stacks return "no route to host" as an OpError without clean errno
	// (we treat unknown stuff as non-definitive, keep scanning).
	return dialOther, "error"
}

func unwrapErrno(err error) (syscall.Errno, bool) {
	// Typical chain: *net.OpError -> *os.SyscallError -> syscall.Errno
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		// Sometimes opErr.Err is already syscall.Errno
		if errno, ok := opErr.Err.(syscall.Errno); ok {
			return errno, true
		}
		var sysErr *os.SyscallError
		if errors.As(opErr.Err, &sysErr) {
			if errno, ok := sysErr.Err.(syscall.Errno); ok {
				return errno, true
			}
		}
	}

	// Direct syscall error
	var sysErr *os.SyscallError
	if errors.As(err, &sysErr) {
		if errno, ok := sysErr.Err.(syscall.Errno); ok {
			return errno, true
		}
	}

	return 0, false
}
