package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rayzalzero/outline-cli/pkg/manifest"
)

func TestStatusDetectsParentChange(t *testing.T) {
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Chdir(oldWd)
	}()

	tmpDir := t.TempDir()
	outlineDir := filepath.Join(tmpDir, ".outline")
	if err := os.MkdirAll(outlineDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config
	configContent := `[remote "origin"]
    url = https://outline.example.com
    collection = coll-123

[auth]
    token = test-token
`
	if err := os.WriteFile(filepath.Join(outlineDir, "config"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create initial manifest with doc in root (no parent)
	m := manifest.Manifest{
		"doc1.md": manifest.Entry{
			ID:       "doc-123",
			Revision: 1,
			Hash:     "hash1",
			ParentID: "", // Root document
		},
	}
	manifestPath := filepath.Join(outlineDir, "manifest.json")
	if err := m.Save(manifestPath); err != nil {
		t.Fatal(err)
	}

	// Create the markdown file matching manifest hash
	docPath := filepath.Join(tmpDir, "doc1.md")
	docContent := `---
outline_id: "doc-123"
outline_collection: "Test Collection"
outline_url: "https://outline.example.com/doc/doc-123"
outline_updated: "2026-07-01T10:00:00Z"
outline_revision: 1
outline_parent_id: ""
---

# Document 1

Content here.
`
	if err := os.WriteFile(docPath, []byte(docContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Update manifest hash to match file
	hash, err := manifest.FileHash(docPath)
	if err != nil {
		t.Fatal(err)
	}
	entry := m["doc1.md"]
	entry.Hash = hash
	m["doc1.md"] = entry
	if err := m.Save(manifestPath); err != nil {
		t.Fatal(err)
	}

	parentDir := filepath.Join(tmpDir, "parent")
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatal(err)
	}
	
	parentDocPath := filepath.Join(parentDir, "parent.md")
	parentContent := `---
outline_id: "parent-456"
outline_collection: "Test Collection"
outline_url: "https://outline.example.com/doc/parent-456"
outline_updated: "2026-07-01T10:00:00Z"
outline_revision: 1
---

# Parent Folder

Parent document.
`
	if err := os.WriteFile(parentDocPath, []byte(parentContent), 0644); err != nil {
		t.Fatal(err)
	}
	
	parentHash, err := manifest.FileHash(parentDocPath)
	if err != nil {
		t.Fatal(err)
	}
	
	m["parent/parent.md"] = manifest.Entry{
		ID:       "parent-456",
		Revision: 1,
		Hash:     parentHash,
		ParentID: "",
	}

	childDir := filepath.Join(parentDir, "child")
	if err := os.MkdirAll(childDir, 0755); err != nil {
		t.Fatal(err)
	}
	newPath := filepath.Join(childDir, "doc1.md")
	if err := os.Rename(docPath, newPath); err != nil {
		t.Fatal(err)
	}

	movedEntry := m["doc1.md"]
	delete(m, "doc1.md")
	
	newHash, err := manifest.FileHash(newPath)
	if err != nil {
		t.Fatal(err)
	}
	movedEntry.Hash = newHash
	m["parent/child/doc1.md"] = movedEntry
	
	if err := m.Save(manifestPath); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	statusErr := runStatus(nil, nil)

	w.Close()
	os.Stdout = oldStdout

	var buf [4096]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	if statusErr != nil {
		t.Fatalf("runStatus failed: %v", statusErr)
	}

	if !contains(output, "moved:") && !contains(output, "Parent changed") {
		t.Errorf("Expected status to show parent change, got:\n%s", output)
	}

	t.Logf("Status output:\n%s", output)
}
