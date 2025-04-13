package cmd

import (
	"encoding/json"
	"fmt"
	"microgit/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func getHead() string {
	headPath := filepath.Join(utils.DEFAULT_PATH, "HEAD")

	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func setHead(hash string) error {
	headPath := filepath.Join(utils.DEFAULT_PATH, "HEAD")
	latestPath := filepath.Join(utils.DEFAULT_PATH, "LATEST")

	errChan := make(chan error, 2) // Buffer for 2 potential errors

	// Write HEAD file asynchronously
	go func() {
		err := os.WriteFile(headPath, []byte(hash), 0644)
		errChan <- err
	}()

	// Write LATEST file asynchronously
	go func() {
		err := os.WriteFile(latestPath, []byte(hash), 0644)
		errChan <- err
	}()

	// Wait for both operations to complete and check for errors
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return fmt.Errorf("failed to write reference file: %w", err)
		}
	}

	return nil
}

func readIndex() (map[string]string, error) {
	index := make(map[string]string)

	indexPath := filepath.Join(utils.DEFAULT_PATH, "index")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return index, err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			index[parts[0]] = parts[1]
		}
	}
	return index, nil
}

func writeSavePointObject(savePoint utils.SavePoint) (string, error) {
	jsonData, err := json.MarshalIndent(savePoint, "", "  ")
	if err != nil {
		return "", err
	}

	// Hash of the entire SavePoint JSON
	hash := utils.HashContent(jsonData)

	err = utils.WriteObject(hash, jsonData)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the current state of staged files",
	Long: `Save the current state of all staged files as a new commit.
This command requires a commit message that describes the changes being saved.
The staged files will be committed and the staging area will be cleared after the save.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: microgit save \"message\"")
			return
		}

		message := args[0]

		index, err := readIndex()
		if err != nil {
			fmt.Println("could not read index: %w", err)
			return
		}

		if len(index) == 0 {
			fmt.Println("No files have been added")
			return
		}

		parent := getHead()
		timestamp := time.Now().Format(time.RFC3339)

		savePoint := utils.SavePoint{
			Message:   message,
			Timestamp: timestamp,
			Parent:    parent,
			Files:     index,
		}

		hash, err := writeSavePointObject(savePoint)
		if err != nil {
			fmt.Println("failed to write commit: %w", err)
			return
		}

		err = setHead(hash)
		if err != nil {
			fmt.Println("failed to update HEAD: %w", err)
			return
		}

		fmt.Printf("Saved: %s\n", hash)

		// Clear the staging area

		indexPath := filepath.Join(utils.DEFAULT_PATH, "index")
		os.WriteFile(indexPath, []byte(""), 0644)
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)
}
