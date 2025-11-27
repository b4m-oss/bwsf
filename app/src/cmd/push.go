package cmd

import (
	"bwenv/src/config"
	"bwenv/src/core"
	"bwenv/src/infra"
	"bwenv/src/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push .env file to Bitwarden",
	Long:  "Push .env file from specified directory to Bitwarden as a note item",
	Run:   runPush,
}

func init() {
	pushCmd.Flags().String("from", ".", "Directory containing .env file")
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) {
	// Check if bw command is installed
	installed, _ := utils.CheckBwCommand()
	if !installed {
		utils.Errorln("[ERROR] ❌ bw command is not installed...")
		os.Exit(1)
	}

	// Get --from flag value
	fromDir, err := cmd.Flags().GetString("from")
	if err != nil {
		utils.Errorln("[ERROR] Failed to get --from flag:", err)
		os.Exit(1)
	}

	// Get current working directory name as project name
	wd, err := os.Getwd()
	if err != nil {
		utils.Errorln("[ERROR] Failed to get current working directory:", err)
		os.Exit(1)
	}
	projectName := filepath.Base(wd)

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Errorln("[ERROR] Failed to load config:", err)
		os.Exit(1)
	}
	if cfg == nil {
		cfg = &config.Config{}
	}

	// Create dependencies
	bw := infra.NewBwClient()
	fs := infra.NewFileSystem()
	logger := infra.NewLogger()

	// Call core logic
	err = core.PushEnvCore(
		fromDir,
		projectName,
		fs,
		bw,
		cfg,
		utils.InputPassword,
		logger,
	)
	if err != nil {
		utils.Errorln("[ERROR]", err)
		os.Exit(1)
	}

	utils.Successln("[INFO] ✅ .env file pushed successfully!")
}
