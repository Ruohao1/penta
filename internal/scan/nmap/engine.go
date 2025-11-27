package nmap

import (
	"context"

	"github.com/Ruohao1/penta/internal/scan"
)

type Engine struct {
	BinaryPath string
	ExtraArgs  []string
}

func NewEngine(bin string, extra []string) *Engine {
	return &Engine{BinaryPath: bin, ExtraArgs: extra}
}

func (e *Engine) ScanHosts(ctx context.Context, target string, opts scan.HostsOptions) ([]string, error) {
	// exec.CommandContext(ctx, e.BinaryPath, "-sn", target, ...)
	// parse output…
	return nil, nil
}

func (e *Engine) ScanPorts(ctx context.Context, target string) ([]string, error) {
	// exec.CommandContext(ctx, e.BinaryPath, "-p", "...", target, ...)
	return nil, nil
}

var _ scan.Engine = (*Engine)(nil)
