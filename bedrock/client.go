package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
	appconfig "github.com/UNILORN/generative-commit-message-for-bedrock.git/config"
)

// Client represents an AWS Bedrock client
type Client struct {
	bedrockClient *bedrockruntime.Client
	modelID       string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new AWS Bedrock client
func NewClient(region, modelID string) (*Client, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create Bedrock client
	bedrockClient := bedrockruntime.NewFromConfig(cfg)

	return &Client{
		bedrockClient: bedrockClient,
		modelID:       modelID,
	}, nil
}

// AnthropicMessage represents a message in the Anthropic API format
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicRequest represents a request to the Anthropic API
type AnthropicRequest struct {
	AnthropicVersion string             `json:"anthropic_version"`
	Messages         []AnthropicMessage `json:"messages"`
	MaxTokens        int                `json:"max_tokens"`
}

// AnthropicContent represents content in the Anthropic API response
type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicResponseMessage represents a message in the Anthropic API response
type AnthropicResponseMessage struct {
	Role    string             `json:"role"`
	Content []AnthropicContent `json:"content"`
}

// AnthropicResponse represents a response from the Anthropic API
type AnthropicResponse struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	Role         string             `json:"role"`
	Model        string             `json:"model"`
	Content      []AnthropicContent `json:"content"`
	StopReason   string             `json:"stop_reason"`
	StopSequence *string            `json:"stop_sequence"`
	Usage        AnthropicUsage     `json:"usage"`
}

// AnthropicUsage represents usage information in the Anthropic API response
type AnthropicUsage struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	OutputTokens             int `json:"output_tokens"`
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	return c.generateWithAnthropic(diff, branch)
}

// generateWithAnthropic generates a commit message using an Anthropic model
func (c *Client) generateWithAnthropic(diff string, branch string) (string, error) {
	// Get config and build prompt
	cfg := appconfig.Get()
	prompt := cfg.BuildPrompt("japanese", branch, diff)

	// Create the request
	request := AnthropicRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 10000,
	}

	// Marshal the request to JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the Bedrock invoke request
	invokeInput := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		ContentType: aws.String("application/json"),
		Body:        requestBytes,
	}

	// Invoke the model
	resp, err := c.bedrockClient.InvokeModel(context.TODO(), invokeInput)
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %w", err)
	}

	// Parse the response
	var response AnthropicResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract the commit message
	if len(response.Content) > 0 && len(response.Content[0].Text) > 0 {
		return response.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}
