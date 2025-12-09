package model

import "time"

type Finding struct {
	Time     time.Time      `json:"time"`
	Check    string         `json:"check"`  // e.g. "tcp_open", "http_probe"
	Target   string         `json:"target"` // IP/hostname
	Port     int            `json:"port,omitempty"`
	Proto    string         `json:"proto,omitempty"` // "tcp", "udp", "http", "https"
	Severity string         `json:"severity,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"` // per-check extras
}
