package cmd

import (
	"bwenv/src/config"
	"bwenv/src/utils"
	"errors"
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
		fmt.Fprintf(os.Stderr, "[ERROR] ❌ bw command is not installed...\n")
		os.Exit(1)
	}

	// Get dotenvs folder ID (with unlock retry if locked)
	folderID, err := getDotenvsFolderIDWithUnlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
		os.Exit(1)
	}

	// List items in the folder (with unlock retry if locked)
	items, err := listItemsInFolderWithUnlock(folderID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
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

// getDotenvsFolderIDWithUnlock attempts to get folder ID, and unlocks if locked
func getDotenvsFolderIDWithUnlock() (string, error) {
	folderID, err := utils.GetDotenvsFolderID()
	if err != nil {
		// Check if error is due to locked state
		if errors.Is(err, utils.ErrBitwardenLocked) || utils.IsLockedError(err) {
			// Try to use login instead of unlock if unlock fails
			// Load config to get email and server URL
			cfg, configErr := config.LoadConfig()
			if configErr == nil && cfg != nil {
				// Prompt for master password
				fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
				masterPassword, err := utils.InputPassword()
				if err != nil {
					return "", fmt.Errorf("failed to get master password: %w", err)
				}

				// Try unlock first
				success, errorMsg := utils.BwUnlock(masterPassword)
				if !success {
					// If unlock fails, try login first, then unlock
					fmt.Println("[INFO] Unlock failed, trying login then unlock...")
					loginSuccess, loginErrorMsg := utils.BwLogin(cfg.Email, masterPassword, cfg.SelfhostedURL)
					if !loginSuccess {
						return "", fmt.Errorf("failed to login Bitwarden CLI: %s", loginErrorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI logged in successfully")
					
					// After login, try unlock again
					success, errorMsg = utils.BwUnlock(masterPassword)
					if !success {
						return "", fmt.Errorf("failed to unlock Bitwarden CLI after login: %s", errorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				} else {
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				}

				// Retry getting folder ID
				folderID, err = utils.GetDotenvsFolderID()
				if err != nil {
					return "", fmt.Errorf("failed to list folders after unlock/login: %w", err)
				}
				return folderID, nil
			}

			// If config is not available, just try unlock
			fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
			masterPassword, err := utils.InputPassword()
			if err != nil {
				return "", fmt.Errorf("failed to get master password: %w", err)
			}

			// Attempt to unlock
			success, errorMsg := utils.BwUnlock(masterPassword)
			if !success {
				return "", fmt.Errorf("failed to unlock Bitwarden CLI: %s", errorMsg)
			}

			fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")

			// Retry getting folder ID
			folderID, err = utils.GetDotenvsFolderID()
			if err != nil {
				return "", fmt.Errorf("failed to list folders after unlock: %w", err)
			}
			return folderID, nil
		}
		return "", err
	}
	return folderID, nil
}

// listItemsInFolderWithUnlock attempts to list items, and unlocks if locked
func listItemsInFolderWithUnlock(folderID string) ([]utils.Item, error) {
	items, err := utils.ListItemsInFolder(folderID)
	if err != nil {
		// Check if error is due to locked state
		if errors.Is(err, utils.ErrBitwardenLocked) || utils.IsLockedError(err) {
			// Try to use login instead of unlock if unlock fails
			// Load config to get email and server URL
			cfg, configErr := config.LoadConfig()
			if configErr == nil && cfg != nil {
				// Prompt for master password
				fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
				masterPassword, err := utils.InputPassword()
				if err != nil {
					return nil, fmt.Errorf("failed to get master password: %w", err)
				}

				// Try unlock first
				success, errorMsg := utils.BwUnlock(masterPassword)
				if !success {
					// If unlock fails, try login first, then unlock
					fmt.Println("[INFO] Unlock failed, trying login then unlock...")
					loginSuccess, loginErrorMsg := utils.BwLogin(cfg.Email, masterPassword, cfg.SelfhostedURL)
					if !loginSuccess {
						return nil, fmt.Errorf("failed to login Bitwarden CLI: %s", loginErrorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI logged in successfully")
					
					// After login, try unlock again
					success, errorMsg = utils.BwUnlock(masterPassword)
					if !success {
						return nil, fmt.Errorf("failed to unlock Bitwarden CLI after login: %s", errorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				} else {
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				}

				// Retry listing items
				items, err = utils.ListItemsInFolder(folderID)
				if err != nil {
					return nil, fmt.Errorf("failed to list items after unlock/login: %w", err)
				}
				return items, nil
			}

			// If config is not available, just try unlock
			fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
			masterPassword, err := utils.InputPassword()
			if err != nil {
				return nil, fmt.Errorf("failed to get master password: %w", err)
			}

			// Attempt to unlock
			success, errorMsg := utils.BwUnlock(masterPassword)
			if !success {
				return nil, fmt.Errorf("failed to unlock Bitwarden CLI: %s", errorMsg)
			}

			fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")

			// Retry listing items
			items, err = utils.ListItemsInFolder(folderID)
			if err != nil {
				return nil, fmt.Errorf("failed to list items after unlock: %w", err)
			}
			return items, nil
		}
		return nil, err
	}
	return items, nil
}

