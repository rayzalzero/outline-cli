package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/rayzalzero/outline-cli/pkg/markdown"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <file>...",
	Short: "Add file contents to the manifest",
	Long: `Add file contents to the manifest for tracking.

This command adds untracked files to the manifest so they can be pushed.

Examples:
  outline add newfile.md
  outline add docs/*.md
  outline add .
`,
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nothing specified, nothing added")
	}

	repoPath, err := config.FindRepository()
	if err != nil {
		return err
	}

	manifestPath := filepath.Join(repoPath, ".outline", "manifest.json")
	m, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	maxIndex := 1
	for _, entry := range m {
		if entry.Index >= maxIndex {
			maxIndex = entry.Index + 1
		}
	}

	added := 0
	for _, arg := range args {
		var filesToAdd []string

		if arg == "." {
			// Add all untracked .md files
			err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
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
				relPath, _ := filepath.Rel(repoPath, path)
				if _, exists := m[relPath]; !exists {
					filesToAdd = append(filesToAdd, path)
				}
				return nil
			})
			if err != nil {
				return err
			}
		} else {
			// Handle specific file or glob pattern
			absPath := arg
			if !filepath.IsAbs(arg) {
				absPath = filepath.Join(repoPath, arg)
			}

			matches, err := filepath.Glob(absPath)
			if err != nil {
				return err
			}

			if len(matches) == 0 {
				// Check if file exists
				if _, err := os.Stat(absPath); os.IsNotExist(err) {
					return fmt.Errorf("pathspec '%s' did not match any files", arg)
				}
				matches = []string{absPath}
			}

			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					continue
				}
				if !info.IsDir() && filepath.Ext(match) == ".md" {
					relPath, _ := filepath.Rel(repoPath, match)
					if _, exists := m[relPath]; !exists {
						filesToAdd = append(filesToAdd, match)
					}
				}
			}
		}

		sortFilesForAdd(filesToAdd, repoPath)
		
		for _, filePath := range filesToAdd {
			relPath, _ := filepath.Rel(repoPath, filePath)

			hash, err := manifest.FileHash(filePath)
			if err != nil {
				return fmt.Errorf("hash file %s: %w", relPath, err)
			}

			entry := manifest.Entry{
				Hash:  hash,
				Index: maxIndex,
			}
			maxIndex++

			// Try to parse frontmatter to get outline_id if exists
			data, err := os.ReadFile(filePath)
			if err == nil {
				if fm, _, err := markdown.Parse(data); err == nil && fm != nil && fm.OutlineID != "" {
					entry.ID = fm.OutlineID
					entry.Revision = fm.OutlineRevision
					entry.Updated = fm.OutlineUpdated
				}
			}
			
			// Calculate parent_id based on folder structure
			parentID := findParentDocumentID(m, relPath)
			entry.ParentID = parentID

			m.Set(relPath, entry)
			fmt.Printf("add '%s'\n", relPath)
			added++
		}
	}

	if added > 0 {
		if err := m.Save(manifestPath); err != nil {
			return fmt.Errorf("save manifest: %w", err)
		}
	} else {
		fmt.Println("No files added")
	}

	return nil
}

func sortFilesForAdd(files []string, repoPath string) {
	sort.Slice(files, func(i, j int) bool {
		relI, _ := filepath.Rel(repoPath, files[i])
		relJ, _ := filepath.Rel(repoPath, files[j])
		
		baseI := filepath.Base(relI)
		baseJ := filepath.Base(relJ)
		
		isRootIndexI := (baseI == "index.md" || baseI == "overview.md") && !strings.Contains(relI, "/")
		isRootIndexJ := (baseJ == "index.md" || baseJ == "overview.md") && !strings.Contains(relJ, "/")
		
		if isRootIndexI != isRootIndexJ {
			return isRootIndexI
		}
		
		depthI := strings.Count(relI, "/")
		depthJ := strings.Count(relJ, "/")
		
		if depthI != depthJ {
			return depthI < depthJ
		}
		
		isFolderIdxI := isFolderIndex(relI)
		isFolderIdxJ := isFolderIndex(relJ)
		
		if isFolderIdxI != isFolderIdxJ {
			return isFolderIdxI
		}
		
		return relI < relJ
	})
}
