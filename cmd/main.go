package main

import (
	"os"

	"fin.orion.dev/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
