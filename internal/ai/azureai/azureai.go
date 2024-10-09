package azureai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Endpoint       string
	APIKey         string
	DeploymentName string
	APIVersion     string
	HTTPClient     *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewClient initializes a new Azure AI client.
func NewClient(endpoint, apiKey, deploymentName, apiVersion string) (*Client, error) {
	if endpoint == "" || apiKey == "" || deploymentName == "" || apiVersion == "" {
		return nil, fmt.Errorf("all parameters must be provided")
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &Client{
		Endpoint:       endpoint,
		APIKey:         apiKey,
		DeploymentName: deploymentName,
		APIVersion:     apiVersion,
		HTTPClient:     httpClient,
	}, nil
}

// Chat sends a request to Azure OpenAI.
func (c *Client) Chat(messages []Message, maxTokens int, temperature float64) (string, error) {
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", c.Endpoint, c.DeploymentName, c.APIVersion)

	requestBody := ChatRequest{
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(chatResponse.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResponse.Choices[0].Message.Content, nil
}
