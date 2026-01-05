package sinks

import (
	"context"

	"github.com/Ruohao1/penta/internal/model"
)

type FilterFunc func(model.Event) bool

type FilterSink struct {
	Filter FilterFunc
	Next   Sink
}

func NewFilterSink(next Sink, filter FilterFunc) Sink {
	return &FilterSink{Filter: filter, Next: next}
}

func (s *FilterSink) Emit(ctx context.Context, ev model.Event) {
	if s.Filter != nil && s.Filter(ev) {
		s.Next.Emit(ctx, ev)
	}
}

func (s *FilterSink) Close() error { return s.Next.Close() }

func OnlyOKFindings(ev model.Event) bool {
	if ev.Type != model.EventFinding || ev.Finding == nil {
		return true
	}
	b, ok := ev.Finding.Meta["ok"].(bool)
	return ok && b
}

func OnlyOpen(ev model.Event) bool {
	if ev.Type != model.EventFinding || ev.Finding == nil {
		return true
	}
	return ev.Finding.Status == "open" || ev.Finding.Status == "ok"
}
