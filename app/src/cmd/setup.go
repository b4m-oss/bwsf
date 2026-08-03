package cmd

import (
	"bwsf/src/config"
	"bwsf/src/core"
	"bwsf/src/infra"
	"bwsf/src/utils"
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
	cfg := loadConfigOrEmpty()
	if cfg.GetBackend() == config.BackendAPI {
		utils.Infoln("[INFO] backend=api: use `bwsf auth` to store a Personal API Key and obtain an Identity token.")
		utils.Infoln("[INFO] Continuing setup will still configure host/email for Identity URL resolution.")
	}
	requireBwCLIIfNeeded(cfg)

	// Create dependencies
	bw := newBwClientFromConfig(cfg)
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
