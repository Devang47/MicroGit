package cmd

import (
	"encoding/json"
	"fmt"
	"microgit/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Switch to a specific commit",
	Long: `Switch to a specific commit in the repository history.

Usage:
  microgit checkout <commit-hash>  - Switch to a specific commit
  microgit checkout latest        - Switch to the most recent commit

This command will:
1. Restore all files to their state at the specified commit
2. Update the HEAD reference to point to the checked out commit
3. Preserve the commit history for future operations`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No commit specified")
			return
		}

		savePointHash := args[0]

		if savePointHash == "latest" {
			latestPath := filepath.Join(utils.DEFAULT_PATH, "LATEST")
			commit, err := os.ReadFile(latestPath)
			if err != nil {
				fmt.Println("Error reading HEAD:", err)
				return
			}

			savePointHash = string(commit)
		}

		savePointPath := filepath.Join(utils.DEFAULT_PATH, "objects", savePointHash)
		data, err := os.ReadFile(savePointPath)
		if err != nil {
			fmt.Printf("savePoint %s not found", savePointHash)
			return
		}

		var savePoint utils.SavePoint
		if err := json.Unmarshal(data, &savePoint); err != nil {
			fmt.Printf("invalid savePoint format: %v", err)
			return
		}

		for path, hash := range savePoint.Files {
			blobPath := filepath.Join(utils.DEFAULT_PATH, "objects", hash)

			content, err := os.ReadFile(blobPath)
			if err != nil {
				fmt.Printf("missing object for file %s", path)
				return
			}

			err = os.WriteFile(path, content, 0644)
			if err != nil {
				fmt.Printf("failed to restore file %s", path)
				return
			}
		}

		headPath := filepath.Join(utils.DEFAULT_PATH, "HEAD")
		err = os.WriteFile(headPath, []byte(savePointHash), 0644)
		if err != nil {
			fmt.Printf("Expected HEAD file to be created")
		}

		fmt.Printf("Successfully checked out commit %s\n", savePointHash)
	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
