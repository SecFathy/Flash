package main

import (
	"flag"
	"fmt"
	"myapp/api"
	"myapp/ascii"
	"myapp/config"
	"myapp/logger"
	"os"
)

func main() {
	// Ensure ASCII art is shown at the beginning
	ascii.PrintASCII()

	// Define command-line arguments
	codeFilePath := flag.String("C", "", "Path to the file containing the code to review")
	outputFilePath := flag.String("O", "", "Optional: Path to save the markdown results (e.g., output.md)")
	help := flag.Bool("H", false, "Show help message")

	// Parse the command-line arguments
	flag.Parse()

	// If no arguments are provided and help is not requested, show help and check API status
	if *codeFilePath == "" && !*help {
		logger.Info("No code file provided. Please use -H for help.")
		checkAPIConfiguration()
		return
	}

	// If help flag is provided, show help message and exit
	if *help {
		showHelp()
		return
	}

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		logger.Error(fmt.Sprintf("Error loading .env file: %v", err))
		return
	}

	// Read the code from the file
	code, err := readCodeFromFile(*codeFilePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error reading code from file: %v", err))
		return
	}

	// Get API configuration
	apiKey, apiUrl, provider, azureConfig := config.GetAPIConfig()
	if apiKey == "" || apiUrl == "" {
		logger.Warn("API configuration missing. Either Azure or OpenAI credentials must be set.")
		logger.Info("No results will be produced without valid API credentials.")
		return
	}

	// Log provider information
	logger.Info(fmt.Sprintf("Using %s API", provider))

	// Send code for analysis
	vulnerabilities, err := api.AnalyzeCode(apiKey, apiUrl, code, azureConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Error analyzing code: %v", err))
		return
	}

	// Check if vulnerabilities are found
	if len(vulnerabilities) == 0 {
		logger.Warn("No vulnerabilities found in the code.")
	}

	// Print the analysis result as Markdown in the console
	api.PrintMarkdown(vulnerabilities)

	// If output file path is provided, save the markdown report to file
	if *outputFilePath != "" {
		err := api.MarkdownReport(vulnerabilities, *outputFilePath)
		if err != nil {
			logger.Error(fmt.Sprintf("Error saving markdown report to %s: %v", *outputFilePath, err))
			return
		}
		// Verify if the file was successfully created and contains content
		if _, err := os.Stat(*outputFilePath); os.IsNotExist(err) {
			logger.Error(fmt.Sprintf("Failed to create markdown file: %s", *outputFilePath))
			return
		}
		fileContent, err := os.ReadFile(*outputFilePath)
		if err != nil {
			logger.Error(fmt.Sprintf("Error reading saved markdown file: %v", err))
			return
		}
		if len(fileContent) == 0 {
			logger.Error("Saved markdown file is empty")
			return
		}
		logger.Info(fmt.Sprintf("Results successfully saved to %s", *outputFilePath))
	}
}

// showHelp displays usage instructions
func showHelp() {
	fmt.Println("Usage:")
	fmt.Println("  -C <path>   : Path to the code file to be analyzed")
	fmt.Println("  -O <path>   : Optional: Path to save the markdown results (e.g., output.md)")
	fmt.Println("  -H          : Show help message")
}

// readCodeFromFile reads the code from the provided file path
func readCodeFromFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to read the file: %w", err)
	}
	return string(content), nil
}

// checkAPIConfiguration checks whether the OpenAI or Azure OpenAI API configuration is present
func checkAPIConfiguration() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	azureKey := os.Getenv("AZURE_API_KEY")

	if apiKey != "" {
		logger.Info("OpenAI API is configured and ready.")
	} else if azureKey != "" {
		logger.Info("Azure OpenAI API is configured and ready.")
	} else {
		logger.Warn("No OpenAI or Azure API configuration found. Set the appropriate environment variables.")
	}
}
