package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvData represents the parsed .env file data
// Lines array preserves the original order of all lines (comments, variables, empty lines)
type EnvData struct {
	Lines []string `json:"lines"`
}

// ParseEnvFile reads and parses a .env file
// It preserves the exact order of all lines (comments, variables, empty lines) and quotes in values
func ParseEnvFile(filePath string) (*EnvData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	data := &EnvData{
		Lines: []string{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Preserve the exact line as-is (including whitespace, quotes, etc.)
		data.Lines = append(data.Lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}

	return data, nil
}

// EnvDataToJSON converts EnvData to JSON string
func EnvDataToJSON(data *EnvData) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal env data to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

// RestoreEnvFileFromJSON restores .env file content from JSON string
func RestoreEnvFileFromJSON(jsonStr string) (string, error) {
	var data EnvData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Join all lines with newlines
	return strings.Join(data.Lines, "\n"), nil
}

// FindEnvFile looks for .env file in the specified directory
func FindEnvFile(dirPath string) (string, error) {
	envPath := filepath.Join(dirPath, ".env")

	// Check if file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return "", fmt.Errorf(".env file not found in %s", dirPath)
	}

	return envPath, nil
}
