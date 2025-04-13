package cmd

import (
	"microgit/utils"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveCmd(t *testing.T) {
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

	err = os.WriteFile("temp.txt", []byte("testing"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed %v", err)
	}

	initCmd.Run(nil, nil)

	// Create test index file
	indexPath := filepath.Join(utils.DEFAULT_PATH, "index")
	testContent := "file1.txt\nfile2.txt\nfile3.txt"
	if err := os.WriteFile(indexPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to update test index: %v", err)
	}

	tests := []struct {
		name          string
		args          []string
		expectedIndex string
	}{
		{
			name:          "Remove single file",
			args:          []string{"file1.txt"},
			expectedIndex: "file2.txt\nfile3.txt",
		},
		{
			name:          "Remove multiple files",
			args:          []string{"file2.txt", "file3.txt"},
			expectedIndex: "file1.txt",
		},
		{
			name:          "Remove all files with dot",
			args:          []string{"."},
			expectedIndex: "",
		},
		{
			name:          "Remove non-existent file",
			args:          []string{"nonexistent.txt"},
			expectedIndex: "file1.txt\nfile2.txt\nfile3.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset index content before each test
			if err := os.WriteFile(indexPath, []byte(testContent), 0644); err != nil {
				t.Fatalf("Failed to reset index: %v", err)
			}

			// Execute remove command
			removeCmd.Run(nil, tt.args)

			// Read resulting index
			got, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatalf("Failed to read index: %v", err)
			}

			if string(got) != tt.expectedIndex {
				t.Errorf("Remove command result = %v, want %v", string(got), tt.expectedIndex)
			}
		})
	}
}
