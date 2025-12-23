package netprobe

import (
	"context"
	"net"
	"time"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type NetDialer struct {
	Dialer net.Dialer
}

func (d NetDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return d.Dialer.DialContext(ctx, network, address)
}

type Result struct {
	OK        bool
	State     string
	Reason    string // "open", "refused", "timeout", "unreachable", "error"
	Err       error  // optional
	ElapsedMs float64
}

func TCPConnect(ctx context.Context, d Dialer, addr string, timeout time.Duration) *Result {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	c, err := d.DialContext(ctx, "tcp", addr)
	elapsed := time.Since(start).Seconds() / 1000.0

	result := &Result{ElapsedMs: elapsed}
	if err != nil {
		return classifyDialErr(result, err)
	}
	_ = c.Close()
	result.OK = true
	result.State = "open"
	result.Reason = "tcp_response"
	return result
}
