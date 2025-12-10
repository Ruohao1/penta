package engine

import "time"

type RunOptions struct {
	// Target selection
	TargetsExpr string   // raw targets (IPs, CIDRs, domains)
	Scope       []string // allow-list enforcement

	// Scanning behavior
	TCPPorts           []int // explicit port list
	EnableDefaultPorts bool  // fallback (22,80,443...) if no ports provided
	ICMP               bool  // do ICMP reachability test
	ARP                bool  // do ARP scan (LAN-only)
	HTTP               bool  // enable HTTP probes (title, server)
	TLS                bool  // enable TLS fingerprinting

	// Performance + limits
	Concurrency int           // worker pool size
	RateLimit   int           // global dial limit per second
	TimeoutTCP  time.Duration // per-dial timeout
	TimeoutHTTP time.Duration // per HTTP request
	TimeoutTLS  time.Duration // TLS handshake deadline
	Retries     int           // retry failed probes

	// OpSec + networking
	Jitter     time.Duration // random delay ± jitter per request
	UserAgent  string        // custom UA for HTTP probes
	Proxy      string        // socks5/http proxy optional
	DNSServers []string      // optional custom resolvers
}
