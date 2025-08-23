package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the diff of all staged changes
func GetStagedDiff() (string, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}

	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Get the diff of staged changes
	cmd = exec.Command("git", "diff", "--staged")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

// GetStagedFiles returns a list of staged files
func GetStagedFiles() ([]string, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return nil, err
	}

	// Get the list of staged files
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(out.String()), "\n")
	// Handle empty output
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}
	return files, nil
}

// GetStagedFilesWithStatus returns a list of staged files with their status
func GetStagedFilesWithStatus() (string, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}

	// Get the list of staged files with status
	cmd := exec.Command("git", "diff", "--staged", "--name-status")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
