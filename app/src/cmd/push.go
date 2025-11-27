package cmd

import (
	"bwenv/src/config"
	"bwenv/src/utils"
	"errors"
	"fmt"
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
		utils.Error("[ERROR] ❌ bw command is not installed...\n")
		os.Exit(1)
	}

	// Get --from flag value
	fromDir, err := cmd.Flags().GetString("from")
	if err != nil {
		utils.Error("[ERROR] Failed to get --from flag: %v\n", err)
		os.Exit(1)
	}

	// Get current working directory name as project name
	wd, err := os.Getwd()
	if err != nil {
		utils.Error("[ERROR] Failed to get current working directory: %v\n", err)
		os.Exit(1)
	}
	projectName := filepath.Base(wd)

	// Find .env file in the specified directory
	envPath, err := utils.FindEnvFile(fromDir)
	if err != nil {
		// Try /project-root if fromDir is "." or ".." (for git repository root)
		if fromDir == "." || fromDir == ".." {
			projectRootPath, rootErr := utils.FindEnvFile("/project-root")
			if rootErr == nil {
				envPath = projectRootPath
				err = nil
			}
		}
		if err != nil {
			utils.Error("[ERROR] %s\n", err)
			utils.Infoln("[INFO] Tip: Use --from flag to specify the directory containing .env file")
			utils.Infoln("[INFO] Example: bwenv push --from /project-root (to use project root)")
			os.Exit(1)
		}
	}

	// Parse .env file
	envData, err := utils.ParseEnvFile(envPath)
	if err != nil {
		utils.Error("[ERROR] Failed to parse .env file: %v\n", err)
		os.Exit(1)
	}

	// Convert to JSON
	jsonData, err := utils.EnvDataToJSON(envData)
	if err != nil {
		utils.Error("[ERROR] Failed to convert to JSON: %v\n", err)
		os.Exit(1)
	}

	// Get dotenvs folder ID (with unlock retry if locked)
	// Use the function from list.go
	folderID, err := getDotenvsFolderIDWithUnlock()
	if err != nil {
		utils.Error("[ERROR] %s\n", err)
		os.Exit(1)
	}

	// Check if item with project name already exists
	existingItem, err := getItemByNameWithUnlock(folderID, projectName)
	if err != nil {
		utils.Error("[ERROR] %s\n", err)
		os.Exit(1)
	}

	// If item exists, ask for confirmation
	if existingItem != nil {
		confirmed, err := utils.ConfirmOverwrite(fmt.Sprintf("Item '%s' already exists. Overwrite? (y/N): ", projectName))
		if err != nil {
			utils.Error("[ERROR] Failed to get confirmation: %v\n", err)
			os.Exit(1)
		}
		if !confirmed {
			utils.Infoln("[INFO] Operation cancelled.")
			return
		}

		// Update existing item
		err = updateItemWithUnlock(existingItem.ID, jsonData)
		if err != nil {
			utils.Error("[ERROR] Failed to update item: %s\n", err)
			os.Exit(1)
		}

		utils.Successln("[INFO] ✅ Item updated successfully!")
	} else {
		// Create new item
		err = createItemWithUnlock(folderID, projectName, jsonData)
		if err != nil {
			utils.Error("[ERROR] Failed to create item: %s\n", err)
			os.Exit(1)
		}

		utils.Successln("[INFO] ✅ Item created successfully!")
	}
}

// getItemByNameWithUnlock attempts to get item by name, and unlocks if locked
func getItemByNameWithUnlock(folderID, itemName string) (*utils.FullItem, error) {
	item, err := utils.GetItemByName(folderID, itemName)
	if err != nil {
		// Check if error is due to locked state
		if errors.Is(err, utils.ErrBitwardenLocked) || utils.IsLockedError(err) {
			// Try to unlock
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

				// Retry getting item
				item, err = utils.GetItemByName(folderID, itemName)
				if err != nil {
					return nil, fmt.Errorf("failed to get item after unlock/login: %w", err)
				}
				return item, nil
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

			// Retry getting item
			item, err = utils.GetItemByName(folderID, itemName)
			if err != nil {
				return nil, fmt.Errorf("failed to get item after unlock: %w", err)
			}
			return item, nil
		}
		return nil, err
	}
	return item, nil
}

// createItemWithUnlock attempts to create item, and unlocks if locked
func createItemWithUnlock(folderID, name, notes string) error {
	err := utils.CreateNoteItem(folderID, name, notes)
	if err != nil {
		// Check if error is due to locked state
		if errors.Is(err, utils.ErrBitwardenLocked) || utils.IsLockedError(err) {
			// Try to unlock
			cfg, configErr := config.LoadConfig()
			if configErr == nil && cfg != nil {
				// Prompt for master password
				fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
				masterPassword, err := utils.InputPassword()
				if err != nil {
					return fmt.Errorf("failed to get master password: %w", err)
				}

				// Try unlock first
				success, errorMsg := utils.BwUnlock(masterPassword)
				if !success {
					// If unlock fails, try login first, then unlock
					fmt.Println("[INFO] Unlock failed, trying login then unlock...")
					loginSuccess, loginErrorMsg := utils.BwLogin(cfg.Email, masterPassword, cfg.SelfhostedURL)
					if !loginSuccess {
						return fmt.Errorf("failed to login Bitwarden CLI: %s", loginErrorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI logged in successfully")

					// After login, try unlock again
					success, errorMsg = utils.BwUnlock(masterPassword)
					if !success {
						return fmt.Errorf("failed to unlock Bitwarden CLI after login: %s", errorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				} else {
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				}

				// Retry creating item
				err = utils.CreateNoteItem(folderID, name, notes)
				if err != nil {
					return fmt.Errorf("failed to create item after unlock/login: %w", err)
				}
				return nil
			}

			// If config is not available, just try unlock
			fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
			masterPassword, err := utils.InputPassword()
			if err != nil {
				return fmt.Errorf("failed to get master password: %w", err)
			}

			// Attempt to unlock
			success, errorMsg := utils.BwUnlock(masterPassword)
			if !success {
				return fmt.Errorf("failed to unlock Bitwarden CLI: %s", errorMsg)
			}

			fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")

			// Retry creating item
			err = utils.CreateNoteItem(folderID, name, notes)
			if err != nil {
				return fmt.Errorf("failed to create item after unlock: %w", err)
			}
			return nil
		}
		return err
	}
	return nil
}

// updateItemWithUnlock attempts to update item, and unlocks if locked
func updateItemWithUnlock(itemID, notes string) error {
	err := utils.UpdateNoteItem(itemID, notes)
	if err != nil {
		// Check if error is due to locked state
		if errors.Is(err, utils.ErrBitwardenLocked) || utils.IsLockedError(err) {
			// Try to unlock
			cfg, configErr := config.LoadConfig()
			if configErr == nil && cfg != nil {
				// Prompt for master password
				fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
				masterPassword, err := utils.InputPassword()
				if err != nil {
					return fmt.Errorf("failed to get master password: %w", err)
				}

				// Try unlock first
				success, errorMsg := utils.BwUnlock(masterPassword)
				if !success {
					// If unlock fails, try login first, then unlock
					fmt.Println("[INFO] Unlock failed, trying login then unlock...")
					loginSuccess, loginErrorMsg := utils.BwLogin(cfg.Email, masterPassword, cfg.SelfhostedURL)
					if !loginSuccess {
						return fmt.Errorf("failed to login Bitwarden CLI: %s", loginErrorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI logged in successfully")

					// After login, try unlock again
					success, errorMsg = utils.BwUnlock(masterPassword)
					if !success {
						return fmt.Errorf("failed to unlock Bitwarden CLI after login: %s", errorMsg)
					}
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				} else {
					fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")
				}

				// Retry updating item
				err = utils.UpdateNoteItem(itemID, notes)
				if err != nil {
					return fmt.Errorf("failed to update item after unlock/login: %w", err)
				}
				return nil
			}

			// If config is not available, just try unlock
			fmt.Println("[INFO] Bitwarden CLI is locked. Please enter your master password to unlock.")
			masterPassword, err := utils.InputPassword()
			if err != nil {
				return fmt.Errorf("failed to get master password: %w", err)
			}

			// Attempt to unlock
			success, errorMsg := utils.BwUnlock(masterPassword)
			if !success {
				return fmt.Errorf("failed to unlock Bitwarden CLI: %s", errorMsg)
			}

			fmt.Println("[INFO] ✅ Bitwarden CLI unlocked successfully")

			// Retry updating item
			err = utils.UpdateNoteItem(itemID, notes)
			if err != nil {
				return fmt.Errorf("failed to update item after unlock: %w", err)
			}
			return nil
		}
		return err
	}
	return nil
}
