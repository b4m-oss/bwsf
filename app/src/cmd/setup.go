package cmd

import (
	"fmt"
	"os"
	"strings"

	"bwsf/src/config"
	"bwsf/src/core"
	"bwsf/src/infra"
	"bwsf/src/utils"

	"github.com/spf13/cobra"
)

var setupFolder string

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Bitwarden host configuration",
	Long:  "Configure Bitwarden host (Cloud or Self-hosted) and login credentials",
	Run:   runSetup,
}

func init() {
	setupCmd.Flags().StringVar(&setupFolder, "folder", "", "Bitwarden folder name for .env notes (default: dotenvs)")
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	// Check if bw command is installed
	installed, _ := utils.CheckBwCommand()
	if !installed {
		utils.Errorln("[ERROR] ❌ bw command is not installed...")
		os.Exit(1)
	}

	folderName := config.DefaultFolderName

	// Persist --folder before core setup so RealBwClient reads it for folder ops.
	if setupFolder != "" {
		if err := config.ValidateFolderName(setupFolder); err != nil {
			utils.Errorln("[ERROR]", err)
			os.Exit(1)
		}
		folderName = strings.TrimSpace(setupFolder)

		cfg, err := config.LoadConfig()
		if err != nil {
			utils.Errorln("[ERROR] Failed to load config:", err)
			os.Exit(1)
		}
		if cfg == nil {
			cfg = &config.Config{}
		}
		cfg.FolderName = folderName
		if err := config.SaveConfig(cfg); err != nil {
			utils.Errorln("[ERROR] Failed to save folder name:", err)
			os.Exit(1)
		}
	} else {
		cfg, err := config.LoadConfig()
		if err != nil {
			utils.Errorln("[ERROR] Failed to load config:", err)
			os.Exit(1)
		}
		folderName = config.ResolveFolderName(cfg)
	}

	// Create dependencies
	bw := infra.NewBwClient()
	fs := infra.NewFileSystem()
	logger := infra.NewLogger()

	// confirmCreateFolder wrapper
	confirmCreateFolder := func() (bool, error) {
		return utils.ConfirmYesNo(fmt.Sprintf("%s folder not found. Create it? (y/N): ", folderName))
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
