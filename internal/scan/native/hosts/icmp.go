package hosts

//
// import (
// 	"context"
// 	"errors"
// 	"net"
// 	"net/netip"
// 	"os"
// 	"syscall"
// 	"time"
//
// 	"github.com/Ruohao1/penta/internal/scan"
// )
//
// type icmpProber struct{}
//
// func (p *icmpProber) Name() string { return "icmp" }
//
// func (p *icmpProber) Probe(ctx context.Context, ip netip.Addr, opts scan.HostsOptions) (scan.Result, error) {
// 	ports := opts.TCPPorts
// 	if len(ports) == 0 {
// 		ports = []int{22, 80, 443}
// 	}
//
// 	result := scan.Result{
// 		Addr:   ip,
// 		Status: scan.StatusUnknown,
// 		Method: scan.MethodTCP,
// 		Meta:   map[string]any{"ports": ports},
// 	}
//
// 	timeout := opts.Timeout
// 	if timeout <= 0 {
// 		timeout = 500 * time.Millisecond
// 	}
//
// 	dialer := net.Dialer{Timeout: timeout}
//
// 	for _, port := range ports {
// 		select {
// 		case <-ctx.Done():
// 			return result, ctx.Err()
// 		default:
// 		}
//
// 		addr := netip.AddrPortFrom(ip, uint16(port))
// 		start := time.Now()
// 		conn, err := dialer.DialContext(ctx, "tcp", addr.String())
// 		elapsed := time.Since(start)
// 		rttMs := float64(elapsed.Microseconds()) / 1000.0
//
// 		// success: host is definitely up
// 		if err == nil {
// 			_ = conn.Close()
// 			result.Status = scan.StatusUp
// 			result.Meta["port"] = port
// 			result.Meta["rtt_ms"] = rttMs
// 			result.Meta["signal"] = "connect_success"
// 			return result, nil
// 		}
//
// 		// classify network errors
// 		var opErr *net.OpError
// 		if errors.As(err, &opErr) {
// 			// unwrap to syscall errno if present
// 			if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
// 				if errno, ok := syscallErr.Err.(syscall.Errno); ok {
// 					switch errno {
// 					case syscall.ECONNREFUSED:
// 						// host is up, port closed
// 						result.Status = scan.StatusUp
// 						result.Meta["port"] = port
// 						result.Meta["rtt_ms"] = rttMs
// 						result.Meta["signal"] = "connect_refused"
// 						return result, nil
//
// 					case syscall.ENETUNREACH, syscall.EHOSTUNREACH:
// 						// definitely not reachable: treat as down and bail
// 						result.Status = scan.StatusDown
// 						result.Meta["signal"] = "unreachable"
// 						result.Meta["error"] = err.Error()
// 						return result, nil
// 					}
// 				}
// 			}
//
// 			// timeout: try next port
// 			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
// 				// Only set signal if nothing else set; we'll keep probing other ports.
// 				if _, ok := result.Meta["signal"]; !ok {
// 					result.Meta["signal"] = "timeout"
// 				}
// 				continue
// 			}
//
// 			// other net error: remember it, but keep scanning
// 			result.Meta["error"] = err.Error()
// 			continue
// 		}
//
// 		// Non-net.OpError: most likely context cancellations etc.
// 		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
// 			return result, err
// 		}
//
// 		// Fallback: some other error; note it but don't claim host up
// 		result.Meta["error"] = err.Error()
// 	}
//
// 	// No port gave us a definitive "up" signal.
// 	if result.Status == scan.StatusUnknown {
// 		result.Status = scan.StatusDown
// 		if _, ok := result.Meta["signal"]; !ok {
// 			result.Meta["signal"] = "no_response"
// 		}
// 	}
// 	return result, nil
// }
