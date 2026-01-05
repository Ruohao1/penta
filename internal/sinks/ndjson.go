package sinks

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/Ruohao1/penta/internal/model"
)

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
