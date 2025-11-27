package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// SelectHostType prompts user to select between Cloud and Self-hosted
func SelectHostType() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Bitwarden Cloud or Self-hosted? (cloud/selfhosted): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "cloud" || input == "selfhosted" {
			return input, nil
		}

		fmt.Println("Invalid input. Please enter 'cloud' or 'selfhosted'")
	}
}

// InputURL prompts user to enter self-hosted URL
func InputURL() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter self-hosted URL: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	url := strings.TrimSpace(input)
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	return url, nil
}

// InputEmail prompts user to enter email address
func InputEmail() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter email address: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	email := strings.TrimSpace(input)
	if email == "" {
		return "", fmt.Errorf("email cannot be empty")
	}

	return email, nil
}

// InputPassword prompts user to enter password (hidden input)
func InputPassword() (string, error) {
	fmt.Print("Enter password: ")

	// Read password without echoing to terminal
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	fmt.Println() // Print newline after password input

	password := string(passwordBytes)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}
