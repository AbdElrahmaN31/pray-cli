// Package main is the entry point for the pray CLI application
package main

import (
	"fmt"
	"os"

	"github.com/anashaat/pray-cli/cmd/pray/cmd"
)

// Version information (set at build time)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info for the root command
	cmd.SetVersionInfo(version, commit, date)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
