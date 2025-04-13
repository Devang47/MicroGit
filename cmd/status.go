package cmd

import (
	"fmt"
	"microgit/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type fileResult struct {
	data interface{}
	err  error
}

func getWorkingFiles() (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip internal directory
		if strings.HasPrefix(path, utils.DEFAULT_PATH) || strings.HasPrefix(path, ".git/") || info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		files[path] = utils.HashContent(content)
		return nil
	})

	return files, err
}

func getCommittedFiles() map[string]string {
	head := getHead()
	if head == "" {
		return map[string]string{}
	}

	commit, err := readCommit(head)
	if err != nil {
		return map[string]string{}
	}
	return commit.Files
}

func getStatusData() (map[string]string, map[string]string, map[string]string, error) {
	indexChan := make(chan fileResult)
	committedChan := make(chan fileResult)
	workingChan := make(chan fileResult)

	// Get index files asynchronously
	go func() {
		index, err := readIndex()
		indexChan <- fileResult{data: index, err: err}
	}()

	// Get committed files asynchronously
	go func() {
		committed := getCommittedFiles()
		committedChan <- fileResult{data: committed, err: nil}
	}()

	// Get working files asynchronously
	go func() {
		working, err := getWorkingFiles()
		workingChan <- fileResult{data: working, err: err}
	}()

	// Collect results
	indexResult := <-indexChan
	committedResult := <-committedChan
	workingResult := <-workingChan

	// Check for errors
	if indexResult.err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read index: %w", indexResult.err)
	}
	if workingResult.err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get working files: %w", workingResult.err)
	}

	return indexResult.data.(map[string]string),
		committedResult.data.(map[string]string),
		workingResult.data.(map[string]string),
		nil
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long: `Display the state of the working directory and the staging area.
Shows which files have been staged for the next commit and which files
are untracked. This helps you understand what will be included in your
next commit.`,
	Run: func(cmd *cobra.Command, args []string) {
		index, committed, working, err := getStatusData()
		if err != nil {
			fmt.Printf("Error getting status: %v\n", err)
			return
		}

		fmt.Println("=== Staged ===")
		for path, hash := range index {
			if committed[path] != hash {
				fmt.Println(path)
			}
		}

		fmt.Println("\n=== Modified but not Staged ===")
		for path, hash := range working {
			if indexHash, ok := index[path]; ok && indexHash != hash {
				fmt.Println(path)
			}
		}

		fmt.Println("\n=== Untracked Files ===")
		for path := range working {
			_, inIndex := index[path]
			_, inCommit := committed[path]
			if !inIndex && !inCommit {
				fmt.Println(path)
			}
		}

		fmt.Println("\n=== Deleted ===")
		seen := map[string]bool{}

		// Deleted files that were in the last commit
		for path := range committed {
			if _, ok := working[path]; !ok {
				fmt.Println(path + " (was saved)")
				seen[path] = true
			}
		}

		// Deleted files that were staged
		for path := range index {
			if _, ok := working[path]; !ok && !seen[path] {
				fmt.Println(path + " (was staged)")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
