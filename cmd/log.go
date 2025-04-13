package cmd

import (
	"encoding/json"
	"fmt"
	"microgit/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func readCommit(hash string) (utils.SavePoint, error) {
	objectPath := filepath.Join(utils.DEFAULT_PATH, "objects", hash)

	data, err := os.ReadFile(objectPath)
	if err != nil {
		return utils.SavePoint{}, fmt.Errorf("could not read commit object: %w", err)
	}

	var commit utils.SavePoint
	err = json.Unmarshal(data, &commit)
	if err != nil {
		return utils.SavePoint{}, fmt.Errorf("failed to parse commit JSON: %w", err)
	}
	return commit, nil
}

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show the commit history",
	Long: `Display the commit history in chronological order, starting from the most recent commit.
For each commit, it shows:
- The commit hash
- The timestamp
- The commit message
- The list of files that were modified`,
	Run: func(cmd *cobra.Command, args []string) {
		head := getHead()
		if head == "" {
			fmt.Println("No commits yet.")
		}

		current := head
		for current != "" {
			commit, err := readCommit(current)
			if err != nil {
				fmt.Println("Error reading commit:", err)
			}

			fmt.Printf("Commit: %s\n", current)
			fmt.Printf("Date: %s\n", commit.Timestamp)
			fmt.Printf("Message: %s\n", commit.Message)
			fmt.Print("Files modified: ")
			for key := range commit.Files {
				fmt.Printf("%s ", key)
			}

			fmt.Print("\n\n")

			current = commit.Parent
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
