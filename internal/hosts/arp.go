package hosts

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"runtime"

	"github.com/Ruohao1/penta/internal/engine"
	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/scan"
	"github.com/vishvananda/netlink"
)

type arpProber struct{}

func (p *arpProber) Name() string { return "arp" }

func (p *arpProber) Probe(ctx context.Context, ip netip.Addr, opts engine.RunOptions) (model.Host, model.Finding, error) {
	finding := model.Finding{
		Check: "arp_probing",
	}
	host := model.Host{
		Addr:            ip,
		State:           model.HostStateUnknown,
		DiscoveryMethod: scan.MethodARP,
	}

	if !canUseARP(ip) {
		host.Reason = fmt.Sprintf("arp_unsupported")
		return host, nil
	}

	neigh, err := lookupARP(ip)
	if err != nil {
		result.Meta["err"] = err.Error()
		result.Status = scan.StatusUnknown
		return result, nil
	}

	switch neigh.State {
	case netlink.NUD_REACHABLE, netlink.NUD_STALE, netlink.NUD_DELAY, netlink.NUD_PROBE:
		result.Status = scan.StatusUp
		result.Meta["mac"] = neigh.HardwareAddr.String()
		result.Meta["dev"] = neigh.LinkIndex
		result.Meta["state"] = neigh.State
		return result, nil

	case netlink.NUD_INCOMPLETE, netlink.NUD_FAILED:
		result.Status = scan.StatusDown
		result.Meta["signal"] = "arp_incomplete"
		result.Meta["state"] = neigh.State
		return result, nil
	}
	return result, nil
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

func canUseARP(target netip.Addr) bool {
	if !target.Is4() || runtime.GOOS != "linux" {
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
			if ipNet.IP.To4() != nil && ipNet.Contains(target.AsSlice()) {
				return true // same L2 subnet
			}
		}
	}
	return false
}
