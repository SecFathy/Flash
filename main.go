package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"source-code-review/config"              // Import config package for ASCII art and colored logging
	"source-code-review/internal/ai/azureai" // Import Azure OpenAI client
	"source-code-review/internal/ai/openai"  // Import OpenAI ChatGPT client
	"source-code-review/internal/markdown"
	"source-code-review/internal/scanner"
)

func main() {
	// Show ASCII art at startup
	config.ShowASCII()

	// Handle graceful shutdown signals (SIGINT, SIGTERM)
	// This channel listens for system signals
	sigs := make(chan os.Signal, 1)
	// Notify the channel if an interrupt or terminate signal is received
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Command line arguments for file or directory and whether to use Azure OpenAI or OpenAI
	filePtr := flag.String("file", "", "File to scan for vulnerabilities")
	dirPtr := flag.String("dir", "", "Directory to scan for vulnerabilities")
	useAzure := flag.Bool("use-azure", false, "Use Azure OpenAI if true, otherwise use OpenAI")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	saveDir := flag.String("save", ".", "Directory to save the results (default is current directory)") // Custom save directory
	flag.Parse()

	// Check if the save directory exists
	if _, err := os.Stat(*saveDir); os.IsNotExist(err) {
		config.Warn(fmt.Sprintf("Save directory does not exist: %s", *saveDir))
		os.Exit(1)
	}

	// Load configuration from file
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		config.Warn(fmt.Sprintf("Failed to load config: %v", err))
		os.Exit(1)
	}

	// Signal handling routine in a separate goroutine
	go func() {
		// Wait for the signal
		sig := <-sigs
		// When a signal is received, attempt graceful shutdown
		config.Info(fmt.Sprintf("Received signal: %s. Attempting graceful shutdown.", sig))
		os.Exit(0)
	}()

	var filesToScan []string
	if *filePtr != "" {
		filesToScan = append(filesToScan, *filePtr)
	} else if *dirPtr != "" {
		files, err := scanner.ScanDirectory(*dirPtr)
		if err != nil {
			config.Warn(fmt.Sprintf("Failed to scan directory: %v", err))
			os.Exit(1)
		}
		filesToScan = files
	} else {
		config.Warn("No file or directory specified")
		os.Exit(1)
	}

	// Iterate over the files and process each one
	for _, file := range filesToScan {
		config.Info(fmt.Sprintf("Scanning %s with Azure OpenAI", file))

		// Read the file content
		content, err := scanner.ScanFile(file)
		if err != nil {
			config.Warn(fmt.Sprintf("Failed to read file: %v", err))
			os.Exit(1)
		}

		var result string

		// Determine whether to use Azure OpenAI or OpenAI ChatGPT
		if *useAzure {
			// Initialize Azure OpenAI client from configuration
			client, err := azureai.NewClient(cfg.AzureOpenAI.Endpoint, cfg.AzureOpenAI.APIKey, cfg.AzureOpenAI.DeploymentName, cfg.AzureOpenAI.APIVersion)
			if err != nil {
				config.Warn(fmt.Sprintf("Failed to create Azure AI client: %v", err))
				os.Exit(1)
			}

			messages := []azureai.Message{
				{Role: "user", Content: content},
			}
			result, err = client.Chat(messages, 800, 0.7)
			if err != nil {
				config.Warn(fmt.Sprintf("Failed to get AI response: %v", err))
				os.Exit(1)
			}

		} else {
			// Initialize OpenAI client from configuration
			client, err := openai.NewClient(cfg.OpenAI.APIKey)
			if err != nil {
				config.Warn(fmt.Sprintf("Failed to create OpenAI client: %v", err))
				os.Exit(1)
			}

			messages := []openai.Message{
				{Role: "user", Content: content},
			}
			result, err = client.Chat(messages, "gpt-3.5-turbo", 800, 0.7)
			if err != nil {
				config.Warn(fmt.Sprintf("Failed to get AI response: %v", err))
				os.Exit(1)
			}
		}

		// Save the result to a markdown file in the specified save directory
		config.Info(fmt.Sprintf("Saving results for %s", file))
		baseName := filepath.Base(file)                       // Get the base file name (without directories)
		resultFile := filepath.Join(*saveDir, baseName+".md") // Save the result in the custom directory
		err = markdown.SaveMarkdown(resultFile, result)
		if err != nil {
			config.Warn(fmt.Sprintf("Failed to save markdown for %s: %v", file, err))
			os.Exit(1)
		}
	}

	config.Info("Scan complete.")
}
