package scan

import (
	"fmt"
	"net/netip"
	"strings"
	"time"
)

type Method int

const (
	MethodTCP Method = iota
	MethodICMP
	MethodARP
)

func (m Method) String() string {
	switch m {
	case MethodTCP:
		return "tcp"
	case MethodICMP:
		return "icmp"
	case MethodARP:
		return "arp"
	default:
		return "unknown"
	}
}

func ParseMethods(ss []string) ([]Method, error) {
	out := make([]Method, 0, len(ss))
	for _, s := range ss {
		switch strings.ToLower(s) {
		case "arp":
			out = append(out, MethodARP)
		case "icmp":
			out = append(out, MethodICMP)
		case "tcp":
			out = append(out, MethodTCP)
		default:
			return nil, fmt.Errorf("invalid method %q (valid: arp|icmp|tcp)", s)
		}
	}
	return out, nil
}

type HostsOptions struct {
	EngineName EngineName
	Methods    []Method
	TCPPorts   []int

	Timeout    time.Duration // per probe
	Rate       int           // global probes/sec
	MaxRetries int           // per method per host
}

type Status string

const (
	StatusUnknown Status = "unknown"
	StatusUp      Status = "up"
	StatusDown    Status = "down"
)

type Result struct {
	Addr   netip.Addr
	Status Status
	Method Method
	Meta   map[string]any
}
