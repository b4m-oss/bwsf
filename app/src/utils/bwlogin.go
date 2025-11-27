package utils

import (
	"fmt"
	"os"
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

	// Build command arguments
	args := []string{"login", email, password}
	if serverURL != "" {
		args = append(args, "--server", serverURL)
	}

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
