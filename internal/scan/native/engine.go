package native

import (
	"context"
	"fmt"
	"time"

	"github.com/Ruohao1/penta/internal/scan"
	"github.com/Ruohao1/penta/internal/scan/native/hosts"
	"github.com/Ruohao1/penta/internal/utils/parser"
)

type Engine struct {
	MaxConcurrency int
	Timeout        time.Duration
}

func NewEngine(maxConcurrency int, timeout time.Duration) *Engine {
	return &Engine{
		MaxConcurrency: maxConcurrency,
		Timeout:        timeout,
	}
}

func (e *Engine) ScanHosts(ctx context.Context, target string, opts scan.HostsOptions) ([]string, error) {
	hostListParser := parser.NewHostListParser()
	ips, err := hostListParser.Parse(target)
	if err != nil {
		return nil, err
	}

	err = hosts.Discover(ctx, ips, opts, func(res scan.Result) error {
		// collect only "up" hosts, but ALWAYS append when they show up
		if res.Status != scan.StatusUp {
			return nil
		}

		fmt.Printf("host=%s status=%s method=%s meta=%v\n",
			res.Addr.String(),
			res.Status,
			res.Method,
			res.Meta,
		)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (e *Engine) ScanPorts(ctx context.Context, target string) ([]string, error) {
	return nil, nil
}

var _ scan.Engine = (*Engine)(nil)
