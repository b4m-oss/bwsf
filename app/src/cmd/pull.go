package cmd

import (
	"bwenv/src/config"
	"bwenv/src/core"
	"bwenv/src/infra"
	"bwenv/src/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull .env file from Bitwarden",
	Long:  "Pull .env file from Bitwarden and save it to the current directory or specified directory",
	Run:   runPull,
}

func init() {
	pullCmd.Flags().String("output", ".", "Directory to save .env file")
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) {
	// Check if bw command is installed
	installed, _ := utils.CheckBwCommand()
	if !installed {
		utils.Errorln("[ERROR] ❌ bw command is not installed...")
		os.Exit(1)
	}

	// Get --output flag value
	outputDir, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.Errorln("[ERROR] Failed to get --output flag:", err)
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

	// confirmOverwrite wrapper
	confirmOverwrite := func(path string) (bool, error) {
		return utils.ConfirmOverwrite(fmt.Sprintf(".env file already exists in %s. Overwrite? (y/N): ", path))
	}

	// Call core logic
	err = core.PullEnvCore(
		outputDir,
		projectName,
		fs,
		bw,
		cfg,
		utils.InputPassword,
		confirmOverwrite,
		logger,
	)
	if err != nil {
		utils.Errorln("[ERROR]", err)
		os.Exit(1)
	}

	utils.Successln("[INFO] ✅ .env file pulled successfully!")
}
