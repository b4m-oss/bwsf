package utils

import (
	"fmt"
	"os"
)

const (
	// ANSI color codes
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

// isColorEnabled checks if color output should be enabled
func isColorEnabled() bool {
	// Check NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if stdout/stderr are terminals
	stdoutIsTerminal := isTerminal(os.Stdout)
	stderrIsTerminal := isTerminal(os.Stderr)

	return stdoutIsTerminal || stderrIsTerminal
}

// isTerminal checks if a file descriptor is a terminal
func isTerminal(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// colorize wraps text with ANSI color codes if color is enabled
func colorize(text string, colorCode string) string {
	if !isColorEnabled() {
		return text
	}
	return colorCode + text + colorReset
}

// Error prints an error message in red to stderr
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorize(message, colorRed)
	fmt.Fprint(os.Stderr, colored)
}

// Errorln prints an error message in red to stderr with a newline
func Errorln(args ...interface{}) {
	message := fmt.Sprint(args...)
	colored := colorize(message, colorRed)
	fmt.Fprintln(os.Stderr, colored)
}

// Success prints a success message in green to stdout
func Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorize(message, colorGreen)
	fmt.Print(colored)
}

// Successln prints a success message in green to stdout with a newline
func Successln(args ...interface{}) {
	message := fmt.Sprint(args...)
	colored := colorize(message, colorGreen)
	fmt.Println(colored)
}

// Warning prints a warning message in yellow to stdout
func Warning(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorize(message, colorYellow)
	fmt.Print(colored)
}

// Warningln prints a warning message in yellow to stdout with a newline
func Warningln(args ...interface{}) {
	message := fmt.Sprint(args...)
	colored := colorize(message, colorYellow)
	fmt.Println(colored)
}

// Info prints an info message in cyan to stdout
func Info(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorize(message, colorCyan)
	fmt.Print(colored)
}

// Infoln prints an info message in cyan to stdout with a newline
func Infoln(args ...interface{}) {
	message := fmt.Sprint(args...)
	colored := colorize(message, colorCyan)
	fmt.Println(colored)
}

// Question prints a question message in magenta to stdout
func Question(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	colored := colorize(message, colorMagenta)
	fmt.Print(colored)
}

// Questionln prints a question message in magenta to stdout with a newline
func Questionln(args ...interface{}) {
	message := fmt.Sprint(args...)
	colored := colorize(message, colorMagenta)
	fmt.Println(colored)
}

// ColorError returns a colored error string
func ColorError(text string) string {
	return colorize(text, colorRed)
}

// ColorSuccess returns a colored success string
func ColorSuccess(text string) string {
	return colorize(text, colorGreen)
}

// ColorWarning returns a colored warning string
func ColorWarning(text string) string {
	return colorize(text, colorYellow)
}

// ColorInfo returns a colored info string
func ColorInfo(text string) string {
	return colorize(text, colorCyan)
}

// ColorQuestion returns a colored question string
func ColorQuestion(text string) string {
	return colorize(text, colorMagenta)
}
