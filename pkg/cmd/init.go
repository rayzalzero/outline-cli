package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new outline repository",
	Long: `Initialize the current directory as an outline repository.

This creates a .outline/ directory with configuration and manifest files.

Example:
  outline init
`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}

	// Check if already initialized
	if config.IsRepository(cwd) {
		return fmt.Errorf("already initialized as outline repository")
	}

	// Create .outline directory
	outlineDir := filepath.Join(cwd, ".outline")
	if err := os.MkdirAll(outlineDir, 0755); err != nil {
		return fmt.Errorf("create .outline directory: %w", err)
	}

	// Get API key from environment
	apiKey := os.Getenv("OUTLINE_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: OUTLINE_API_KEY not set in environment")
		fmt.Println("Set it before running clone/pull/push commands:")
		fmt.Println("  export OUTLINE_API_KEY='ol_api_...'")
	}

	// Get base URL from environment or use default
	baseURL := os.Getenv("OUTLINE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://outline-rbi.jatismobile.com"
	}

	// Create config
	cfg := &config.Config{
		RemoteURL:        baseURL,
		ConflictStrategy: "prompt",
		APIDelay:         "300ms",
	}

	if err := cfg.Save(cwd); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Create empty manifest
	manifestPath := filepath.Join(outlineDir, "manifest.json")
	m := make(manifest.Manifest)
	if err := m.Save(manifestPath); err != nil {
		return fmt.Errorf("create manifest: %w", err)
	}

	fmt.Printf("Initialized empty Outline repository in %s/.outline/\n", cwd)
	
	if apiKey != "" {
		fmt.Println("\nYou can now:")
		fmt.Println("  outline clone <collection-url>")
		fmt.Println("  outline pull")
	}

	return nil
}
