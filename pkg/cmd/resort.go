package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rayzalzero/outline-cli/pkg/api"
	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/spf13/cobra"
)

var resortCmd = &cobra.Command{
	Use:   "resort",
	Short: "Resort documents in Outline based on new sorting logic",
	Long: `Resort all documents in the collection by moving them to correct positions.
This uses the documents.move API to reorder documents without re-uploading content.`,
	RunE: runResort,
}

func init() {
	rootCmd.AddCommand(resortCmd)
}

func runResort(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	if !config.IsRepository(cwd) {
		return fmt.Errorf("not an outline repository")
	}

	apiKey := os.Getenv("OUTLINE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OUTLINE_TOKEN")
	}
	if apiKey == "" {
		return fmt.Errorf("OUTLINE_API_KEY or OUTLINE_TOKEN not set")
	}

	baseURL := os.Getenv("OUTLINE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://outline-rbi.jatismobile.com"
	}

	client := api.NewClient(baseURL, apiKey)

	manifestPath := filepath.Join(cwd, ".outline", "manifest.json")
	m, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	collectionID := cfg.CollectionID
	if collectionID == "" {
		return fmt.Errorf("no collection ID in config")
	}

	var files []string
	for relPath := range m {
		files = append(files, relPath)
	}

	sortFilesForPush(files, m)

	fmt.Printf("Resorting %d documents...\n\n", len(files))

	moved := 0
	for index, relPath := range files {
		entry, exists := m.Get(relPath)
		if !exists || entry.ID == "" {
			continue
		}

		parentID := entry.ParentID
		
		_, err = client.MoveDocumentWithCollection(entry.ID, parentID, collectionID)
		if err != nil {
			fmt.Printf("  ✗ %s (move error: %v)\n", relPath, err)
			continue
		}

		fmt.Printf("  ✓ [%3d] %s\n", index+1, relPath)
		moved++
	}

	fmt.Printf("\nResorted %d documents\n", moved)
	return nil
}
