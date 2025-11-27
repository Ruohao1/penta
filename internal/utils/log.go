package utils

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(human bool, level zerolog.Level) zerolog.Logger {
	var w io.Writer = os.Stderr
	if human {
		w = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	}

	logger := zerolog.New(w).
		Level(level).
		With().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(level)

	return logger
}
