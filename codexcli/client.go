package codexcli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-ai-tool/client"
	appconfig "github.com/UNILORN/generative-commit-message-for-ai-tool/config"
)

// Client represents an OpenAI Codex CLI client
type Client struct {
	model string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Codex CLI client
func NewClient(model string) (*Client, error) {
	// Check if codex command is available
	if _, err := exec.LookPath("codex"); err != nil {
		return nil, fmt.Errorf("codex command not found in PATH: %w", err)
	}

	// Set default model if not provided
	if model == "" {
		model = "o4-mini"
	}

	return &Client{
		model: model,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Get config and build prompt (English version)
	cfg := appconfig.Get()
	prompt := cfg.BuildPromptEnglish(branch, diff)

	// Execute codex exec with stdin piping to handle large prompts
	// codex exec - reads the prompt from stdin in non-interactive mode
	cmd := exec.Command("codex", "exec", "--model", c.model, "-")
	cmd.Stdin = strings.NewReader(prompt)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute codex command: %w\nstderr: %s", err, stderr.String())
	}

	response := strings.TrimSpace(stdout.String())
	if response == "" {
		return "", fmt.Errorf("empty response from codex command")
	}

	// Remove usage statistics from the end first
	if idx := strings.Index(response, "\n\nTotal usage"); idx > 0 {
		response = response[:idx]
	}
	if idx := strings.Index(response, "\nTotal usage"); idx > 0 {
		response = response[:idx]
	}

	// Get list of Semantic Release prefixes from config
	prefixes := cfg.GetPrefixList()

	// Try to find the first line that starts with a Semantic Release prefix
	lines := strings.Split(response, "\n")
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

	// Fallback: return the whole response (without usage stats)
	return response, nil
}
