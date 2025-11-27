package cmd

import (
	"bwenv/src/core"
	"bwenv/src/infra"
	"bwenv/src/utils"
	"os"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Bitwarden host configuration",
	Long:  "Configure Bitwarden host (Cloud or Self-hosted) and login credentials",
	Run:   runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	// Check if bw command is installed
	installed, _ := utils.CheckBwCommand()
	if !installed {
		utils.Errorln("[ERROR] ❌ bw command is not installed...")
		os.Exit(1)
	}

	// Create dependencies
	bw := infra.NewBwClient()
	fs := infra.NewFileSystem()
	logger := infra.NewLogger()

	// confirmCreateFolder wrapper
	confirmCreateFolder := func() (bool, error) {
		return utils.ConfirmYesNo("dotenvs folder not found. Create it? (y/N): ")
	}

	// Call core logic
	err := core.SetupBitwardenCore(
		fs,
		bw,
		logger,
		utils.SelectHostType,
		utils.InputURL,
		utils.InputEmail,
		utils.InputPassword,
		confirmCreateFolder,
	)
	if err != nil {
		utils.Errorln("[ERROR]", err)
		os.Exit(1)
	}

	// Success message
	utils.Successln("[INFO] ✅ Sign in to Bitwarden was successful!")
}
