package geminicli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
	appconfig "github.com/UNILORN/generative-commit-message-for-bedrock.git/config"
)

// Client represents a Gemini CLI client
type Client struct {
	model string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Gemini CLI client
func NewClient(model string) (*Client, error) {
	// Check if gemini command is available
	if _, err := exec.LookPath("gemini"); err != nil {
		return nil, fmt.Errorf("gemini command not found in PATH: %w", err)
	}

	// Set default model if not provided
	if model == "" {
		model = "gemini-2.5-pro"
	}

	return &Client{
		model: model,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Get config and build prompt
	cfg := appconfig.Get()
	prompt := cfg.BuildPrompt("japanese", branch, diff)

	// Execute gemini command
	cmd := exec.Command("gemini", "--model", c.model, "--prompt", prompt)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute gemini command: %w\nstderr: %s", err, stderr.String())
	}

	response := strings.TrimSpace(stdout.String())
	if response == "" {
		return "", fmt.Errorf("empty response from gemini command")
	}

	return response, nil
}