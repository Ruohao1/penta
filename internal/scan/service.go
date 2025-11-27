package scan

import (
	"context"
	"fmt"
	"strings"
)

type EngineName string

const (
	NativeEngine EngineName = "native"
	NmapEngine   EngineName = "nmap"
)

func (e *EngineName) String() string {
	if e == nil {
		return ""
	}
	return string(*e)
}

func (e *EngineName) Set(s string) error {
	switch strings.ToLower(s) {
	case "native":
		*e = NativeEngine
	case "nmap":
		*e = NmapEngine
	default:
		return fmt.Errorf("invalid engine %q (valid: native|nmap)", s)
	}
	return nil
}

func (e *EngineName) Type() string {
	return "engine"
}

type Engine interface {
	ScanHosts(ctx context.Context, target string, opts HostsOptions) ([]string, error)
	ScanPorts(ctx context.Context, target string) ([]string, error)
}

type Service struct {
	engine Engine
}

func NewService(engine Engine) *Service {
	return &Service{engine: engine}
}

func (s *Service) ScanHosts(ctx context.Context, target string, opts HostsOptions) ([]string, error) {
	return s.engine.ScanHosts(ctx, target, opts)
}

func (s *Service) ScanPorts(ctx context.Context, target string) ([]string, error) {
	return s.engine.ScanPorts(ctx, target)
}
