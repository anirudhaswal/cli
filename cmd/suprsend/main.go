/*
Copyright © 2025 SuprSend
*/
package main

import (
	_ "embed"
	"os"

	"github.com/suprsend/cli/internal/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
