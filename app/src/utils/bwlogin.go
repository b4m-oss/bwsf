package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// BwLogin executes bw login command and returns success status and error message
func BwLogin(email, password, serverURL string) (bool, string) {
	// Check if bw command exists
	_, err := exec.LookPath("bw")
	if err != nil {
		return false, "bw command is not installed"
	}

	// Start spinner
	StartSpinner("Logging in...")
	defer StopSpinner()

	// If self-hosted, configure server first
	// But first check if we need to logout
	if serverURL != "" {
		// Check current server config
		checkCmd := exec.Command("bw", "config", "server")
		checkOutput, _ := checkCmd.CombinedOutput()
		currentServer := strings.TrimSpace(string(checkOutput))

		// If server URL is different, logout first
		if currentServer != "" && currentServer != serverURL {
			logoutCmd := exec.Command("bw", "logout")
			logoutCmd.Run() // Ignore errors, just try to logout
		}

		configCmd := exec.Command("bw", "config", "server", serverURL)
		configOutput, err := configCmd.CombinedOutput()
		if err != nil {
			errorMsg := strings.TrimSpace(string(configOutput))
			if errorMsg == "" {
				errorMsg = err.Error()
			}
			// If error is about logout required, try logout and retry
			if strings.Contains(errorMsg, "Logout required") {
				logoutCmd := exec.Command("bw", "logout")
				logoutCmd.Run() // Ignore errors
				// Retry config
				configOutput, err = configCmd.CombinedOutput()
				if err != nil {
					errorMsg = strings.TrimSpace(string(configOutput))
					if errorMsg == "" {
						errorMsg = err.Error()
					}
					return false, fmt.Sprintf("Failed to configure server: %s", errorMsg)
				}
			} else {
				return false, fmt.Sprintf("Failed to configure server: %s", errorMsg)
			}
		}
	}

	// Build login command arguments
	args := []string{"login", email, password}

	// Execute bw login command
	cmd := exec.Command("bw", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Extract error message from output
		errorMsg := strings.TrimSpace(string(output))
		if errorMsg == "" {
			errorMsg = err.Error()
		}
		return false, errorMsg
	}

	// Check if login was successful
	outputStr := strings.TrimSpace(string(output))
	if strings.Contains(outputStr, "You are logged in") ||
		strings.Contains(outputStr, "You're logged in") ||
		strings.Contains(outputStr, "logged in!") {
		return true, ""
	}

	// If output doesn't contain success message, treat as failure
	return false, outputStr
}

// CheckBwCommand checks if bw command is installed
func CheckBwCommand() (bool, string) {
	path, err := exec.LookPath("bw")
	if err != nil {
		return false, ""
	}
	return true, path
}
