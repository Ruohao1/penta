package hosts

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"runtime"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/vishvananda/netlink"
)

type arpProber struct{}

func (p *arpProber) Name() string { return "arp" }

func (p *arpProber) Probe(ctx context.Context, target model.Target, opts model.RunOptions) (model.Finding, error) {
	host := target.MakeHost()
	host.State = model.HostStateUnknown

	finding := model.Finding{
		Check: "arp_probe",
		Proto: model.ProtocolARP,
		Host:  &host,
		Meta:  map[string]any{},
	}

	if !canUseARP(target) {
		finding.Reason = fmt.Sprintf("arp_unsupported")
		return finding, nil
	}

	neigh, err := lookupARP(target.Addr)
	if err != nil {
		host.State = model.HostStateDown

		finding.Meta["err"] = err.Error()
		return finding, nil
	}

	switch neigh.State {
	case netlink.NUD_REACHABLE, netlink.NUD_STALE, netlink.NUD_DELAY, netlink.NUD_PROBE:
		host.State = model.HostStateUp
		host.MAC = neigh.HardwareAddr.String()

		finding.Reason = "arp_reachable"
		return finding, nil

	case netlink.NUD_INCOMPLETE, netlink.NUD_FAILED:
		host.State = model.HostStateDown

		finding.Reason = "arp_incomplete"
		return finding, nil
	}
	return finding, nil
}

func lookupARP(ip netip.Addr) (netlink.Neigh, error) {
	linkIndex, err := lookupLinkIndex(ip)
	if err != nil {
		return netlink.Neigh{}, err
	}

	neighbors, err := netlink.NeighList(linkIndex, netlink.FAMILY_V4)
	if err != nil {
		return netlink.Neigh{}, err
	}
	for _, n := range neighbors {
		if n.IP.Equal(net.ParseIP(ip.String())) {
			return n, nil
		}
	}
	return netlink.Neigh{}, nil
}

func lookupLinkIndex(ip netip.Addr) (int, error) {
	linkList, err := netlink.LinkList()
	if err != nil {
		return -1, err
	}
	for _, link := range linkList {
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			return -1, err
		}
		for _, a := range addrs {
			ipNet := a.IPNet
			if ipNet.IP.To4() != nil && ipNet.Contains(ip.AsSlice()) {
				return link.Attrs().Index, nil
			}
		}
	}
	return -1, fmt.Errorf("link not found for ip %s", ip.String())
}

func canUseARP(target model.Target) bool {
	if !target.Addr.Is4() || runtime.GOOS != "linux" {
		return false
	}
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			ipNet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			if ipNet.IP.To4() != nil && ipNet.Contains(target.Addr.AsSlice()) {
				return true // same L2 subnet
			}
		}
	}
	return false
}
