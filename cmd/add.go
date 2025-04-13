package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"microgit/utils"
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
	objectPath := filepath.Join(utils.DEFAULT_PATH, "objects", hash)
	return os.WriteFile(objectPath, content, 0644)
}

// updateIndex writes or updates the index file with path -> hash
func updateIndex(filePath, hash string) error {
	indexPath := filepath.Join(utils.DEFAULT_PATH, "index")
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

	newIndex := strings.Trim(strings.Join(lines, "\n"), "\n")
	return os.WriteFile(indexPath, []byte(newIndex), 0644)
}

func stageFile(path, fileName string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", fileName, err)
		return nil
	}

	// Calculate hash
	hash := hashContent(content)

	// Write object
	if err := writeObject(hash, content); err != nil {
		fmt.Printf("Error writing object for '%s': %v\n", fileName, err)
		return nil
	}

	// Update index
	if err := updateIndex(path, hash); err != nil {
		fmt.Printf("Error updating index for '%s': %v\n", fileName, err)
		return nil
	}

	fmt.Printf("Added %s (hash: %s)\n", path, hash)

	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add files to the staging area",
	Long: `Add files to the staging area for the next commit.

Usage:
  microgit add <file1> [file2 ...]  - Stage specific files
  microgit add .                    - Stage all files in current directory

The add command will:
1. Calculate a SHA-256 hash of the file content
2. Store the file content in the objects directory
3. Update the index with the file path and corresponding hash

Files in the .microgit/ and .git/ directories are automatically ignored.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No files specified")
			return
		}

		if args[0] == "." {
			err := filepath.WalkDir(".", func(path string, file os.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if file.IsDir() || strings.HasPrefix(path, utils.DEFAULT_PATH) || strings.HasPrefix(path, ".git/") {
					return nil
				}

				return stageFile(path, file.Name())
			})
			if err != nil {
				fmt.Printf("Error reading directory: %v\n", err)
				return
			}
			return
		}

		for _, file := range args {
			err := stageFile(file, file)
			if err != nil {
				fmt.Printf("Error staging file: %v", err)
				continue
			}
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
