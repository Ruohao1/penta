package model

import (
	"net/netip"
	"time"
)

type Host struct {
	Addr      netip.Addr `json:"addr"`
	Hostnames []string   `json:"hostnames,omitempty"`

	MAC    string `json:"mac,omitempty"`
	Vendor string `json:"vendor,omitempty"`

	// Discovery...
	DiscoveryMethod string  `json:"discovery_method,omitempty"`
	DiscoveryProto  string  `json:"discovery_proto,omitempty"`
	DiscoveryPort   int     `json:"discovery_port,omitempty"`
	DiscoveryRTTms  float64 `json:"discovery_rtt_ms,omitempty"`
	DiscoverySignal string  `json:"discovery_signal,omitempty"`

	Ports []Port `json:"ports,omitempty"`

	// Network-level context
	NetworkCIDR string `json:"network_cidr,omitempty"` // "192.168.2.0/24"
	Scope       string `json:"scope,omitempty"`        // "in","out","unknown"

	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}
