package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rayzalzero/outline-cli/pkg/api"
	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/rayzalzero/outline-cli/pkg/markdown"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local changes to Outline",
	Long: `Push local changes to the remote Outline instance.

This command:
  1. Scans for modified markdown files
  2. Updates documents via Outline API
  3. Updates manifest with new revisions

Examples:
  outline push
  outline push --dry-run
  outline push --force
`,
	RunE: runPush,
}

var (
	pushDryRun bool
	pushForce  bool
)

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVar(&pushDryRun, "dry-run", false, "show what would be pushed without actually pushing")
	pushCmd.Flags().BoolVar(&pushForce, "force", false, "force push even if remote is newer")
}

func runPush(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	if !config.IsRepository(cwd) {
		return fmt.Errorf("not an outline repository (no .outline directory found)")
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

	modified := []string{}
	
	for relPath, entry := range m {
		filePath := filepath.Join(cwd, relPath)
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		currentHash, err := manifest.FileHash(filePath)
		if err != nil {
			continue
		}

		if currentHash != entry.Hash {
			modified = append(modified, relPath)
		}
	}

	if len(modified) == 0 {
		fmt.Println("Nothing to push")
		return nil
	}

	fmt.Printf("Modified files: %d\n", len(modified))
	for _, path := range modified {
		fmt.Printf("  - %s\n", path)
	}

	if pushDryRun {
		fmt.Println("\n(Dry run - no changes made)")
		return nil
	}

	fmt.Println("\nPushing changes...")
	
	pushed := 0
	for _, relPath := range modified {
		filePath := filepath.Join(cwd, relPath)
		entry, _ := m.Get(relPath)

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("  ✗ %s (read error: %v)\n", relPath, err)
			continue
		}

		frontmatter, text, err := markdown.Parse(content)
		if err != nil {
			fmt.Printf("  ✗ %s (parse error: %v)\n", relPath, err)
			continue
		}

		if frontmatter == nil {
			fmt.Printf("  ✗ %s (no frontmatter found)\n", relPath)
			continue
		}

		docID := frontmatter.OutlineID
		if docID == "" {
			fmt.Printf("  ✗ %s (no outline_id in frontmatter)\n", relPath)
			continue
		}

		doc, err := client.UpdateDocument(docID, text, entry.Revision)
		if err != nil {
			fmt.Printf("  ✗ %s (API error: %v)\n", relPath, err)
			continue
		}

		newHash, _ := manifest.FileHash(filePath)
		entry.Hash = newHash
		entry.Revision = doc.Revision
		entry.Updated = doc.UpdatedAt
		m.Set(relPath, entry)

		fmt.Printf("  ✓ %s (revision %d)\n", relPath, doc.Revision)
		pushed++
	}

	if pushed > 0 {
		if err := m.Save(manifestPath); err != nil {
			return fmt.Errorf("save manifest: %w", err)
		}
	}

	fmt.Printf("\nPushed %d/%d files\n", pushed, len(modified))
	return nil
}
