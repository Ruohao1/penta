package model

import "time"

// EventType tells the consumer what kind of event this is.
type EventType string

const (
	// Findings / results
	EventFinding EventType = "finding"

	// Host lifecycle
	EventHostUp   EventType = "host_up"
	EventHostDown EventType = "host_down"

	// Port / service
	EventPortOpen   EventType = "port_open"
	EventPortClosed EventType = "port_closed"

	// Probe execution
	EventProbeStart EventType = "probe_start"
	EventProbeDone  EventType = "probe_done"

	// Engine lifecycle
	EventEngineStart EventType = "engine_start"
	EventEngineStop  EventType = "engine_stop"
	EventEngineDone  EventType = "engine_done"

	// State / control
	EventIdle EventType = "idle"
	EventDone EventType = "done"

	// Observability
	EventError EventType = "error"
	EventLog   EventType = "log"
)

// Event is what Engine emits over the channel.
// It's intentionally generic: one stream, multiple event kinds.
type Event struct {
	Time time.Time `json:"time"`
	Type EventType `json:"type"`

	// Optional common fields for routing / display
	Target string `json:"target,omitempty"` // usually IP/hostname

	// Payloads (only one of these is typically non-nil per event)
	Finding  *Finding  `json:"finding,omitempty"`  // for EventTypeFinding
	Host     *Host     `json:"host,omitempty"`     // for host_start/host_done snapshots
	Progress *Progress `json:"progress,omitempty"` // for EventTypeProgress

	// For log/error events
	Message string `json:"message,omitempty"` // human-readable message
	Err     string `json:"err,omitempty"`     // error text if any
}

func NewEvent(t EventType) Event {
	return Event{Type: t}
}

func NewEventWithProgress(t EventType, total int) Event {
	progress := &Progress{TotalTargets: total}
	return Event{Type: t, Progress: progress}
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
