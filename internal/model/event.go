package model

import "time"

// EventType tells the consumer what kind of event this is.
type EventType int

const (
	EventFinding EventType = 1 << iota

	EventHostUp
	EventPortOpen

	EventProbeStart
	EventProbeDone

	EventEngineStart
	EventEngineStop
	EventEngineDone

	EventIdle
	EventError
	EventLog
	EventDone
)

// Event is what Engine emits over the channel.
// It's intentionally generic: one stream, multiple event kinds.
type Event struct {
	Time time.Time `json:"time"`
	Type EventType `json:"type"`

	// Optional common fields for routing / display
	Target string `json:"target,omitempty"` // usually IP/hostname
	Check  string `json:"check,omitempty"`  // check name, e.g. "tcp_open","http_probe"

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

// Progress is for high-level progress reporting (TUI / verbose mode).
type Progress struct {
	TotalTargets   int `json:"total_targets,omitempty"`   // number of hosts planned
	ProcessedHosts int `json:"processed_hosts,omitempty"` // hosts fully done
	ActiveHosts    int `json:"active_hosts,omitempty"`    // currently being scanned

	// Optional fine-grained metrics
	TotalFindings int     `json:"total_findings,omitempty"`
	Percent       float64 `json:"percent,omitempty"` // 0.0–100.0 best-effort
}
