package model

import (
	"time"
)

type Finding struct {
	Time     time.Time `json:"time"`
	Check    string    `json:"check"`           // e.g. "tcp_open", "http_probe"
	Proto    Protocol  `json:"proto,omitempty"` // "tcp", "udp", "http", "https"
	Severity string    `json:"severity,omitempty"`

	Host   *Host   `json:"host,omitempty"`
	RTTMs  float64 `json:"rtt_ms,omitempty"`
	Reason string  `json:"reason,omitempty"`

	Meta map[string]any `json:"meta,omitempty"` // per-check extras
}
