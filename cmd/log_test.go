package cmd

import (
	"microgit/utils"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadCommit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "microgit-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to the temporary directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	objectsDir := filepath.Join(utils.DEFAULT_PATH, "objects")

	err = os.WriteFile("temp.txt", []byte("testing"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed %v", err)
	}

	initCmd.Run(nil, nil)
	addCmd.Run(nil, []string{"temp.txt"})
	saveCmd.Run(nil, []string{"temp"})

	headPath := filepath.Join(utils.DEFAULT_PATH, "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		t.Errorf("Error reading HEAD file %v", err)
	}
	commitHash := strings.TrimSpace(string(data))

	t.Run("successful commit read", func(t *testing.T) {
		// Test reading the commit
		result, err := readCommit(commitHash)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Compare results
		if result.Message != "temp" {
			t.Errorf("Expected message %q, got %q", "temp", result.Message)
		}
		if result.Parent != "" {
			t.Errorf("Expected parent %q, got %q", "", result.Parent)
		}
		if len(result.Files) != 1 {
			t.Errorf("Expected %d files, got %d", 1, len(result.Files))
		}
	})

	t.Run("nonexistent commit", func(t *testing.T) {
		_, err := readCommit("nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent commit, got nil")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		// Write invalid JSON to file
		commitHash := "invalid123"
		invalidJSON := []byte("{invalid json")
		commitPath := filepath.Join(objectsDir, commitHash)
		if err := os.WriteFile(commitPath, invalidJSON, 0644); err != nil {
			t.Fatalf("Failed to write invalid commit: %v", err)
		}

		_, err := readCommit(commitHash)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}
