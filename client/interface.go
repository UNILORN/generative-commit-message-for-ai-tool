package client

// AIClient represents an interface for AI clients that can generate commit messages
type AIClient interface {
	// GenerateCommitMessage generates a commit message based on the provided diff and branch
	GenerateCommitMessage(diff string, branch string) (string, error)
}