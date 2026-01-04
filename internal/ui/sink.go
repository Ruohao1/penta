package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/Ruohao1/penta/internal/model"
)

type EventSink interface {
	Emit(context.Context, model.Event)
	Close() error
}

type SinkOptions struct {
	JSON bool
	Out  io.Writer
	Err  io.Writer
}

func NewSink(opt SinkOptions) EventSink {
	if opt.JSON {
		return NewNDJSONSink(opt.Out)
	}
	return NewHumanSink(opt.Out, opt.Err)
}

// MultiSink
type MultiSink struct {
	sinks []EventSink
}

func NewMultiSink(sinks ...EventSink) *MultiSink {
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

type NDJSONSink struct {
	mu sync.Mutex
	w  io.Writer
}

func NewNDJSONSink(w io.Writer) *NDJSONSink { return &NDJSONSink{w: w} }

func (s *NDJSONSink) Emit(ctx context.Context, ev model.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, err := json.Marshal(ev)
	if err != nil {
		return // or write an error event; don't panic in sink
	}
	s.w.Write(b)
	s.w.Write([]byte("\n"))
}

func (s *NDJSONSink) Close() error { return nil }

type HumanSink struct {
	out, err io.Writer
}

func NewHumanSink(out, err io.Writer) *HumanSink { return &HumanSink{out: out, err: err} }

func (s *HumanSink) Emit(ctx context.Context, ev model.Event) {
	fmt.Println(ev)
	// s.out.Write([]byte(ev.String() + "\n"))
}

func (s *HumanSink) Close() error { return nil }
