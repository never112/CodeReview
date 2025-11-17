package review

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func PerformCodeReview(ctx context.Context, repoDir, reviewCommand string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"directory": repoDir,
		"command":   reviewCommand,
	}).Info("Starting code review")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Change to the repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoDir); err != nil {
		return "", fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Execute the review command
	cmd := exec.CommandContext(ctx, "sh", "-c", reviewCommand)

	// Capture stdout and stderr
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Set environment variables
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", repoDir),
		"CODE_REVIEW=true",
	)

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start)

	logrus.WithFields(logrus.Fields{
		"duration":       duration,
		"Review Command": reviewCommand,
		"exit_code":      cmd.ProcessState.ExitCode(),
	}).Info("Code review command completed")

	// Prepare the result
	var result strings.Builder
	//result.WriteString(fmt.Sprintf("**Review Command:** `%s`\n\n", reviewCommand))
	//result.WriteString(fmt.Sprintf("**Duration:** %v\n\n", duration))
	//result.WriteString(fmt.Sprintf("**Exit Code:** %d\n\n", cmd.ProcessState.ExitCode()))

	// Add output sections
	if stdout.Len() > 0 {
		result.WriteString("### 📋 Review Output\n\n```\n")
		result.WriteString(stdout.String())
		result.WriteString("\n```\n\n")
	}

	if stderr.Len() > 0 {
		result.WriteString("### ⚠️ Warnings/Errors\n\n```\n")
		result.WriteString(stderr.String())
		result.WriteString("\n```\n\n")
	}

	// Add summary
	if err != nil {
		result.WriteString("### ❌ Status\n\nReview completed with errors.\n")
	} else {
		result.WriteString("### ✅ Status\n\nReview completed successfully.\n")
	}

	// Add file information if available
	if files, err := getChangedFiles(repoDir); err == nil && len(files) > 0 {
		result.WriteString(fmt.Sprintf("### 📁 Files Analyzed\n\n%d files were included in this review.\n\n", len(files)))

		// Add a summary of file types
		fileTypes := make(map[string]int)
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file))
			if ext == "" {
				ext = "no_extension"
			}
			fileTypes[ext]++
		}

		if len(fileTypes) > 0 {
			result.WriteString("**File Types:**\n")
			for ext, count := range fileTypes {
				result.WriteString(fmt.Sprintf("- `%s`: %d files\n", ext, count))
			}
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

func getChangedFiles(repoDir string) ([]string, error) {
	// Get list of all non-binary files in the repository
	var files []string

	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// Skip .git directory
		if strings.Contains(path, ".git") {
			return nil
		}

		// Skip binary files (basic check)
		ext := strings.ToLower(filepath.Ext(path))
		skipExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true,
			".ico": true, ".svg": true, ".pdf": true, ".zip": true, ".tar": true,
			".gz": true, ".exe": true, ".dll": true, ".so": true, ".dylib": true,
		}

		if !skipExts[ext] {
			// Get relative path
			relPath, err := filepath.Rel(repoDir, path)
			if err == nil {
				files = append(files, relPath)
			}
		}

		return nil
	})

	return files, err
}
