package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// hashContent returns the SHA-256 hash of the file content
func hashContent(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

// writeObject saves the file content to objects/<hash>
func writeObject(hash string, content []byte) error {
	objectPath := filepath.Join(DEFAULT_PATH, "objects", hash)
	return os.WriteFile(objectPath, content, 0644)
}

// updateIndex writes or updates the index file with path -> hash
func updateIndex(filePath, hash string) error {
	indexPath := filepath.Join(DEFAULT_PATH, "index")
	existing := ""

	if data, err := os.ReadFile(indexPath); err == nil {
		existing = string(data)
	}

	lines := strings.Split(existing, "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, filePath+" ") {
			lines[i] = filePath + " " + hash
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, filePath+" "+hash)
	}

	newIndex := strings.Join(lines, "\n")
	return os.WriteFile(indexPath, []byte(newIndex), 0644)
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add files to the staging area",
	Long: `Add files to the staging area. This command takes one or more file paths
as arguments and stages them for the next commit.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No files specified")
			return
		}

		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		for _, file := range args {
			fullFilePath := filepath.Join(currentDir, file)

			// Read file content
			content, err := os.ReadFile(fullFilePath)
			if err != nil {
				fmt.Printf("Error reading file '%s': %v\n", file, err)
				continue
			}

			// Calculate hash
			hash := hashContent(content)

			// Write object
			if err := writeObject(hash, content); err != nil {
				fmt.Printf("Error writing object for '%s': %v\n", file, err)
				continue
			}

			// Update index
			if err := updateIndex(file, hash); err != nil {
				fmt.Printf("Error updating index for '%s': %v\n", file, err)
				continue
			}

			fmt.Printf("Added %s (hash: %s)\n", file, hash)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
