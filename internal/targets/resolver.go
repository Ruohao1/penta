package targets

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

func Resolve(expr string) ([]netip.Addr, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, fmt.Errorf("empty targets expression")
	}

	var out []netip.Addr

	parts := strings.SplitSeq(expr, ",")
	for raw := range parts {
		part := strings.TrimSpace(raw)
		if part == "" {
			continue
		}

		switch {
		case strings.Contains(part, "/"):
			ips, err := expandCIDR(part, true)
			if err != nil {
				return nil, fmt.Errorf("parse %q as cidr: %w", part, err)
			}
			out = append(out, ips...)

		case strings.Contains(part, "-"):
			ips, err := expandRange(part)
			if err != nil {
				return nil, fmt.Errorf("parse %q as range: %w", part, err)
			}
			out = append(out, ips...)

		default:
			ip, err := netip.ParseAddr(part)
			if err != nil {
				return nil, fmt.Errorf("parse %q as ip: %w", part, err)
			}
			out = append(out, ip)
		}
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no valid targets in %q", expr)
	}

	return out, nil
}

func expandCIDR(expr string, skipReservedAddr bool) ([]netip.Addr, error) {
	pfx, err := netip.ParsePrefix(expr)
	if err != nil {
		return nil, err
	}

	pfx = pfx.Masked()

	if pfx.Addr().Is4() {
		bits := pfx.Bits()
		hostBits := 32 - bits
		if hostBits > 16 { // > 65,536 addresses
			return nil, fmt.Errorf("CIDR %s too large to expand", expr)
		}
	}

	var (
		res      []netip.Addr
		ip       = pfx.Addr()
		isIPv4   = ip.Is4()
		ones     = pfx.Bits()
		hasBcast = isIPv4 && ones <= 30
	)

	var (
		networkAddr   netip.Addr
		broadcastAddr netip.Addr
	)

	if hasBcast && skipReservedAddr {
		networkAddr = ip
		broadcastAddr = ipv4Broadcast(pfx)
	}

	for cur := ip; pfx.Contains(cur); cur = cur.Next() {
		if skipReservedAddr && hasBcast {
			if cur == networkAddr || cur == broadcastAddr {
				continue
			}
		}
		res = append(res, cur)
	}

	return res, nil
}

// ipv4Broadcast computes the broadcast address for an IPv4 prefix.
func ipv4Broadcast(pfx netip.Prefix) netip.Addr {
	ip := pfx.Addr()
	ip4 := ip.As4()

	bits := pfx.Bits()

	hostBits := 32 - bits
	base := uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
	bcast := base | ((1 << hostBits) - 1)

	return netip.AddrFrom4([4]byte{
		byte(bcast >> 24),
		byte(bcast >> 16),
		byte(bcast >> 8),
		byte(bcast),
	})
}

func expandRange(expr string) ([]netip.Addr, error) {
	parts := strings.Split(strings.TrimSpace(expr), ".")
	if len(parts) != 4 {
		return nil, fmt.Errorf("range %q is malformed (got %d octets)", expr, len(parts))
	}

	// ranges[i] = all possible values for octet i
	ranges := make([][]int, 4)

	for i, raw := range parts {
		seg := strings.TrimSpace(raw)
		if seg == "" {
			return nil, fmt.Errorf("empty octet in %q", expr)
		}

		bounds := strings.SplitN(seg, "-", 2)
		if len(bounds) == 1 {
			// single value, like "10"
			v, err := strconv.Atoi(bounds[0])
			if err != nil || v < 0 || v > 255 {
				return nil, fmt.Errorf("invalid octet %q in %q", seg, expr)
			}
			ranges[i] = []int{v}
			continue
		}

		// range, like "1-2"
		start, err := strconv.Atoi(strings.TrimSpace(bounds[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid range start %q in %q", seg, expr)
		}
		end, err := strconv.Atoi(strings.TrimSpace(bounds[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid range end %q in %q", seg, expr)
		}

		if start < 0 || start > 255 || end < 0 || end > 255 || end < start {
			return nil, fmt.Errorf("invalid octet range %d-%d in %q", start, end, expr)
		}

		vals := make([]int, 0, end-start+1)
		for v := start; v <= end; v++ {
			vals = append(vals, v)
		}
		ranges[i] = vals
	}

	// guard against explosion
	total := 1
	for _, r := range ranges {
		total *= len(r)
	}
	const maxAddrs = 1_000_000
	if total > maxAddrs {
		return nil, fmt.Errorf("expression %q expands to %d addresses (> %d)", expr, total, maxAddrs)
	}

	out := make([]netip.Addr, 0, total)
	for _, a := range ranges[0] {
		for _, b := range ranges[1] {
			for _, c := range ranges[2] {
				for _, d := range ranges[3] {
					addr := netip.AddrFrom4([4]byte{
						byte(a), byte(b), byte(c), byte(d),
					})
					out = append(out, addr)
				}
			}
		}
	}
	return out, nil
}
