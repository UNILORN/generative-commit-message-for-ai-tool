package copilotcli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-bedrock.git/client"
)

// Client represents a Copilot CLI client
type Client struct {
	model string
}

// Ensure Client implements the AIClient interface
var _ client.AIClient = (*Client)(nil)

// NewClient creates a new Copilot CLI client
func NewClient(model string) (*Client, error) {
	// Check if copilot command is available
	if _, err := exec.LookPath("copilot"); err != nil {
		return nil, fmt.Errorf("copilot command not found in PATH: %w", err)
	}

	// Set default model if not provided
	if model == "" {
		model = "claude-sonnet-4.5"
	}

	return &Client{
		model: model,
	}, nil
}

// GenerateCommitMessage generates a commit message based on the provided diff
func (c *Client) GenerateCommitMessage(diff string, branch string) (string, error) {
	// Create the prompt that requests only the commit message without any explanation
	prompt := fmt.Sprintf(`You are a commit message generator. Output ONLY the commit message, nothing else.

CRITICAL RULES:
- Output ONLY the commit message itself
- NO explanations, NO analysis, NO comments
- NO markdown code blocks
- NO "Here is...", NO "The commit message is..."
- Just the raw commit message text

Commit Message Guidelines:
- Start with a short summary line (50-72 characters)
- Use imperative mood
- After a blank line, include detailed explanation as bullet points (using - for each point)
- Focus on "why" and "what" rather than "how"
- Use Semantic Release prefix format
- If branch name contains number, include it like: 'feat: message #1234'

Current branch: %s

Semantic Release Prefixes:
- "feat: :sparkles:" : New feature
- "fix: :bug:" : Bug fix
- "refactor: :hammer:" : Code refactoring
- "test: :white_check_mark:" : Test changes
- "docs: :memo:" : Documentation only
- "config: :wrench:" : Configuration files
- "lint: :rotating_light:" : Linter warnings fix
- "ci: :construction_worker:" : CI changes
- "remove: :wastebasket:" : Deletion
- "improve: :zap:" : Performance improvement
- "try: :bulb:" : Experimental changes
- "wip: :construction:" : WIP
- "update: :up:" : Library updates
- "release: :rocket:" : Release
- "merge: :twisted_rightwards_arrows:" : Merge

Git Diff:
%s

OUTPUT FORMAT: Just output the commit message in Japanese. No JSON, no code blocks, no explanations.`, branch, diff)

	// Execute copilot command with -p flag for prompt and --model for model specification
	cmd := exec.Command("copilot", "-p", prompt, "--model", c.model)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute copilot command: %w\nstderr: %s", err, stderr.String())
	}

	response := strings.TrimSpace(stdout.String())
	if response == "" {
		return "", fmt.Errorf("empty response from copilot command")
	}

	// Remove usage statistics from the end first
	if idx := strings.Index(response, "\n\nTotal usage"); idx > 0 {
		response = response[:idx]
	}
	if idx := strings.Index(response, "\nTotal usage"); idx > 0 {
		response = response[:idx]
	}

	// Remove leading bullet point (●) that Copilot CLI often adds
	response = strings.TrimPrefix(response, "●")
	response = strings.TrimPrefix(response, "● ")
	response = strings.TrimSpace(response)

	// List of Semantic Release prefixes to detect commit message start
	prefixes := []string{
		"feat:", "fix:", "refactor:", "test:", "docs:", "config:",
		"lint:", "ci:", "remove:", "improve:", "try:", "wip:",
		"update:", "release:", "merge:",
	}

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
