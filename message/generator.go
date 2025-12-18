package message

import (
	"fmt"
	"strings"

	"github.com/UNILORN/generative-commit-message-for-bedrock/client"
	"github.com/UNILORN/generative-commit-message-for-bedrock/git"
)

// Generate generates a commit message based on the provided diff
func Generate(aiClient client.AIClient, diff string, branch string) (string, error) {
	// If diff is empty, try to get more context from staged files
	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no diff provided")
	}

	// Get the list of staged files with status for additional context
	filesWithStatus, err := git.GetStagedFilesWithStatus()
	if err != nil {
		return "", fmt.Errorf("failed to get staged files: %w", err)
	}

	// If we have a lot of files, we might want to include a summary
	// in the prompt to help the AI generate a better commit message
	if len(filesWithStatus) > 0 {
		diff = fmt.Sprintf("Files changed:\n%s\n\nDiff:\n%s", filesWithStatus, diff)
	}

	// Generate the commit message using the AI client
	commitMsg, err := aiClient.GenerateCommitMessage(diff, branch)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	return commitMsg, nil
}

// ApplyCommitMessage applies the generated commit message using git commit
// This is left as a future enhancement
func ApplyCommitMessage(message string) error {
	// This could be implemented to automatically commit with the generated message
	// For now, we just print the message and let the user commit manually
	return nil
}
