package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	cfg, err := config.Load(cwd)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	baseURL := cfg.RemoteURL
	if baseURL == "" {
		baseURL = os.Getenv("OUTLINE_BASE_URL")
		if baseURL == "" {
			baseURL = "https://outline-rbi.jatismobile.com"
		}
	}

	client := api.NewClient(baseURL, cfg.APIKey)

	manifestPath := filepath.Join(cwd, ".outline", "manifest.json")
	m, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("load manifest: %w", err)
	}

	modified := []string{}
	newFiles := []string{}
	parentChanged := []string{}
	
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

		expectedParentID := findParentDocumentID(m, relPath)
		if entry.ParentID != "" && expectedParentID != entry.ParentID {
			parentChanged = append(parentChanged, relPath)
		}
	}

	totalChanges := len(modified) + len(newFiles) + len(parentChanged)
	
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

	if len(parentChanged) > 0 {
		fmt.Printf("Parent changed (hierarchy moved): %d\n", len(parentChanged))
		for _, path := range parentChanged {
			fmt.Printf("  moved: %s\n", path)
		}
	}

	if pushDryRun {
		fmt.Println("\n(Dry run - no changes made)")
		return nil
	}

	collectionID := cfg.CollectionID
	if collectionID == "" {
		return fmt.Errorf("no collection ID in config")
	}

	fmt.Println("\nPushing changes...")
	
	pushed := 0
	created := 0
	moved := 0
	
	for _, relPath := range parentChanged {
		entry, exists := m.Get(relPath)
		if !exists || entry.ID == "" {
			continue
		}

		expectedParentID := findParentDocumentID(m, relPath)
		if expectedParentID == entry.ParentID {
			continue
		}

		_, err = client.MoveDocument(entry.ID, expectedParentID)
		if err != nil {
			fmt.Printf("  ✗ %s (move error: %v)\n", relPath, err)
			continue
		}

		entry.ParentID = expectedParentID
		m[relPath] = entry

		fmt.Printf("  ↔ %s (moved to new parent)\n", relPath)
		moved++
	}
	
	allFiles := append(newFiles, modified...)
	sortFilesForPush(allFiles, m)
	
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
		title := getTitleFromFile(relPath, text)
		
		doc, err = client.CreateDocument(title, text, collectionID)
		if err != nil && (contains(err.Error(), "authorization_error") || contains(err.Error(), "403")) {
			fallbackParentID := findAnyDocumentID(m)
			if fallbackParentID != "" {
				doc, err = client.CreateDocumentWithParent(title, text, collectionID, fallbackParentID)
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
		
		newHash, _ := manifest.FileHash(filePath)
		entry.Hash = newHash
		entry.ID = doc.ID
		entry.Revision = doc.Revision
		entry.Updated = doc.UpdatedAt
		entry.Collection = collectionID
		entry.ParentID = doc.ParentDocumentID
		m.Set(relPath, entry)
		
		parentID := findParentDocumentID(m, relPath)
		if parentID != "" {
			doc, err = client.MoveDocument(doc.ID, parentID)
			if err != nil {
				fmt.Printf("  ✗ %s (move to parent error: %v)\n", relPath, err)
			} else {
				fmt.Printf("  ↔ %s (moved to parent)\n", relPath)
				entry.ParentID = doc.ParentDocumentID
				m.Set(relPath, entry)
			}
		}
	} else {
		expectedParentID := findParentDocumentID(m, relPath)
		if expectedParentID != entry.ParentID {
			_, err = client.MoveDocument(docID, expectedParentID)
			if err != nil {
				fmt.Printf("  ✗ %s (move error: %v)\n", relPath, err)
				continue
		}
		fmt.Printf("  ↔ %s (moved to new parent)\n", relPath)
		}
		
		title := getTitleFromFile(relPath, text)
		
		doc, err = client.UpdateDocument(docID, text, title, entry.Revision)
		if err != nil {
			fmt.Printf("  ✗ %s (update error: %v)\n", relPath, err)
			continue
		}
		
		newHash, _ := manifest.FileHash(filePath)
		entry.Hash = newHash
		entry.ID = doc.ID
		entry.Revision = doc.Revision
		entry.Updated = doc.UpdatedAt
		entry.Collection = collectionID
		entry.ParentID = doc.ParentDocumentID
		m.Set(relPath, entry)
			
		fmt.Printf("  ✓ %s (updated, revision %d)\n", relPath, doc.Revision)
	}

		pushed++
	}

	if pushed > 0 {
		if err := m.Save(manifestPath); err != nil {
			return fmt.Errorf("save manifest: %w", err)
		}
	}
	
	if created > 0 {
		fmt.Println("\nSorting documents...")
		
		var sortFiles []string
		for _, relPath := range allFiles {
			entry, exists := m.Get(relPath)
			if exists && entry.ID != "" {
				sortFiles = append(sortFiles, relPath)
			}
		}
		
		sort.Slice(sortFiles, func(i, j int) bool {
			entryA, _ := m.Get(sortFiles[i])
			entryB, _ := m.Get(sortFiles[j])
			return entryA.Index < entryB.Index
		})
		
		sorted := 0
		for idx, relPath := range sortFiles {
			entry, _ := m.Get(relPath)
			
			index := entry.Index
			if index == 0 {
				index = idx
			}
			
			parentID := findParentDocumentID(m, relPath)
			if parentID == "" {
				_, err := client.MoveDocumentWithIndex(entry.ID, "", collectionID, index)
				if err != nil {
					fmt.Printf("  ✗ %s (sort error: %v)\n", relPath, err)
				} else {
					fmt.Printf("  ↕ %s (index: %d)\n", relPath, index)
					sorted++
				}
			} else {
				_, err := client.MoveDocumentWithIndex(entry.ID, parentID, "", index)
				if err != nil {
					fmt.Printf("  ✗ %s (sort error: %v)\n", relPath, err)
				} else {
					fmt.Printf("  ↕ %s (index: %d)\n", relPath, index)
					sorted++
				}
			}
		}
		fmt.Printf("Sorted %d document(s)\n", sorted)
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
	
	fileName := filepath.Base(filePath)
	folderName := filepath.Base(dir)
	
	if fileName == folderName+".md" {
		parentDir := filepath.Dir(dir)
		if parentDir != "." && parentDir != "" {
			parentFolderName := filepath.Base(parentDir)
			parentDocPath := filepath.Join(parentDir, parentFolderName+".md")
			if entry, exists := m.Get(parentDocPath); exists && entry.ID != "" {
				return entry.ID
			}
		}
		return ""
	}
	
	parentDocPath := filepath.Join(dir, folderName+".md")
	if entry, exists := m.Get(parentDocPath); exists && entry.ID != "" {
		return entry.ID
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

func getPushPriority(path string, allFiles []string) int {
	depth := strings.Count(path, "/")
	parent := filepath.Dir(path)
	isIndex := strings.HasSuffix(path, "/index.md") || strings.HasSuffix(path, "/overview.md") || path == "index.md" || path == "overview.md"
	isFolderIdx := isFolderIndex(path)
	
	if depth == 0 {
		if isIndex {
			return 10000
		}
		return 9000
	}
	
	maxSubfolderDepth := depth
	for _, f := range allFiles {
		fParent := filepath.Dir(f)
		if strings.HasPrefix(fParent, parent+"/") {
			fDepth := strings.Count(f, "/")
			if fDepth > maxSubfolderDepth {
				maxSubfolderDepth = fDepth
			}
		}
	}
	
	if maxSubfolderDepth > depth {
		return maxSubfolderDepth*1000 + 800
	}
	
	if isIndex {
		return depth*1000 + 900
	}
	
	if isFolderIdx {
		return depth*1000 + 850
	}
	
	return depth*1000 + 800
}

func isFolderIndex(path string) bool {
	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}
	
	base := filepath.Base(path)
	folderName := filepath.Base(dir)
	
	return base == folderName+".md"
}

func sortFilesForPush(files []string, m manifest.Manifest) {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		
		entryA, okA := m.Get(a)
		entryB, okB := m.Get(b)
		
		if okA && okB {
			if entryA.Index != entryB.Index {
				return entryA.Index < entryB.Index
			}
			
			aHasParent := entryA.ParentID != ""
			bHasParent := entryB.ParentID != ""
			
			if aHasParent != bHasParent {
				return !aHasParent
			}
		}
		
		aPriority := getPushPriority(a, files)
		bPriority := getPushPriority(b, files)
		
		if aPriority != bPriority {
			return aPriority < bPriority
		}
		
		return a < b
	})
}

func getTitleFromFile(relPath, text string) string {
	filename := filepath.Base(relPath)
	
	if filename == "index.md" || filename == "overview.md" {
		return "Overview"
	}
	
	title := extractH1FromMarkdown(text)
	if title == "" {
		title = filename[:len(filename)-3]
	}
	return title
}

func extractH1FromMarkdown(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return ""
}
