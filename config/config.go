package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config structure to hold all credentials
type Config struct {
	AzureOpenAI struct {
		Endpoint       string `json:"endpoint"`
		APIKey         string `json:"api_key"`
		DeploymentName string `json:"deployment_name"`
		APIVersion     string `json:"api_version"`
	} `json:"azure_openai"`
	OpenAI struct {
		APIKey string `json:"api_key"`
	} `json:"openai"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return &config, nil
}
