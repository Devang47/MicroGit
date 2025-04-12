package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
`,
	Run: func(cmd *cobra.Command, args []string) {
		repoDir := DEFAULT_PATH
		objectsDir := DEFAULT_PATH + "/objects"

		if _, err := os.Stat(repoDir); !os.IsNotExist(err) {
			fmt.Println("\nRepository already initialized.")
			return
		}

		os.Mkdir(repoDir, 0755)
		os.Mkdir(objectsDir, 0755)

		// Create index and HEAD files
		// Staging area
		os.WriteFile(repoDir+"/index", []byte(""), 0644)
		// Pointer to the current commit
		os.WriteFile(repoDir+"/HEAD", []byte(""), 0644)

		fmt.Printf("Initialized empty SCM repository in %s/\n", DEFAULT_PATH)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
