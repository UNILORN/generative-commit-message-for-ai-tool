package git

import (
	"bytes"
	"os/exec"
	"strings"
)

// GetCurrentBranch returns the name of the current branch
func GetCurrentBranch() (string, error) {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return "", err
	}
	// Get the current branch name
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
