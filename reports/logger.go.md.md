Your `logger` package provides a simple way to log messages with different severity levels (Info, Warn, and Error) using ANSI color codes for better visibility in the terminal. Here are a few suggestions to enhance this package:

1. **Timestamping**: Adding a timestamp to each log message can be very useful for tracking when events occurred.
2. **Log Levels**: Allow the user to set a log level to control which messages get printed.
3. **Custom Output**: Allow the user to specify where the log messages should be written (e.g., a file, standard output, etc.).
4. **Package Documentation**: Adding comments and documentation to your functions and constants can help others understand and use your package effectively.

Here's an improved version of your `logger` package with these enhancements:

```go
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
)

// LogLevel defines the severity of the log message
type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

// Logger represents a simple logger
type Logger struct {
	level  LogLevel
	output io.Writer
}

// New creates a new Logger
func New(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// Info logs informational messages
func (l *Logger) Info(message string) {
	if l.level <= INFO {
		l.log(ColorGreen, "INF", message)
	}
}

// Warn logs warning messages
func (l *Logger) Warn(message string) {
	if l.level <= WARN {
		l.log(ColorYellow, "WRN", message)
	}
}

// Error logs error messages
func (l *Logger) Error(message string) {
	if l.level <= ERROR {
		l.log(ColorRed, "ERR", message)
	}
}

// log formats and outputs a log message
func (l *Logger) log(color, level, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(l.output, "%s[%s] [%s] %s%s\n", color, level, timestamp, message, ColorReset)
}

// Default logger instance
var defaultLogger = New(INFO, os.Stdout)

// Info logs informational messages using the default logger
func Info(message string) {
	defaultLogger.Info(message)
}

// Warn logs warning messages using the default logger
func Warn(message string) {
	defaultLogger.Warn(message)
}

// Error logs error messages using the default logger
func Error(message string) {
	defaultLogger.Error(message)
}
```

### Usage Example

```go
package main

import (
	"os"
	"github.com/yourusername/logger"
)

func main() {
	// Use default logger
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Create a custom logger
	file, _ := os.Create("app.log")
	defer file.Close()
	customLogger := logger.New(logger.WARN, file)

	customLogger.Info("This info message won't be logged because the level is set to WARN")
	customLogger.Warn("This is a custom warning message")
	customLogger.Error("This is a custom error message")
}
```

This new version of the logger package includes timestamping, log levels, and the ability to specify custom output destinations. It also provides a default logger for convenience.