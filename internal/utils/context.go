package utils

import (
	"context"

	"github.com/Ruohao1/penta/internal/config"
	"github.com/rs/zerolog"
)

type (
	ctxLoggerKey struct{}
	ctxConfigKey struct{}
)

func WithLogger(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

func LoggerFrom(ctx context.Context) zerolog.Logger {
	l, ok := ctx.Value(ctxLoggerKey{}).(zerolog.Logger)
	if !ok {
		return zerolog.Nop()
	}
	return l
}

func WithConfig(ctx context.Context, cfg *config.Config) context.Context {
	return context.WithValue(ctx, ctxConfigKey{}, cfg)
}

func ConfigFrom(ctx context.Context) *config.Config {
	cfg, _ := ctx.Value(ctxConfigKey{}).(*config.Config)
	return cfg
}
