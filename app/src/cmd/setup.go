package cmd

import (
	"bwenv/src/config"
	"bwenv/src/utils"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "[ERROR] ❌ bw command is not installed...\n")
		os.Exit(1)
	}

	// Load existing config (if any)
	existingConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to load existing config: %v\n", err)
		os.Exit(1)
	}

	// If config exists, inform user that it will be overwritten
	if existingConfig != nil {
		fmt.Println("[INFO] Existing configuration found. It will be overwritten.")
	}

	// Step 1: Select host type
	hostType, err := utils.SelectHostType()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to select host type: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Get self-hosted URL if needed
	var selfhostedURL string
	if hostType == "selfhosted" {
		selfhostedURL, err = utils.InputURL()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR] Failed to get URL: %v\n", err)
			os.Exit(1)
		}
	}

	// Step 3: Get email
	email, err := utils.InputEmail()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to get email: %v\n", err)
		os.Exit(1)
	}

	// Step 4: Get password
	password, err := utils.InputPassword()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to get password: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Attempt login with bw command
	success, errorMsg := utils.BwLogin(email, password, selfhostedURL)
	if !success {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", errorMsg)
		os.Exit(1)
	}

	// Step 6: Save configuration
	newConfig := &config.Config{
		HostType:      hostType,
		SelfhostedURL: selfhostedURL,
		Email:         email,
	}

	if err := config.SaveConfig(newConfig); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to save configuration: %v\n", err)
		os.Exit(1)
	}

	// Step 7: Success message
	fmt.Println("[INFO] ✅ Sign in to Bitwarden was successful!")
}
