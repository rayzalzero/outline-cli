package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rayzalzero/outline-cli/pkg/api"
	"github.com/rayzalzero/outline-cli/pkg/config"
	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/rayzalzero/outline-cli/pkg/markdown"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <collection-id-or-url> [directory]",
	Short: "Clone a collection to a local directory",
	Long: `Clone an Outline collection to a local directory.

This downloads all documents from the specified collection and sets up
a local repository for syncing.

You can provide:
  - Collection UUID: 2e317a13-b7fa-469f-aef8-27474cf336ed
  - Full URL: https://outline.com/collection/name-2e317a13-b7fa-469f-aef8-27474cf336ed
  - Collection path: /collection/name-2e317a13-b7fa-469f-aef8-27474cf336ed

Examples:
  outline clone 2e317a13-b7fa-469f-aef8-27474cf336ed jns-docs
  outline clone https://outline-rbi.jatismobile.com/collection/jns-yY1zI9VRK3 jns-docs
  outline clone /collection/jns-yY1zI9VRK3 jns-docs
  outline clone --all ~/outline-workspace
`,
	Args: cobra.MinimumNArgs(1),
	RunE: runClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().Bool("all", false, "clone all accessible collections")
}

func extractCollectionID(input string) string {
	// If already a UUID, return as-is
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if uuidPattern.MatchString(input) {
		return input
	}

	// Extract slug from URL patterns:
	// https://outline.com/collection/slug-xyz
	// /collection/slug-xyz
	// Returns slug that will be resolved via API
	patterns := []string{
		`/collection/([^/]+)$`,  // Extract slug from path
		`^([a-z0-9-]+)$`,        // Plain slug
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(input); len(matches) > 1 {
			return matches[1]
		}
	}

	return input
}

func runClone(cmd *cobra.Command, args []string) error {
	collectionID := args[0]
	
	// Extract collection ID from URL if needed
	collectionID = extractCollectionID(collectionID)
	
	// Determine target directory
	var targetDir string
	if len(args) > 1 {
		targetDir = args[1]
	} else {
		targetDir = "." // Current directory
	}

	apiKey := os.Getenv("OUTLINE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OUTLINE_TOKEN")
	}
	if apiKey == "" {
		return fmt.Errorf("OUTLINE_API_KEY or OUTLINE_TOKEN not set")
	}

	// Get base URL
	baseURL := os.Getenv("OUTLINE_BASE_URL")
	if baseURL == "" {
		baseURL = "https://outline-rbi.jatismobile.com"
	}

	// Create API client
	client := api.NewClient(baseURL, apiKey)

	// Resolve slug to UUID if needed (collections.info accepts both)
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(collectionID) {
		collection, err := client.GetCollection(collectionID)
		if err != nil {
			return fmt.Errorf("resolve collection: %w", err)
		}
		collectionID = collection.ID
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Change to target directory
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("change directory: %w", err)
	}

	cwd, _ := os.Getwd()

	// Check if already initialized
	if config.IsRepository(cwd) {
		return fmt.Errorf("directory already initialized as outline repository")
	}

	fmt.Printf("Cloning collection %s...\n", collectionID)

	// Fetch collection info
	collection, err := client.GetCollection(collectionID)
	if err != nil {
		return fmt.Errorf("fetch collection: %w", err)
	}

	fmt.Printf("Collection: %s\n", collection.Name)

	// Create .outline directory
	outlineDir := filepath.Join(cwd, ".outline")
	if err := os.MkdirAll(outlineDir, 0755); err != nil {
		return fmt.Errorf("create .outline directory: %w", err)
	}

	// Save collection metadata
	collectionPath := filepath.Join(outlineDir, "collection.json")
	collData, _ := json.MarshalIndent(collection, "", "  ")
	if err := os.WriteFile(collectionPath, collData, 0644); err != nil {
		return fmt.Errorf("save collection metadata: %w", err)
	}

	// Create config
	cfg := &config.Config{
		RemoteURL:        baseURL,
		CollectionID:     collectionID,
		ConflictStrategy: "prompt",
		APIDelay:         "300ms",
	}
	if err := cfg.Save(cwd); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	// Fetch document tree
	fmt.Println("Fetching document tree...")
	docTree, err := client.GetCollectionDocuments(collectionID)
	if err != nil {
		return fmt.Errorf("fetch document tree: %w", err)
	}

	// Initialize manifest
	m := make(manifest.Manifest)

	// Download documents recursively
	totalDocs := 0
	var processNode func(node api.DocumentNode, parentPath string) error
	processNode = func(node api.DocumentNode, parentPath string) error {
		// Fetch full document
		doc, err := client.GetDocument(node.ID)
		if err != nil {
			return fmt.Errorf("fetch document %s: %w", node.ID, err)
		}

		// Generate file path
		slug := slugify(doc.Title)
		var filePath string
		if parentPath == "" {
			filePath = slug + ".md"
		} else {
			filePath = filepath.Join(parentPath, slug, slug+".md")
		}

		// Create parent directories
		if dir := filepath.Dir(filePath); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("create directory %s: %w", dir, err)
			}
		}

		// Create frontmatter
		fm := &markdown.Frontmatter{
			OutlineID:         doc.ID,
			OutlineCollection: collection.Name,
			OutlineURL:        doc.URL,
			OutlineUpdated:    doc.UpdatedAt,
			OutlineRevision:   doc.Revision,
		}

		// Serialize markdown with frontmatter
		fileContent, err := markdown.Serialize(fm, doc.Text)
		if err != nil {
			return fmt.Errorf("serialize markdown: %w", err)
		}

		// Write file
		if err := os.WriteFile(filePath, fileContent, 0644); err != nil {
			return fmt.Errorf("write file %s: %w", filePath, err)
		}

		// Compute hash
		hash := fmt.Sprintf("%x", md5.Sum(fileContent))

		// Add to manifest
		m.Set(filePath, manifest.Entry{
			ID:         doc.ID,
			Revision:   doc.Revision,
			Hash:       hash,
			Updated:    doc.UpdatedAt,
			Collection: collection.Name,
		})

		totalDocs++
		fmt.Printf("  [%d] %s\n", totalDocs, filePath)

		// Process children
		if len(node.Children) > 0 {
			childParent := slug
			if parentPath != "" {
				childParent = filepath.Join(parentPath, slug)
			}
			for _, child := range node.Children {
				if err := processNode(child, childParent); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Process all root documents
	for _, node := range docTree {
		if err := processNode(node, ""); err != nil {
			return err
		}
	}

	// Save manifest
	manifestPath := filepath.Join(outlineDir, "manifest.json")
	if err := m.Save(manifestPath); err != nil {
		return fmt.Errorf("save manifest: %w", err)
	}

	fmt.Printf("\nCloned %d documents from %s\n", totalDocs, collection.Name)

	return nil
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	
	// Replace spaces and special chars with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	
	// Remove leading/trailing hyphens
	s = strings.Trim(s, "-")
	
	// Limit length
	if len(s) > 50 {
		s = s[:50]
	}
	
	if s == "" {
		return "untitled"
	}
	
	return s
}
