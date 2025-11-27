package cmd

import (
	"bwenv/src/utils"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "bwenv",
	Short:   "CLI tool to manage .env files using Bitwarden",
	Long:    "bwenv is a CLI tool that uses Bitwarden to manage .env files",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Error("Error: %v\n", err)
		os.Exit(1)
	}
}
