package logger

import "fmt"

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
)

// Info logs informational messages
func Info(message string) {
	fmt.Printf("%s[INF] %s%s\n", ColorGreen, message, ColorReset)
}

// Warn logs warning messages
func Warn(message string) {
	fmt.Printf("%s[WRN] %s%s\n", ColorYellow, message, ColorReset)
}

// Error logs error messages
func Error(message string) {
	fmt.Printf("%s[ERR] %s%s\n", ColorRed, message, ColorReset)
}
