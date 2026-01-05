package sinks

import (
	"context"

	"github.com/Ruohao1/penta/internal/model"
)

type MultiSink struct {
	sinks []Sink
}

func NewMultiSink(sinks ...Sink) *MultiSink {
	return &MultiSink{sinks: sinks}
}

func (m *MultiSink) Emit(ctx context.Context, ev model.Event) {
	for _, s := range m.sinks {
		s.Emit(ctx, ev)
	}
}

func (m *MultiSink) Close() error {
	var first error
	for _, s := range m.sinks {
		if err := s.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
