package main

import (
	"os"

	"github.com/Fire-Dragon-DoL/clickup-cli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
