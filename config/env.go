package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv checks if the .env file exists and creates it if not, then loads the environment variables
func LoadEnv() error {
	envFile := ".env"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// If the .env file does not exist, create it with default content
		fmt.Println(".env file not found. Creating a new one with default settings...")
		createEnvFile(envFile)
	} else {
		fmt.Println(".env file found. Loading configuration...")
	}

	// Load the .env file
	return godotenv.Load(envFile)
}

// createEnvFile creates the default .env file with required placeholders if it doesn't exist
func createEnvFile(filePath string) {
	content := `# Choose which API to use (OpenAI or Azure OpenAI)
# Set USE_OPENAI=true if using OpenAI API
# Set USE_AZURE=true if using Azure OpenAI API
USE_OPENAI=true
USE_AZURE=false

# OpenAI API configuration
OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-3.5-turbo

# Azure OpenAI API configuration
AZURE_API_KEY=your_azure_api_key_here
AZURE_API_VERSION=2023-03-15-preview
AZURE_API_ENDPOINT=your_azure_endpoint_here
AZURE_DEPLOYMENT_NAME=your_azure_deployment_name_here
`
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
		return
	}

	fmt.Println(".env file created with default settings. Please update it with your API keys.")
}

// GetAPIConfig determines whether to use OpenAI or Azure OpenAI based on the environment variables
func GetAPIConfig() (string, string, string, map[string]string) {
	useOpenAI := os.Getenv("USE_OPENAI") == "true"
	useAzure := os.Getenv("USE_AZURE") == "true"

	if useOpenAI {
		return os.Getenv("OPENAI_API_KEY"), "https://api.openai.com/v1/chat/completions", "OpenAI", nil
	} else if useAzure {
		// Fetch Azure-specific settings
		azureConfig := map[string]string{
			"api_key":         os.Getenv("AZURE_API_KEY"),
			"api_version":     os.Getenv("AZURE_API_VERSION"),
			"api_endpoint":    os.Getenv("AZURE_API_ENDPOINT"),
			"deployment_name": os.Getenv("AZURE_DEPLOYMENT_NAME"),
		}

		// Ensure all required variables are set
		for key, value := range azureConfig {
			if value == "" {
				fmt.Printf("Missing Azure OpenAI Configuration: %s is not set in the .env file\n", key)
				return "", "", "", nil
			}
		}

		return azureConfig["api_key"], azureConfig["api_endpoint"], "Azure OpenAI", azureConfig
	}

	return "", "", "", nil
}
