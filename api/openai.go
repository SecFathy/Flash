package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Message represents a message object sent to the OpenAI or Azure API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestBody represents the structure of the request body for OpenAI or Azure API
type RequestBody struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// ErrorResponse represents the Azure OpenAI error response structure
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// ResponseChoice represents each choice in the OpenAI/Azure API response
type ResponseChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}

// APIResponse represents the structure of the API response
type APIResponse struct {
	Choices []ResponseChoice `json:"choices"`
}

// Vulnerability represents a single vulnerability with its details
type Vulnerability struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	ProofOfConcept string `json:"proof_of_concept"`
	Severity       string `json:"severity"`
	VulnerableCode string `json:"vulnerable_code"`
	RecommendedFix string `json:"recommended_fix"`
}

// AnalyzeCode sends the code to OpenAI or Azure OpenAI for vulnerability analysis
func AnalyzeCode(apiKey, apiUrl, code string, azureConfig map[string]string) ([]Vulnerability, error) {
	log.Println("Starting code analysis...")

	// Step 1: Initial analysis
	initialPrompt := "You are a cybersecurity expert. Review this code for vulnerabilities. Provide a high-level summary of potential vulnerabilities you've identified."
	initialAnalysis, err := makeAPICall(apiKey, apiUrl, initialPrompt, code, azureConfig)
	if err != nil {
		return nil, fmt.Errorf("initial analysis failed: %v", err)
	}
	log.Printf("Initial analysis result: %s", initialAnalysis)

	// Step 2: Detailed analysis for each vulnerability
	detailedPrompt := fmt.Sprintf("Based on the initial analysis: %s\n\nNow, for each identified vulnerability, provide the following details:\n- Title\n- Description\n- Proof of Concept\n- Severity (Critical, High, Medium, Low)\n- Vulnerable Code\n- Recommended Fix", initialAnalysis)
	detailedAnalysis, err := makeAPICall(apiKey, apiUrl, detailedPrompt, code, azureConfig)
	if err != nil {
		return nil, fmt.Errorf("detailed analysis failed: %v", err)
	}
	log.Printf("Detailed analysis result: %s", detailedAnalysis)

	// Step 3: Parse and structure the vulnerabilities
	vulnerabilities, err := parseVulnerabilities(detailedAnalysis)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vulnerabilities: %v", err)
	}
	log.Printf("Parsed %d vulnerabilities", len(vulnerabilities))

	return vulnerabilities, nil
}

func makeAPICall(apiKey, apiUrl, prompt, code string, azureConfig map[string]string) (string, error) {
	log.Printf("Making API call with prompt: %s", prompt)

	messages := []Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: code},
	}

	requestBody := RequestBody{
		Model:     "gpt-3.5-turbo",
		Messages:  messages,
		MaxTokens: 2000,
	}

	if azureConfig != nil {
		requestBody.Model = azureConfig["deployment_name"]
		apiUrl = fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
			azureConfig["api_endpoint"],
			azureConfig["deployment_name"],
			azureConfig["api_version"])
	}

	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error preparing request body: %v", err)
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("error creating API request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey) 

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading API response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-OK status: %d", resp.StatusCode)
		log.Printf("Response body: %s", string(body))
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorResp.Error.Code == "" && errorResp.Error.Message == "" {
				return "", fmt.Errorf("API returned status code %d but no error message. Raw response: %s", resp.StatusCode, string(body))
			}
			return "", fmt.Errorf("API error: %s - %s", errorResp.Error.Code, errorResp.Error.Message)
		}
		return "", fmt.Errorf("API returned status code %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("error parsing API response: %v", err)
	}

	if len(apiResponse.Choices) > 0 {
		return apiResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content in API response")
}

func parseVulnerabilities(content string) ([]Vulnerability, error) {
	log.Printf("Raw content received: %s", content)

	var vulnerabilities []Vulnerability
	sections := strings.Split(content, "###") // Split the content by sections marked by "###"

	for _, section := range sections {
		lines := strings.Split(strings.TrimSpace(section), "\n")
		vuln := Vulnerability{}
		for i, line := range lines {
			line = strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(line, "**Title**:"):
				vuln.Title = strings.TrimPrefix(line, "**Title**: ")
			case strings.HasPrefix(line, "**Description**:"):
				vuln.Description = extractContent(lines[i:], "**Description**:")
			case strings.HasPrefix(line, "**Proof of Concept**:"):
				vuln.ProofOfConcept = extractContent(lines[i:], "**Proof of Concept**:")
			case strings.HasPrefix(line, "**Severity**:"):
				vuln.Severity = strings.TrimPrefix(line, "**Severity**: ")
			case strings.HasPrefix(line, "**Vulnerable Code**:"):
				vuln.VulnerableCode = extractContent(lines[i:], "**Vulnerable Code**:")
			case strings.HasPrefix(line, "**Recommended Fix**:"):
				vuln.RecommendedFix = extractContent(lines[i:], "**Recommended Fix**:")
			}
		}
		if vuln.Title != "" {
			vulnerabilities = append(vulnerabilities, vuln)
		}
	}

	if len(vulnerabilities) == 0 {
		return vulnerabilities, fmt.Errorf("no vulnerabilities found in the response")
	}

	log.Printf("Parsed %d vulnerabilities", len(vulnerabilities))
	return vulnerabilities, nil
}

func extractContent(lines []string, prefix string) string {
	var content []string
	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "###") || strings.HasPrefix(line, "**") {
			break
		}
		content = append(content, strings.TrimSpace(line))
	}
	return strings.Join(content, "\n")
}

func MarkdownReport(vulnerabilities []Vulnerability, outputPath string) error {
	var markdown strings.Builder
	markdown.WriteString("# Vulnerability Report\n\n")

	for i, vuln := range vulnerabilities {
		markdown.WriteString(fmt.Sprintf("## %d. %s\n", i+1, vuln.Title))
		markdown.WriteString(fmt.Sprintf("**Description**: %s\n\n", vuln.Description))
		markdown.WriteString(fmt.Sprintf("**Proof of Concept**:\n```\n%s\n```\n\n", vuln.ProofOfConcept))
		markdown.WriteString(fmt.Sprintf("**Severity**: %s\n\n", vuln.Severity))
		markdown.WriteString(fmt.Sprintf("**Vulnerable Code**:\n```\n%s\n```\n\n", vuln.VulnerableCode))
		markdown.WriteString(fmt.Sprintf("**Recommended Fix**:\n```\n%s\n```\n\n", vuln.RecommendedFix))
		markdown.WriteString("---\n\n")
	}

	// Log the generated markdown content for debugging
	log.Printf("Generated markdown content: \n%s", markdown.String())

	// Attempt to write the markdown to the file
	err := os.WriteFile(outputPath, []byte(markdown.String()), 0644)
	if err != nil {
		log.Printf("Failed to write markdown file: %v", err)
		return fmt.Errorf("failed to write markdown file: %v", err)
	}

	log.Printf("Markdown report saved to: %s", outputPath)
	return nil
}

// PrintMarkdown prints the markdown report to the console
func PrintMarkdown(vulnerabilities []Vulnerability) {
	var markdown strings.Builder
	markdown.WriteString("# Vulnerability Report\n\n")

	for i, vuln := range vulnerabilities {
		markdown.WriteString(fmt.Sprintf("## %d. %s\n", i+1, vuln.Title))
		markdown.WriteString(fmt.Sprintf("**Description**: %s\n\n", vuln.Description))
		markdown.WriteString(fmt.Sprintf("**Proof of Concept**:\n```\n%s\n```\n\n", vuln.ProofOfConcept))
		markdown.WriteString(fmt.Sprintf("**Severity**: %s\n\n", vuln.Severity))
		markdown.WriteString(fmt.Sprintf("**Vulnerable Code**:\n```\n%s\n```\n\n", vuln.VulnerableCode))
		markdown.WriteString(fmt.Sprintf("**Recommended Fix**:\n```\n%s\n```\n\n", vuln.RecommendedFix))
		markdown.WriteString("---\n\n")
	}

	fmt.Println(markdown.String())
}
