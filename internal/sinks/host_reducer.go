package sinks

import (
	"context"
	"sync"
	"time"

	"github.com/Ruohao1/penta/internal/model"
)

type HostReducerSink struct {
	next Sink

	mu sync.Mutex
	up map[string]bool // key: host/ip string
}

func NewHostReducerSink(next Sink) *HostReducerSink {
	return &HostReducerSink{
		next: next,
		up:   make(map[string]bool),
	}
}

func (s *HostReducerSink) Emit(ctx context.Context, ev model.Event) {
	// Always forward the event (or not, depending on how you compose)
	// s.next.Emit(ctx, ev)

	if ev.Type != model.EventFinding || ev.Finding == nil {
		return
	}

	f := ev.Finding
	host := f.Endpoint.Key()
	if host == "" {
		return
	}

	ok, _ := f.Meta["ok"].(bool)
	isUpEvidence := ok || f.Status == "refused"

	if !isUpEvidence {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.up[host] {
		return // already marked up
	}
	s.up[host] = true

	hs := model.HostStateEvent{
		Host:  host,
		State: model.HostStateUp,
		// Reason: "tcp_connect"
	}
	s.next.Emit(ctx, model.Event{
		EmittedAt: time.Now().UTC(),
		Type:      model.EventHostState,
		Stage:     ev.Stage,
		HostState: &hs,
	})
}

func (s *HostReducerSink) Close() error { return s.next.Close() }
