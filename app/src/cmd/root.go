package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bwenv",
	Short: "CLI tool to manage .env files using Bitwarden",
	Long:  "bwenv is a CLI tool that uses Bitwarden to manage .env files",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
