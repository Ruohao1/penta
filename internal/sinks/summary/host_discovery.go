package summary

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/sinks"
)

type HostDiscoverySink struct {
	next sinks.Sink
	out  io.Writer

	mu    sync.Mutex
	up    []string
	seen  map[string]bool
	stage string
}

func NewHostDiscoverySink(next sinks.Sink, out io.Writer) *HostDiscoverySink {
	return &HostDiscoverySink{
		next: next,
		out:  out,
		seen: make(map[string]bool),
	}
}

func (s *HostDiscoverySink) Emit(ctx context.Context, ev model.Event) {
	// keep streaming to next sink if you want (or skip for quiet mode)
	// s.next.Emit(ctx, ev)
	fmt.Println(ev)

	if ev.Type != model.EventHostState || ev.HostState == nil {
		return
	}
	if ev.HostState.State != model.HostStateUp {
		return
	}
	fmt.Println("here")

	host := ev.HostState.Host
	if host == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen[host] {
		return
	}
	s.seen[host] = true
	s.up = append(s.up, host)
	if s.stage == "" {
		s.stage = ev.Stage
	}
}

func (s *HostDiscoverySink) Close() error {
	// close downstream first (flushes files, etc.)
	_ = s.next.Close()

	s.mu.Lock()
	defer s.mu.Unlock()

	sort.Strings(s.up)

	// nice compact output
	if len(s.up) == 0 {
		_, _ = io.WriteString(s.out, "No hosts up.\n")
		return nil
	}

	_, _ = fmt.Fprintf(s.out, "\nHosts up (%d) [stage=%s]\n", len(s.up), s.stage)
	for _, h := range s.up {
		_, _ = fmt.Fprintf(s.out, "  %s\n", h)
	}
	return nil
}
