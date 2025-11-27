package cmd

import (
	"bwenv/src/config"
	"bwenv/src/core"
	"bwenv/src/infra"
	"bwenv/src/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List items in the dotenvs folder",
	Long:  "List all items in the dotenvs folder from Bitwarden",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	// Check if bw command is installed
	installed, _ := utils.CheckBwCommand()
	if !installed {
		utils.Errorln("[ERROR] ❌ bw command is not installed...")
		os.Exit(1)
	}

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
	logger := infra.NewLogger()

	// Call core logic
	items, err := core.ListDotenvsCore(
		bw,
		cfg,
		utils.InputPassword,
		logger,
	)
	if err != nil {
		utils.Errorln("[ERROR]", err)
		os.Exit(1)
	}

	// Output item names (one per line)
	if len(items) == 0 {
		fmt.Println("No items found in dotenvs folder")
		return
	}

	for _, item := range items {
		fmt.Println(item.Name)
	}
}
