package copilotsdk

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	copilot "github.com/github/copilot-sdk/go"

	"github.com/UNILORN/generative-commit-message-for-ai-tool/client"
	appconfig "github.com/UNILORN/generative-commit-message-for-ai-tool/config"
)

// Client represents a Copilot SDK client that uses the GitHub Copilot SDK
// for programmatic access to Copilot CLI, providing session management
// and streaming capabilities.
type Client struct {
	model string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Copilot SDK client.
// Unlike copilotcli which executes the CLI directly, copilotsdk uses
// the GitHub Copilot SDK for programmatic access with session management.
func NewClient(model string) (*Client, error) {
	// Check if copilot command is available (required by the SDK)
	if _, err := exec.LookPath("copilot"); err != nil {
		return nil, fmt.Errorf("copilot command not found in PATH (required by Copilot SDK): %w", err)
	}

	// Set default model if not provided
	if model == "" {
		model = "gpt-4o"
	}

	return &Client{
		model: model,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Get config and build prompt
	cfg := appconfig.Get()
	prompt := cfg.BuildPromptEnglish(branch, diff)

	// Create Copilot client
	copilotClient := copilot.NewClient(nil)

	// Start the client (this starts the Copilot CLI server)
	if err := copilotClient.Start(); err != nil {
		return "", fmt.Errorf("failed to start copilot client: %w", err)
	}
	defer copilotClient.Stop()

	// Create a session with the specified model
	session, err := copilotClient.CreateSession(&copilot.SessionConfig{
		Model: c.model,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Destroy()

	// Send the prompt and wait for the response
	// Timeout is set to 120 seconds to allow for complex diffs and model processing time
	response, err := session.SendAndWait(copilot.MessageOptions{
		Prompt: prompt,
	}, 120*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	if response == nil || response.Data.Content == nil {
		return "", fmt.Errorf("empty response from copilot SDK")
	}

	responseText := strings.TrimSpace(*response.Data.Content)
	if responseText == "" {
		return "", fmt.Errorf("empty response from copilot SDK")
	}

	// Remove leading bullet point (●) that Copilot sometimes adds
	responseText = strings.TrimPrefix(responseText, "●")
	responseText = strings.TrimPrefix(responseText, "● ")
	responseText = strings.TrimSpace(responseText)

	// Get list of Semantic Release prefixes from config
	prefixes := cfg.GetPrefixList()

	// Try to find the first line that starts with a Semantic Release prefix
	lines := strings.Split(responseText, "\n")
	startIdx := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				startIdx = i
				break
			}
		}
		if startIdx >= 0 {
			break
		}
	}

	// If we found a line starting with a prefix, extract from there
	if startIdx >= 0 {
		commitLines := lines[startIdx:]
		commitMessage := strings.Join(commitLines, "\n")
		return strings.TrimSpace(commitMessage), nil
	}

	// Fallback: return the whole response
	return responseText, nil
}
