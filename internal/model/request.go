package model

import (
	"context"
)

type Request struct {
	Mode    Mode        // hosts/ports/web
	Backend BackendName // internal/nmap/...

	// External backend passthrough (tool owns parsing)
	ToolArgs []string // everything after `--` (e.g. ["-sn","10.0.0.0/24"])
}

type Mode string

const (
	ModeHosts Mode = "hosts"
	ModePorts Mode = "ports"
	ModeWeb   Mode = "web"
)

type BackendName string

const (
	BackendInternal BackendName = "internal"
	BackendNmap     BackendName = "nmap"
	BackendNuclei   BackendName = "nuclei"
)

type ExternalTool string

const (
	ToolNmap   ExternalTool = "nmap"
	ToolNuclei ExternalTool = "nuclei"
)

type ToolJob struct {
	Tool ExternalTool
	Args []string       // full argv (or user args + enforced output flags)
	Meta map[string]any // provenance, temp file paths, etc.
}

// ===== Execution interfaces =====

type Checker interface {
	Name() string
	Supports(w WorkItem) bool // quick guard; keeps engine simple
	Run(ctx context.Context, w WorkItem) ([]Finding, error)
}

type Backend interface {
	Name() BackendName
	Run(ctx context.Context, req Request, emit func(Event)) error
}
