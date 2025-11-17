package main

import (
	"github.com/Ruohao1/penta/internal/app"
)

func main() {
	if err := app.Execute(); err != nil {
		panic(err)
	}
}
