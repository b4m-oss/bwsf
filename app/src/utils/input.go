package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// SelectHostType prompts user to select between Cloud and Self-hosted
func SelectHostType() (string, error) {
	prompt := promptui.Select{
		Label: "Bitwarden Cloud or Self-hosted?",
		Items: []string{"cloud", "selfhosted"},
	}

	index, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("failed to select host type: %w", err)
	}

	_ = index // index is not used but returned by prompt.Run()
	return result, nil
}

// InputURL prompts user to enter self-hosted URL
func InputURL() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	Question("Enter self-hosted URL: ")
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

	Question("Enter email address: ")
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
	Question("Enter password: ")

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
