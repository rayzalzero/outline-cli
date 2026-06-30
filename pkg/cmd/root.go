package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "outline",
	Short: "Outline CLI - Git-like workflow for Outline wiki",
	Long: `A command-line tool for managing Outline wiki documents with a Git-like workflow.

Features:
  - Clone collections to local directories
  - Pull updates from remote
  - Push local changes
  - Track working tree status
  - Handle conflicts gracefully

Example usage:
  outline clone https://outline.example.com/collection/docs docs/
  cd docs/
  outline status
  outline pull
  outline push
`,
	Version: "0.1.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: .outline/config)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet mode (errors only)")
}
