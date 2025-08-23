package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupGitRepo creates a temporary Git repository for testing
func setupGitRepo(t *testing.T) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize a Git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize git repository: %v", err)
	}

	// Configure Git user for the test repository
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	return tempDir
}

// createAndStageFile creates a file and stages it with Git
func createAndStageFile(t *testing.T, repoDir, filename, content string) {
	// Create a file
	filePath := filepath.Join(repoDir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Stage the file
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}
}

func TestGetStagedDiff(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git is not installed, skipping test")
	}

	// Setup a temporary Git repository
	repoDir := setupGitRepo(t)
	defer os.RemoveAll(repoDir)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Change to the test repository directory
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create and stage a file
	createAndStageFile(t, repoDir, "test.txt", "Hello, World!")

	// Get the staged diff
	diff, err := GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	// Verify the diff contains the expected content
	if !strings.Contains(diff, "test.txt") || !strings.Contains(diff, "Hello, World!") {
		t.Errorf("Diff does not contain expected content: %s", diff)
	}
}

func TestGetStagedFiles(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git is not installed, skipping test")
	}

	// Setup a temporary Git repository
	repoDir := setupGitRepo(t)
	defer os.RemoveAll(repoDir)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Change to the test repository directory
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create and stage multiple files
	createAndStageFile(t, repoDir, "file1.txt", "Content 1")
	createAndStageFile(t, repoDir, "file2.txt", "Content 2")

	// Get the staged files
	files, err := GetStagedFiles()
	if err != nil {
		t.Fatalf("GetStagedFiles failed: %v", err)
	}

	// Verify the files list contains the expected files
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	foundFile1 := false
	foundFile2 := false
	for _, file := range files {
		if file == "file1.txt" {
			foundFile1 = true
		}
		if file == "file2.txt" {
			foundFile2 = true
		}
	}

	if !foundFile1 || !foundFile2 {
		t.Errorf("Files list does not contain expected files: %v", files)
	}
}

func TestGetStagedFilesWithStatus(t *testing.T) {
	// Skip if git is not installed
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git is not installed, skipping test")
	}

	// Setup a temporary Git repository
	repoDir := setupGitRepo(t)
	defer os.RemoveAll(repoDir)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Change to the test repository directory
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create and stage a new file
	createAndStageFile(t, repoDir, "new-file.txt", "New file content")

	// Get the staged files with status
	filesWithStatus, err := GetStagedFilesWithStatus()
	if err != nil {
		t.Fatalf("GetStagedFilesWithStatus failed: %v", err)
	}

	// Verify the output contains the expected file and status
	if !strings.Contains(filesWithStatus, "new-file.txt") {
		t.Errorf("Files with status does not contain expected file: %s", filesWithStatus)
	}

	// The status should be 'A' for added files
	if !strings.Contains(filesWithStatus, "A") {
		t.Errorf("Files with status does not contain expected status: %s", filesWithStatus)
	}
}
