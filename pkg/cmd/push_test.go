package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/rayzalzero/outline-cli/pkg/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPushDetectsFolderMove(t *testing.T) {
	originalCwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		os.Chdir(originalCwd)
	}()

	tmpDir := t.TempDir()

	err = os.MkdirAll(filepath.Join(tmpDir, "a", "b", "c"), 0755)
	require.NoError(t, err)

	// Create document file at original location
	docPath := filepath.Join(tmpDir, "a", "b", "c", "c.md")
	docContent := `---
outline_id: c-id
outline_collection: test-collection
outline_url: /doc/c
outline_parent_id: b-id
outline_revision: 1
---

# C Document

Content here.
`
	err = ioutil.WriteFile(docPath, []byte(docContent), 0644)
	require.NoError(t, err)

	// Create .outline directory and config
	outlineDir := filepath.Join(tmpDir, ".outline")
	err = os.MkdirAll(outlineDir, 0755)
	require.NoError(t, err)

	// Create config file
	configContent := `[remote "origin"]
    url = https://outline-rbi.jatismobile.com
    collection = test-collection-id

[auth]
    token = test-token
`
	err = ioutil.WriteFile(filepath.Join(outlineDir, "config"), []byte(configContent), 0644)
	require.NoError(t, err)

	// Create manifest with document entries and parent relationships
	manifestPath := filepath.Join(outlineDir, "manifest.json")
	manifestData := map[string]interface{}{
		"a/b/b.md": map[string]interface{}{
			"id":         "b-id",
			"revision":   1,
			"hash":       "hash-b",
			"updated":    "2026-07-01T00:00:00Z",
			"collection": "test-collection-id",
			"parentId":   "a-id",
		},
		"a/b/c/c.md": map[string]interface{}{
			"id":         "c-id",
			"revision":   1,
			"hash":       "hash-c-old",
			"updated":    "2026-07-01T00:00:00Z",
			"collection": "test-collection-id",
			"parentId":   "b-id",
		},
		"a/x/x.md": map[string]interface{}{
			"id":         "x-id",
			"revision":   1,
			"hash":       "hash-x",
			"updated":    "2026-07-01T00:00:00Z",
			"collection": "test-collection-id",
			"parentId":   "a-id",
		},
	}
	manifestJSON, _ := json.MarshalIndent(manifestData, "", "  ")
	err = ioutil.WriteFile(manifestPath, manifestJSON, 0644)
	require.NoError(t, err)

	// User action: Move file from a/b/c/ to a/x/c/
	newDir := filepath.Join(tmpDir, "a", "x", "c")
	err = os.MkdirAll(newDir, 0755)
	require.NoError(t, err)

	newPath := filepath.Join(newDir, "c.md")
	err = os.Rename(docPath, newPath)
	require.NoError(t, err)

	m, err := manifest.Load(manifestPath)
	require.NoError(t, err)
	
	oldEntry := m["a/b/c/c.md"]
	delete(m, "a/b/c/c.md")
	
	newHash, err := manifest.FileHash(newPath)
	require.NoError(t, err)
	oldEntry.Hash = newHash
	
	m["a/x/c/c.md"] = oldEntry
	err = m.Save(manifestPath)
	require.NoError(t, err)

	var moveAPICalls []map[string]interface{}
	var updateAPICalls []map[string]interface{}

	// Mock Outline API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/documents.move" {
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)
			moveAPICalls = append(moveAPICalls, req)

			// Return success response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"documents": []interface{}{},
				},
			})
			return
		}

		// Mock documents.info for hash comparison
		if r.URL.Path == "/api/documents.info" {
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":               "c-id",
					"title":            "C",
					"text":             "Content here.",
					"parentDocumentId": "b-id", // Old parent
					"revision":         1,
					"updatedAt":        "2026-07-01T00:00:00Z",
				},
			})
			return
		}

		// Mock documents.update
		if r.URL.Path == "/api/documents.update" {
			var req map[string]interface{}
			json.NewDecoder(r.Body).Decode(&req)
			updateAPICalls = append(updateAPICalls, req)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":               "c-id",
					"title":            "C",
					"text":             "Content here.",
					"parentDocumentId": "b-id",
					"revision":         2,
					"updatedAt":        "2026-07-01T00:00:00Z",
				},
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	// Set environment variables for API
	oldAPIKey := os.Getenv("OUTLINE_API_KEY")
	oldBaseURL := os.Getenv("OUTLINE_BASE_URL")
	defer func() {
		if oldAPIKey == "" {
			os.Unsetenv("OUTLINE_API_KEY")
		} else {
			os.Setenv("OUTLINE_API_KEY", oldAPIKey)
		}
		if oldBaseURL == "" {
			os.Unsetenv("OUTLINE_BASE_URL")
		} else {
			os.Setenv("OUTLINE_BASE_URL", oldBaseURL)
		}
	}()

	os.Setenv("OUTLINE_API_KEY", "test-token")
	os.Setenv("OUTLINE_BASE_URL", server.URL)

	err = runPush(nil, nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(moveAPICalls), "documents.move should be called once when document is moved to new parent folder")

	moveCall := moveAPICalls[0]
	assert.Equal(t, "c-id", moveCall["id"], "Should move document c")
	assert.Equal(t, "x-id", moveCall["parentDocumentId"], "New parent should be x (x-id)")
}
