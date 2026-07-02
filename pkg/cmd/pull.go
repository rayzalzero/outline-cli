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

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from Outline",
	Long: `Pull changes from the remote Outline instance.

This command:
  1. Fetches document updates from Outline
  2. Detects hierarchy changes (parent moves)
  3. Reorganizes local files to match Outline tree structure
  4. Updates content and manifest

Examples:
  outline pull
  outline pull --force
`,
	RunE: runPull,
}

var pullForce bool

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.Flags().BoolVar(&pullForce, "force", false, "overwrite local changes")
}

func runPull(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repoPath, err := config.FindRepository()
	if err != nil {
		return err
	}

	if repoPath != cwd {
		if err := os.Chdir(repoPath); err != nil {
			return err
		}
		cwd = repoPath
	}

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

	client := api.NewClient(cfg.RemoteURL, cfg.APIKey)

	fmt.Println("Fetching updates from Outline...")

	updated := 0
	hierarchyMoved := 0

	for relPath, entry := range m {
		if entry.ID == "" {
			continue
		}

		doc, err := client.GetDocument(entry.ID)
		if err != nil {
			fmt.Printf("  ✗ %s (fetch error: %v)\n", relPath, err)
			continue
		}

		filePath := filepath.Join(cwd, relPath)
		
		contentChanged := doc.Revision > entry.Revision
		parentChanged := doc.ParentDocumentID != entry.ParentID

		if !contentChanged && !parentChanged {
			continue
		}

		if parentChanged {
			newPath := calculateNewPath(m, doc, relPath)
			if newPath != relPath {
				if err := moveLocalFile(cwd, relPath, newPath); err != nil {
					fmt.Printf("  ✗ %s (move error: %v)\n", relPath, err)
					continue
				}
				
				delete(m, relPath)
				relPath = newPath
				filePath = filepath.Join(cwd, newPath)
				
				fmt.Printf("  ↔ %s (hierarchy moved)\n", newPath)
				hierarchyMoved++
			}
		}

		if contentChanged {
			frontmatter := &markdown.Frontmatter{
				OutlineID:         doc.ID,
				OutlineCollection: collectionID,
				OutlineURL:        doc.URL,
				OutlineUpdated:    doc.UpdatedAt,
				OutlineRevision:   doc.Revision,
				OutlineParentID:   doc.ParentDocumentID,
			}

			newContent, err := markdown.Serialize(frontmatter, doc.Text)
			if err != nil {
				fmt.Printf("  ✗ %s (serialize error: %v)\n", relPath, err)
				continue
			}

			if !pullForce {
				localHash, err := manifest.FileHash(filePath)
				if err == nil && localHash != entry.Hash {
					fmt.Printf("  ✗ %s (conflict: local changes exist, use --force to overwrite)\n", relPath)
					continue
				}
			}

			if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
				fmt.Printf("  ✗ %s (write error: %v)\n", relPath, err)
				continue
			}

			fmt.Printf("  ✓ %s (updated to revision %d)\n", relPath, doc.Revision)
			updated++
		}

		newHash, _ := manifest.FileHash(filePath)
		m[relPath] = manifest.Entry{
			ID:         doc.ID,
			Revision:   doc.Revision,
			Hash:       newHash,
			Updated:    doc.UpdatedAt,
			Collection: collectionID,
			ParentID:   doc.ParentDocumentID,
		}
	}

	if err := m.Save(manifestPath); err != nil {
		return fmt.Errorf("save manifest: %w", err)
	}

	if updated == 0 && hierarchyMoved == 0 {
		fmt.Println("Already up to date")
		return nil
	}

	fmt.Printf("\nPulled %d updates", updated)
	if hierarchyMoved > 0 {
		fmt.Printf(", %d hierarchy moves", hierarchyMoved)
	}
	fmt.Println()

	return nil
}

func calculateNewPath(m manifest.Manifest, doc *api.Document, currentPath string) string {
	if doc.ParentDocumentID == "" {
		return filepath.Base(currentPath)
	}

	for path, entry := range m {
		if entry.ID == doc.ParentDocumentID {
			parentDir := filepath.Dir(path)
			fileName := filepath.Base(currentPath)
			return filepath.Join(parentDir, fileName)
		}
	}

	return currentPath
}

func moveLocalFile(baseDir, oldPath, newPath string) error {
	oldFullPath := filepath.Join(baseDir, oldPath)
	newFullPath := filepath.Join(baseDir, newPath)

	newDir := filepath.Dir(newFullPath)
	if err := os.MkdirAll(newDir, 0755); err != nil {
		return err
	}

	return os.Rename(oldFullPath, newFullPath)
}
