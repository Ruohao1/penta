package hosts

import (
	"context"

	"github.com/Ruohao1/penta/internal/model"
)

type Prober interface {
	Name() string
	Probe(ctx context.Context, target model.Target, opts model.RunOptions) (model.Finding, error)
}
