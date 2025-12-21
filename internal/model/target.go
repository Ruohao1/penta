package model

import (
	"net/netip"
	"strconv"
)

type Target struct {
	Type TargetType

	// Type == URL
	URL string

	// Type == Host
	Hostname string
	Addr     netip.Addr
}

type TargetType string

const (
	TargetTypeHost TargetType = "host"
	TargetTypeURL  TargetType = "url"
)

func NewTargetURL(url string) *Target {
	return &Target{
		Type: TargetTypeURL,
		URL:  url,
	}
}

func NewTargetHostFromIP(addr netip.Addr) *Target {
	return NewTargetHost("", addr)
}

func NewTargetHostFromHostname(hostname string) *Target {
	return NewTargetHost(hostname, netip.Addr{})
}

func NewTargetHost(hostname string, addr netip.Addr) *Target {
	return &Target{
		Type:     TargetTypeHost,
		Hostname: hostname,
		Addr:     addr,
	}
}

func (t *Target) Address(port Port) string {
	var address string
	switch {
	case t.Type == TargetTypeURL:
		address = t.URL
	case t.Hostname != "":
		address = t.Hostname
	case t.Addr.IsValid():
		address = t.Addr.String()
	default:
		return ""
	}

	return address + ":" + strconv.Itoa(port.Port)
}
