package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
	appconfig "github.com/UNILORN/generative-commit-message-for-bedrock.git/config"
)

// Client represents a Claude API client
type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Claude API client
func NewClient(model string) (*Client, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is not set")
	}

	return &Client{
		apiKey:  apiKey,
		model:   model,
		baseURL: "https://api.anthropic.com",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// ClaudeMessage represents a message in the Claude API format
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeRequest represents a request to the Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeResponseContent represents content in the Claude API response
type ClaudeResponseContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeUsage represents usage information in the Claude API response
type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ClaudeResponse represents a response from the Claude API
type ClaudeResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Role       string                  `json:"role"`
	Model      string                  `json:"model"`
	Content    []ClaudeResponseContent `json:"content"`
	StopReason string                  `json:"stop_reason"`
	Usage      ClaudeUsage             `json:"usage"`
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Get config and build prompt
	cfg := appconfig.Get()
	prompt := cfg.BuildPrompt("japanese", branch, diff)

	// Create the request
	request := ClaudeRequest{
		Model:     c.model,
		MaxTokens: 10000,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Marshal the request to JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/v1/messages", bytes.NewBuffer(requestBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers according to Claude API documentation
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the response
	var response ClaudeResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the commit message
	if len(response.Content) > 0 && len(response.Content[0].Text) > 0 {
		return response.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}