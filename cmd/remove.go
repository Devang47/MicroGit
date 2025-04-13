package cmd

import (
	"fmt"
	"microgit/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove files from the staging area",
	Long: `Remove files from the staging area, effectively un-staging them.

Usage:
  microgit remove <file1> [file2 ...]  - Remove specific files from staging
  microgit remove .                    - Remove all files from staging

This command will:
1. Remove the specified files from the index
2. Keep the files in your working directory
3. Allow you to re-stage them later if needed`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No files specified")
			return
		}

		if args[0] == "." {
			indexPath := filepath.Join(utils.DEFAULT_PATH, "index")
			err := os.WriteFile(indexPath, []byte(""), 0644)
			if err != nil {
				fmt.Printf("Failed to unstage files: %v", err)
				return
			}

			return
		}

		for _, file := range args {
			indexPath := filepath.Join(utils.DEFAULT_PATH, "index")

			data, err := os.ReadFile(indexPath)
			if err != nil {
				fmt.Println("Failed to read index file", err)
			}

			lines := strings.Split(string(data), "\n")

			var stagedFiles []string

			for _, line := range lines {
				if !strings.HasPrefix(line, file) {
					stagedFiles = append(stagedFiles, line)
				}
			}

			newIndex := strings.Trim(strings.Join(stagedFiles, "\n"), "\n")
			os.WriteFile(indexPath, []byte(newIndex), 0644)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
