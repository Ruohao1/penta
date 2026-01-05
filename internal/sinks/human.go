package sinks

import (
	"context"
	"io"

	"github.com/Ruohao1/penta/internal/model"
)

type HumanSink struct {
	out, err io.Writer
}

func NewHumanSink(out, err io.Writer) *HumanSink { return &HumanSink{out: out, err: err} }

func (s *HumanSink) Emit(ctx context.Context, ev model.Event) {
	line := ev.String() + "\n"
	switch ev.Type {
	case model.EventUnknown, model.EventError:
		s.err.Write([]byte(line))
	default:
		s.out.Write([]byte(line))
	}
}

func (s *HumanSink) Close() error { return nil }
