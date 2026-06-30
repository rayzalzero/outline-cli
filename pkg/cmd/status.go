package cmd

import (
	"fmt"

	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long: `Show the working tree status.

Displays which files have been modified, deleted, or are untracked.

Example:
  outline status
`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Find repository
	repoPath, err := config.FindRepository()
	if err != nil {
		return err
	}

	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Println("\nStatus command not yet implemented.")
	fmt.Println("Coming soon: show modified, deleted, and untracked files")

	return nil
}
