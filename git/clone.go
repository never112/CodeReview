package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func CloneRepository(ctx context.Context, cloneURL, branch, workDir string) (string, error) {
	// Create a unique directory name for this repository
	repoName := strings.TrimSuffix(filepath.Base(cloneURL), ".git")
	tempDir := filepath.Join(workDir, repoName+"-"+branch)

	// Remove existing directory if it exists
	if _, err := os.Stat(tempDir); err == nil {
		if err := os.RemoveAll(tempDir); err != nil {
			return "", fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create the working directory if it doesn't exist
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create work directory: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"url":    cloneURL,
		"branch": branch,
		"dir":    tempDir,
	}).Info("Cloning repository")

	// Clone the repository
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--branch", branch, cloneURL, tempDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	logrus.WithField("directory", tempDir).Info("Successfully cloned repository")

	return tempDir, nil
}

func CleanupRepository(repoDir string) error {
	logrus.WithField("directory", repoDir).Info("Cleaning up repository")

	if err := os.RemoveAll(repoDir); err != nil {
		return fmt.Errorf("failed to cleanup repository directory: %w", err)
	}

	return nil
}