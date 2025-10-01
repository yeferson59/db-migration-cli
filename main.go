package main

import (
	"os"

	"github.com/yeferson59/db-migration-cli/cmd/db-migration-cli/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
