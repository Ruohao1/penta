package model

import "net/netip"

type Network struct {
	CIDR netip.Prefix `json:"cidr"` // e.g. 192.168.2.0/24

	// Optional human metadata
	Name        string `json:"name,omitempty"`        // "office-lan", "dmz"
	Description string `json:"description,omitempty"` // free text

	// Scope / policy
	Scope    string `json:"scope,omitempty"`    // "in","out","unknown"
	Source   string `json:"source,omitempty"`   // "scope.txt","inferred"
	Priority int    `json:"priority,omitempty"` // which rule matched

	// Topology hints
	GatewayIP  netip.Addr `json:"gateway_ip,omitempty"`
	GatewayMAC string     `json:"gateway_mac,omitempty"`
	VLANID     int        `json:"vlan_id,omitempty"` // if you ever learn it

	// Aggregate stats (computed in report, not during scan)
	HostCount      int `json:"host_count,omitempty"`
	UpHosts        int `json:"up_hosts,omitempty"`
	OpenPortsTotal int `json:"open_ports_total,omitempty"`
}
