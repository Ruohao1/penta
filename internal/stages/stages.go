package stages

import (
	"context"

	"github.com/Ruohao1/penta/internal/model"
	"github.com/Ruohao1/penta/internal/runner"
	"github.com/Ruohao1/penta/internal/sinks"
)

type Stage interface {
	Name() string
	Build(ctx context.Context, task model.Task, opts model.RunOptions, sink sinks.Sink) ([]runner.Job, error)
	// Optional: called after pool completes for this stage
	After(ctx context.Context, task model.Task, opts model.RunOptions, sink sinks.Sink) error
}
