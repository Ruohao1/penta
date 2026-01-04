package tcpconnect

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/checks"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/netprobe"
)

var _ checks.Checker = (*Checker)(nil)

type Input struct {
	Endpoint model.Endpoint
	Opts     model.RunOptions
}

type Checker struct {
	Dialer netprobe.Dialer
}

func New() Checker {
	return Checker{Dialer: netprobe.NetDialer{}}
}

func (p *Checker) Name() string { return "tcp_connect" }

func (p *Checker) Check() checks.CheckFn {
	return func(ctx context.Context, in any, emit checks.EmitFn) error {
		req, ok := in.(Input)
		if !ok {
			return fmt.Errorf("%s: want %T, got %T", p.Name(), Input{}, in)
		}

		// Emit "raw" domain outputs (job/stage will wrap to model.Event)
		finding := model.Finding{
			ObservedAt: time.Now().UTC(),
			Check:      p.Name(),
			Proto:      model.ProtocolTCP,
			Endpoint:   req.Endpoint,

			Severity: "info",
			Meta:     map[string]any{},
		}

		if req.Endpoint.Kind != model.EndpointNet {

			ev := model.Event{
				EmittedAt: time.Now().UTC(),
				Type:      model.EventError,
				Target:    req.Endpoint.String(),
				Finding:   &finding,

				Err: fmt.Sprintf("unsupported endpoint kind: %s", req.Endpoint.Kind),
			}
			emit(ev)

			return nil
		}

		res := netprobe.TCPConnect(ctx, p.Dialer, req.Endpoint.String(), req.Opts.Timeouts.TCP)
		finding.RTTMs = res.ElapsedMs
		finding.Status = res.Reason
		finding.Meta["ok"] = res.OK

		ev := model.NewFindingEvent(&finding)
		emit(ev)
		return nil
	}
}
