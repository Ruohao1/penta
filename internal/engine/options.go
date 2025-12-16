package engine

import (
	"time"

	"github.com/Ruohao1/penta/internal/targets"
)

type RunOptions struct {
	// Target selection
	Targets []targets.Target // parsed targets
	Scope   []string         // allow-list enforcement

	// Scanning behavior
	TCPPorts           []int // explicit port list
	EnableDefaultPorts bool  // fallback (22,80,443...) if no ports provided
	ICMP               bool  // do ICMP reachability test
	ARP                bool  // do ARP scan (LAN-only)
	HTTP               bool  // enable HTTP probes (title, server)
	TLS                bool  // enable TLS fingerprinting

	// Performance + limits
	Concurrency int // worker pool size
	MinRate     int
	MaxRate     int
	MaxRetries  int // retry failed probes

	Timeout     time.Duration
	TimeoutTCP  time.Duration // per-dial timeout
	TimeoutHTTP time.Duration // per HTTP request
	TimeoutTLS  time.Duration // TLS handshake deadline

	// OpSec + networking
	Jitter     time.Duration // random delay ± jitter per request
	UserAgent  string        // custom UA for HTTP probes
	Proxy      string        // socks5/http proxy optional
	DNSServers []string      // optional custom resolvers
}
