package model

import (
	"context"
	"strconv"
)

type TaskKind string

const (
	TaskHostDiscovery TaskKind = "host_discovery"
	TaskPortScan      TaskKind = "port_scan"
	TaskScan          TaskKind = "scan"
	TaskWebProbe      TaskKind = "web_probe"
)

// Task is user intent: WHAT to do + WHAT to touch.
type Task struct {
	Kind TaskKind `json:"kind"`

	Targets []Target `json:"targets"` // parsed expressions: CIDR, range, IP, hostname
	Ports   []int    `json:"ports,omitempty"`

	Wait func(ctx context.Context) error // optional external wait hook, can be nil
}

func NewScanTask(targetsExpr string, portsExpr []string) (Task, error) {
	targets, err := NewTargets(targetsExpr)
	if err != nil {
		return Task{}, err
	}
	ports, err := ParsePorts(portsExpr)
	if err != nil {
		return Task{}, err
	}
	return Task{
		Kind:    TaskScan,
		Targets: targets,
		Ports:   ports,
	}, nil
}

func NewHostDiscoveryTask(targetsExpr string, portsExpr []string) (Task, error) {
	targets, err := NewTargets(targetsExpr)
	if err != nil {
		return Task{}, err
	}
	ports, err := ParsePorts(portsExpr)
	if err != nil {
		return Task{}, err
	}
	return Task{
		Kind:    TaskHostDiscovery,
		Targets: targets,
		Ports:   ports,
	}, nil
}

func NewPortScanTask(targetsExpr string, portsExpr []string) (Task, error) {
	targets, err := NewTargets(targetsExpr)
	if err != nil {
		return Task{}, err
	}
	ports, err := ParsePorts(portsExpr)
	if err != nil {
		return Task{}, err
	}
	return Task{
		Kind:    TaskPortScan,
		Targets: targets,
		Ports:   ports,
	}, nil
}

func ParsePorts(ports []string) ([]int, error) {
	portsInt := make([]int, len(ports))
	for i, port := range ports {
		p, err := strconv.Atoi(port)
		if err != nil {
			return []int{}, err
		}
		portsInt[i] = p
	}
	return portsInt, nil
}

type DiscoveryMethod string

const (
	DiscoveryTCP  DiscoveryMethod = "tcp"  // TCP connect to probe port(s)
	DiscoveryICMP DiscoveryMethod = "icmp" // ping
	DiscoveryARP  DiscoveryMethod = "arp"  // LAN-only
)

func (t Task) ExpandAllTargetsExpr() ([]string, error) {
	out := []string{}
	for _, target := range t.Targets {
		expand, err := target.ExpandAll()
		if err != nil {
			return nil, err
		}
		out = append(out, expand...)
	}
	return out, nil
}
