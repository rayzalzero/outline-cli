package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/rayzalzero/outline-cli/pkg/markdown"
)

// TestHierarchySyncIntegration tests the complete folder hierarchy sync workflow
func TestHierarchySyncIntegration(t *testing.T) {
	// Create test directory
	tmpDir := t.TempDir()
	
	// Setup: Create initial structure
	//   root.md
	//   parent/
	//     child1.md
	//     child2.md
	
	// Step 1: Create manifest with parent relationships
	m := make(manifest.Manifest)
	m["root.md"] = manifest.Entry{
		ID:         "root-id",
		Revision:   1,
		Hash:       "hash1",
		Updated:    time.Now(),
		Collection: "test-collection",
		ParentID:   "", // root has no parent
	}
	m["parent/child1.md"] = manifest.Entry{
		ID:         "child1-id",
		Revision:   1,
		Hash:       "hash2",
		Updated:    time.Now(),
		Collection: "test-collection",
		ParentID:   "root-id", // child1's parent is root
	}
	m["parent/child2.md"] = manifest.Entry{
		ID:         "child2-id",
		Revision:   1,
		Hash:       "hash3",
		Updated:    time.Now(),
		Collection: "test-collection",
		ParentID:   "root-id", // child2's parent is root
	}
	
	// Create actual files
	os.MkdirAll(filepath.Join(tmpDir, "parent"), 0755)
	
	writeTestFile(t, filepath.Join(tmpDir, "root.md"), "root-id", "", "# Root")
	writeTestFile(t, filepath.Join(tmpDir, "parent", "child1.md"), "child1-id", "root-id", "# Child 1")
	writeTestFile(t, filepath.Join(tmpDir, "parent", "child2.md"), "child2-id", "root-id", "# Child 2")
	
	// Save manifest
	manifestPath := filepath.Join(tmpDir, ".outline", "manifest.json")
	os.MkdirAll(filepath.Dir(manifestPath), 0755)
	if err := m.Save(manifestPath); err != nil {
		t.Fatalf("Failed to save manifest: %v", err)
	}
	
	// Step 2: Simulate folder move (user moves child1.md to different parent)
	// Move parent/child1.md → newparent/child1.md
	os.MkdirAll(filepath.Join(tmpDir, "newparent"), 0755)
	os.Rename(
		filepath.Join(tmpDir, "parent", "child1.md"),
		filepath.Join(tmpDir, "newparent", "child1.md"),
	)
	
	// Update manifest to reflect the move
	// Delete old entry
	delete(m, "parent/child1.md")
	// Add new entry at new location (simulating what push would do)
	m["newparent/child1.md"] = manifest.Entry{
		ID:         "child1-id",
		Revision:   1,
		Hash:       "hash2",
		Updated:    time.Now(),
		Collection: "test-collection",
		ParentID:   "root-id", // Parent should be detected from folder structure
	}
	
	// Step 3: Test findParentDocumentID detection
	// Since "newparent/" folder has no document yet, parent should be "" (root level)
	// Let's create newparent.md first
	writeTestFile(t, filepath.Join(tmpDir, "newparent.md"), "newparent-id", "", "# New Parent")
	m["newparent.md"] = manifest.Entry{
		ID:         "newparent-id",
		Revision:   1,
		Hash:       "hash4",
		Updated:    time.Now(),
		Collection: "test-collection",
		ParentID:   "",
	}
	
	// Now child1's parent should be detected as newparent-id
	detectedParent := findParentDocumentID(m, "newparent/child1.md")
	if detectedParent != "newparent-id" {
		t.Errorf("Expected parent 'newparent-id', got '%s'", detectedParent)
	}
	
	// Step 4: Verify status would detect the change
	oldParent := m["newparent/child1.md"].ParentID // Still "root-id" in manifest
	if oldParent == detectedParent {
		t.Errorf("Parent should have changed: old='%s', new='%s'", oldParent, detectedParent)
	}
	
	t.Logf("✅ Hierarchy detection works:")
	t.Logf("   Old parent: %s", oldParent)
	t.Logf("   New parent (from folder): %s", detectedParent)
}

func writeTestFile(t *testing.T, path, id, parentID, content string) {
	fm := &markdown.Frontmatter{
		OutlineID:         id,
		OutlineCollection: "test-collection",
		OutlineURL:        "https://example.com/doc/" + id,
		OutlineUpdated:    time.Now(),
		OutlineRevision:   1,
		OutlineParentID:   parentID,
	}
	
	data, err := markdown.Serialize(fm, content)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
}
