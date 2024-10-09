package config

import "fmt"

const (
	Yellow = "\033[33m" // Yellow color for normal information
	Red    = "\033[31m" // Red color for warnings and errors
	Reset  = "\033[0m"  // Reset the color
)

func Info(message string) {
	fmt.Printf("%s[INF] %s%s\n", Yellow, message, Reset)
}
func Warn(message string) {
	fmt.Printf("%s[WRN] %s%s\n", Red, message, Reset)
}
func ShowASCII() {
	asciiArt := `
	_____ __    _____ _____ _____ 
	|   __|  |  |  _  |   __|  |  |
	|   __|  |__|     |__   |     |
	|__|  |_____|__|__|_____|__|__|  v1.0.1
					
						By Mohammed Fathy @Secfathy`
	fmt.Println(asciiArt)
	fmt.Println("")
}
