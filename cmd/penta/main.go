package main

import (
	"os"

	"github.com/Ruohao1/penta/internal/app"
)

func main() {
	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}
