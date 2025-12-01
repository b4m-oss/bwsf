package utils

import (
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

var (
	currentSpinner *spinner.Spinner
	spinnerMutex   sync.Mutex
)

// isSpinnerEnabled checks if spinner should be enabled
// Spinner is disabled in non-TTY environments or when NO_COLOR is set
func isSpinnerEnabled() bool {
	// Check NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout is a terminal
	return isTerminal(os.Stdout)
}

// StartSpinner starts a spinner with the given message
// Returns true if spinner was started, false if disabled
func StartSpinner(message string) bool {
	if !isSpinnerEnabled() {
		return false
	}

	spinnerMutex.Lock()
	defer spinnerMutex.Unlock()

	// Stop any existing spinner
	if currentSpinner != nil {
		currentSpinner.Stop()
	}

	// Create new spinner with dots style (CharSets[14]: ⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Writer = os.Stderr
	s.Start()

	currentSpinner = s
	return true
}

// StopSpinner stops the current spinner (success case)
func StopSpinner() {
	spinnerMutex.Lock()
	defer spinnerMutex.Unlock()

	if currentSpinner != nil {
		currentSpinner.Stop()
		currentSpinner = nil
	}
}

// StopSpinnerWithError stops the current spinner (error case)
func StopSpinnerWithError() {
	spinnerMutex.Lock()
	defer spinnerMutex.Unlock()

	if currentSpinner != nil {
		currentSpinner.Stop()
		currentSpinner = nil
	}
}

// UpdateSpinnerMessage updates the message of the current spinner
func UpdateSpinnerMessage(message string) {
	spinnerMutex.Lock()
	defer spinnerMutex.Unlock()

	if currentSpinner != nil {
		currentSpinner.Suffix = " " + message
	}
}



