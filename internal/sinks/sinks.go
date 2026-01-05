package sinks

import (
	"context"
	"io"

	"github.com/Ruohao1/penta/internal/model"
)

type Sink interface {
	Emit(context.Context, model.Event)
	Close() error
}

type SinkOptions struct {
	JSON bool
	Out  io.Writer
	Err  io.Writer
}

func NewSink(opt SinkOptions) Sink {
	if opt.JSON {
		return NewNDJSONSink(opt.Out)
	}
	return NewHumanSink(opt.Out, opt.Err)
}
