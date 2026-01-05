package model

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Event is what Engine emits over the channel.
// It's intentionally generic: one stream, multiple event kinds.
type Event struct {
	EmittedAt time.Time `json:"emitted_at"`
	Type      EventType `json:"type"`
	Stage     string    `json:"stage,omitempty"`
	Target    string    `json:"target,omitempty"` // usually IP/hostname

	// Payloads (only one of these is typically non-nil per event)
	Finding   *Finding        `json:"finding,omitempty"` // for EventTypeFinding
	HostState *HostStateEvent `json:"host_state,omitempty"`
	Progress  *Progress       `json:"progress,omitempty"` // for EventTypeProgress

	// For log/error events
	Message string `json:"message,omitempty"` // human-readable message
	Err     string `json:"err,omitempty"`     // error text if any
}

// EventType tells the consumer what kind of event this is.
type EventType string

const (
	// Findings / results
	EventFinding EventType = "finding"

	// Host lifecycle
	EventHostState EventType = "host_state"

	// Port / service
	EventPortOpen   EventType = "port_open"
	EventPortClosed EventType = "port_closed"

	// Probe execution
	EventProbeStart EventType = "probe_start"
	EventProbeDone  EventType = "probe_done"

	// Scan execution
	EventScanStart EventType = "scan_start"
	EventScanDone  EventType = "scan_done"

	// Engine lifecycle
	EventEngineStart EventType = "engine_start"
	EventEngineStop  EventType = "engine_stop"
	EventEngineDone  EventType = "engine_done"

	// State / control
	EventIdle EventType = "idle"
	EventDone EventType = "done"

	// Observability
	EventError   EventType = "error"
	EventLog     EventType = "log"
	EventUnknown EventType = "unknown"
)

type HostStateEvent struct {
	Host   string         `json:"host"`
	State  HostState      `json:"state"` // up/down
	Via    string         `json:"via,omitempty"`
	Port   int            `json:"port,omitempty"`
	Reason string         `json:"reason,omitempty"`
	Meta   map[string]any `json:"meta,omitempty"`
}

func NewEvent(t EventType) Event {
	return Event{Type: t}
}

func NewEventWithProgress(t EventType, total int) Event {
	progress := &Progress{TotalTargets: total}
	return Event{Type: t, Progress: progress}
}

func NewFindingEvent(f *Finding) Event {
	if f == nil {
		fmt.Println("nil finding")
	}
	return Event{Type: EventFinding, Finding: f}
}

// Progress is for high-level progress reporting (TUI / verbose mode).
type Progress struct {
	TotalTargets   int `json:"total_targets,omitempty"`   // number of hosts planned
	ProcessedHosts int `json:"processed_hosts,omitempty"` // hosts fully done
	ActiveHosts    int `json:"active_hosts,omitempty"`    // currently being scanned

	// Optional fine-grained metrics
	TotalFindings int     `json:"total_findings,omitempty"`
	Percent       float64 `json:"percent,omitempty"` // 0.0–100.0 best-effort
}

func (ev Event) String() string {
	ts := ev.EmittedAt
	if ts.IsZero() {
		ts = time.Now()
	}

	parts := []string{
		ts.Format(time.RFC3339),
		string(ev.Type),
	}

	if ev.Stage != "" {
		parts = append(parts, "stage="+ev.Stage)
	}
	if ev.Target != "" {
		parts = append(parts, "target="+ev.Target)
	}

	switch ev.Type {
	case EventFinding:
		if ev.Finding != nil {
			parts = append(parts, findingSummary(*ev.Finding)...)
		}

	case EventHostState:
		if ev.HostState != nil {
			hs := ev.HostState
			if hs.Host != "" {
				parts = append(parts, "host="+hs.Host)
			}
			parts = append(parts, "state="+string(hs.State))
			if hs.Via != "" {
				parts = append(parts, "via="+hs.Via)
			}
			if hs.Port != 0 {
				parts = append(parts, fmt.Sprintf("port=%d", hs.Port))
			}
			if hs.Reason != "" {
				parts = append(parts, "reason="+short(hs.Reason, 80))
			}
		}

	case EventError:
		if ev.Err != "" {
			parts = append(parts, "err="+short(ev.Err, 140))
		} else if ev.Message != "" {
			parts = append(parts, "err="+short(ev.Message, 140))
		}

	case EventLog:
		if ev.Message != "" {
			parts = append(parts, "msg="+short(ev.Message, 140))
		}
		if ev.Err != "" {
			parts = append(parts, "err="+short(ev.Err, 140))
		}

	default:
		if ev.Progress != nil {
			parts = append(parts, progressSummary(*ev.Progress)...)
		}
		if ev.Message != "" {
			parts = append(parts, "msg="+short(ev.Message, 140))
		}
		if ev.Err != "" {
			parts = append(parts, "err="+short(ev.Err, 140))
		}
	}

	return strings.Join(parts, " ")
}

func findingSummary(f Finding) []string {
	out := make([]string, 0, 10)

	if !f.ObservedAt.IsZero() {
		out = append(out, "obs="+f.ObservedAt.Format(time.RFC3339))
	}

	if f.Check != "" {
		out = append(out, "check="+f.Check)
	}

	// Protocol is usually a string alias; fmt.Sprint handles it safely.
	if string(f.Proto) != "" {
		out = append(out, "proto="+string(f.Proto))
	}

	if !f.Endpoint.IsZero() {
		if ep := f.Endpoint.String(); ep != "" {
			out = append(out, "ep="+ep)
		}
	}

	if f.Status != "" {
		out = append(out, "status="+f.Status)
	}

	if f.Severity != "" {
		out = append(out, "sev="+f.Severity)
	}

	if f.RTTMs > 0 {
		out = append(out, fmt.Sprintf("rtt=%.2fms", f.RTTMs))
	}

	// Optional: include only a *small* meta hint (avoid dumping huge blobs)
	if len(f.Meta) > 0 {
		out = append(out, "meta="+metaHint(f.Meta, 3))
	}

	return out
}

func progressSummary(p Progress) []string {
	out := make([]string, 0, 6)
	if p.TotalTargets != 0 {
		out = append(out, fmt.Sprintf("total=%d", p.TotalTargets))
	}
	if p.ProcessedHosts != 0 {
		out = append(out, fmt.Sprintf("done=%d", p.ProcessedHosts))
	}
	if p.ActiveHosts != 0 {
		out = append(out, fmt.Sprintf("active=%d", p.ActiveHosts))
	}
	if p.TotalFindings != 0 {
		out = append(out, fmt.Sprintf("findings=%d", p.TotalFindings))
	}
	if p.Percent != 0 {
		out = append(out, fmt.Sprintf("pct=%.1f", p.Percent))
	}
	return out
}

func metaHint(m map[string]any, maxKeys int) string {
	if maxKeys <= 0 {
		return "{}"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) > maxKeys {
		keys = keys[:maxKeys]
	}
	// only show keys (not values) to avoid secrets / huge prints
	return "{" + strings.Join(keys, ",") + "}"
}

func short(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 || len(s) <= max {
		return s
	}
	if max <= 1 {
		return "…"
	}
	return s[:max-1] + "…"
}
