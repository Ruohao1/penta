package hosts

import (
	"context"

	"github.com/Ruohao1/penta/internal/scan"
	"github.com/Ruohao1/penta/internal/targets"
)

type Prober interface {
	Name() string
	Probe(ctx context.Context, target targets.Target, opts scan.HostsOptions) (scan.Result, error)
}

func NewProbers(methods []scan.Method) []Prober {
	ps := []Prober{}
	for _, m := range methods {
		switch m {
		case scan.MethodARP:
			ps = append(ps, &arpProber{})
		case scan.MethodTCP:
			ps = append(ps, &tcpProber{})
			// case scan.MethodICMP:
			// 	ps = append(ps, &icmpProber{})
		}
	}
	return ps
}
