package ports

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/netprobe"
)

type TCPScanner struct {
	Dialer netprobe.Dialer
}

func (s *TCPScanner) ScanPort(ctx context.Context, host *model.Host, port *model.Port, opts model.RunOptions) (model.Finding, error) {
	addr := fmt.Sprintf("%s:%d", host.Address(), port.Number)
	res := netprobe.TCPConnect(ctx, s.Dialer, addr, opts.TimeoutTCP)

	f := model.Finding{
		Time:  time.Now().UTC(),
		Check: "tcp_port_scan",
		Proto: "tcp",

		Host:   host,
		RTTMs:  res.ElapsedMs,
		Reason: res.Reason,

		Meta:     map[string]any{"reason": res.Reason},
		Severity: "info",
	}

	// You decide your semantics:
	// open => open
	// refused => closed (or “reachable but closed”)
	return f, nil
}
