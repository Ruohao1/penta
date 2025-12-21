package model

import (
	"net/netip"
	"time"
)

type HostState string

const (
	HostStateUnknown HostState = "unknown"
	HostStateUp      HostState = "up"
	HostStateDown    HostState = "down"
)

type Host struct {
	Addr      netip.Addr `json:"addr"`
	Hostnames []string   `json:"hostnames,omitempty"`
	State     HostState  `json:"state,omitempty"`

	MAC    string `json:"mac,omitempty"`
	Vendor string `json:"vendor,omitempty"`

	Ports []Port `json:"ports,omitempty"`

	// Network-level context
	NetworkCIDR string `json:"network_cidr,omitempty"` // "192.168.2.0/24"
	Scope       string `json:"scope,omitempty"`        // "in","out","unknown"

	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

func (t Target) MakeHost() Host {
	host := Host{
		Addr:      t.Addr,
		Hostnames: []string{t.Hostname},
	}
	return host
}
