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
	newFiles := []string{}
	
	for relPath, entry := range m {
		filePath := filepath.Join(cwd, relPath)
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		if entry.ID == "" {
			newFiles = append(newFiles, relPath)
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

	totalChanges := len(modified) + len(newFiles)
	
	if totalChanges == 0 {
		fmt.Println("Nothing to push")
		return nil
	}

	if len(newFiles) > 0 {
		fmt.Printf("New files: %d\n", len(newFiles))
		for _, path := range newFiles {
			fmt.Printf("  + %s\n", path)
		}
	}
	
	if len(modified) > 0 {
		fmt.Printf("Modified files: %d\n", len(modified))
		for _, path := range modified {
			fmt.Printf("  ~ %s\n", path)
		}
	}

	if pushDryRun {
		fmt.Println("\n(Dry run - no changes made)")
		return nil
	}

	cfg, err := config.Load(cwd)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	collectionID := cfg.CollectionID
	if collectionID == "" {
		return fmt.Errorf("no collection ID in config")
	}

	fmt.Println("\nPushing changes...")
	
	pushed := 0
	created := 0
	
	allFiles := append(newFiles, modified...)
	
	for _, relPath := range allFiles {
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

		docID := ""
		if frontmatter != nil {
			docID = frontmatter.OutlineID
		}

		var doc *api.Document
		
		if docID == "" {
			title := filepath.Base(relPath)
			title = title[:len(title)-3]
			
			parentID := findParentDocumentID(m, relPath)
			
			if parentID != "" {
				doc, err = client.CreateDocumentWithParent(title, text, collectionID, parentID)
			} else {
				doc, err = client.CreateDocument(title, text, collectionID)
				if err != nil && (contains(err.Error(), "authorization_error") || contains(err.Error(), "403")) {
					fallbackParentID := findAnyDocumentID(m)
					if fallbackParentID != "" {
						doc, err = client.CreateDocumentWithParent(title, text, collectionID, fallbackParentID)
					}
				}
			}
			
			if err != nil {
				if contains(err.Error(), "authorization_error") || contains(err.Error(), "403") {
					fmt.Printf("  ✗ %s (create failed: no permission to create documents in this collection)\n", relPath)
				} else {
					fmt.Printf("  ✗ %s (create error: %v)\n", relPath, err)
				}
				continue
			}
			
			newContent := fmt.Sprintf("---\noutline_id: %s\noutline_collection: %s\noutline_url: %s\noutline_updated: %s\noutline_revision: %d\n---\n\n%s",
				doc.ID, collectionID, doc.URL, doc.UpdatedAt.Format("2006-01-02T15:04:05.000Z"), doc.Revision, text)
			
			if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
				fmt.Printf("  ✗ %s (write frontmatter error: %v)\n", relPath, err)
				continue
			}
			
			fmt.Printf("  ✓ %s (created, revision %d)\n", relPath, doc.Revision)
			created++
		} else {
			doc, err = client.UpdateDocument(docID, text, entry.Revision)
			if err != nil {
				fmt.Printf("  ✗ %s (update error: %v)\n", relPath, err)
				continue
			}
			
			fmt.Printf("  ✓ %s (updated, revision %d)\n", relPath, doc.Revision)
		}

		newHash, _ := manifest.FileHash(filePath)
		entry.Hash = newHash
		entry.ID = doc.ID
		entry.Revision = doc.Revision
		entry.Updated = doc.UpdatedAt
		entry.Collection = collectionID
		m.Set(relPath, entry)

		pushed++
	}

	if pushed > 0 {
		if err := m.Save(manifestPath); err != nil {
			return fmt.Errorf("save manifest: %w", err)
		}
	}

	if created > 0 {
		fmt.Printf("\nCreated %d new document(s), updated %d\n", created, pushed-created)
	} else {
		fmt.Printf("\nPushed %d/%d files\n", pushed, totalChanges)
	}
	
	return nil
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func findParentDocumentID(m manifest.Manifest, filePath string) string {
	dir := filepath.Dir(filePath)
	if dir == "." {
		return ""
	}
	
	for path, entry := range m {
		entryDir := filepath.Dir(path)
		if entryDir == dir && entry.ID != "" {
			return entry.ID
		}
	}
	
	return ""
}

func findAnyDocumentID(m manifest.Manifest) string {
	for _, entry := range m {
		if entry.ID != "" {
			return entry.ID
		}
	}
	return ""
}
