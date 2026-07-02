package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
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
	repoPath, err := config.FindRepository()
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(repoPath, ".outline", "manifest.json")
	m, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	collectionPath := filepath.Join(repoPath, ".outline", "collection.json")
	collectionData, _ := os.ReadFile(collectionPath)
	
	fmt.Printf("On collection: %s\n", repoPath)
	if len(collectionData) > 0 {
		fmt.Printf("Collection info: .outline/collection.json\n")
	}
	fmt.Println()

	modified := []string{}
	deleted := []string{}
	untracked := []string{}
	parentChanged := []string{}
	tracked := make(map[string]bool)

	for relPath, entry := range m {
		tracked[relPath] = true
		filePath := filepath.Join(repoPath, relPath)
		
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			deleted = append(deleted, relPath)
			continue
		}

		currentHash, err := manifest.FileHash(filePath)
		if err != nil {
			continue
		}

		if currentHash != entry.Hash {
			modified = append(modified, relPath)
		}

		expectedParentID := findParentDocumentID(m, relPath)
		if expectedParentID != entry.ParentID {
			parentChanged = append(parentChanged, relPath)
		}
	}

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if info.Name() == ".outline" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return nil
		}

		if !tracked[relPath] {
			untracked = append(untracked, relPath)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("scan files: %w", err)
	}

	if len(modified) == 0 && len(deleted) == 0 && len(untracked) == 0 && len(parentChanged) == 0 {
		fmt.Println("nothing to commit, working tree clean")
		return nil
	}

	if len(parentChanged) > 0 {
		fmt.Println("Parent changed (hierarchy moved):")
		for _, path := range parentChanged {
			fmt.Printf("  moved:      %s\n", path)
		}
		fmt.Println()
	}

	if len(modified) > 0 {
		fmt.Println("Changes not pushed:")
		for _, path := range modified {
			fmt.Printf("  modified:   %s\n", path)
		}
		fmt.Println()
	}

	if len(deleted) > 0 {
		fmt.Println("Deleted files:")
		for _, path := range deleted {
			fmt.Printf("  deleted:    %s\n", path)
		}
		fmt.Println()
	}

	if len(untracked) > 0 {
		fmt.Println("Untracked files:")
		fmt.Println("  (use \"outline add <file>...\" to track)")
		fmt.Println()
		for _, path := range untracked {
			fmt.Printf("  %s\n", path)
		}
		fmt.Println()
	}

	return nil
}
