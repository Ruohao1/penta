package runner

import "sync"

// HostGate enforces at most N concurrent operations per host key.
type HostGate struct {
	mu    sync.Mutex
	sem   map[string]chan struct{}
	limit int
}

func NewHostGate(limit int) *HostGate {
	if limit < 1 {
		limit = 1
	}
	return &HostGate{
		sem:   make(map[string]chan struct{}),
		limit: limit,
	}
}

func (g *HostGate) Acquire(hostKey string) func() {
	g.mu.Lock()
	ch, ok := g.sem[hostKey]
	if !ok {
		ch = make(chan struct{}, g.limit)
		g.sem[hostKey] = ch
	}
	g.mu.Unlock()

	ch <- struct{}{}
	return func() { <-ch }
}
