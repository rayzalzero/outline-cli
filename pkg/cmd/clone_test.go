package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCloneStoresParentInManifest(t *testing.T) {
	// Given: Mock Outline API returning parent-child tree
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock collections.info endpoint
		if r.URL.Path == "/api/collections.info" {
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ok": true,
				"data": map[string]interface{}{
					"id":   req["id"],
					"name": "Test Collection",
				},
			})
			return
		}

		// Mock collections.documents endpoint
		if r.URL.Path == "/api/collections.documents" {
			response := map[string]interface{}{
				"ok": true,
				"data": []map[string]interface{}{
					{
						"id":    "parent-id",
						"title": "Parent",
						"url":   "/doc/parent",
						"children": []map[string]interface{}{
							{
								"id":    "child-id",
								"title": "Child",
								"url":   "/doc/child",
								"children": []interface{}{},
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Mock documents.info endpoint
		if r.URL.Path == "/api/documents.info" {
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)

			if req["id"] == "parent-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"ok": true,
					"data": map[string]interface{}{
						"document": map[string]interface{}{
							"id":               "parent-id",
							"title":            "Parent",
							"text":             "Parent content",
							"parentDocumentId": "",
						},
					},
				})
			} else if req["id"] == "child-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"ok": true,
					"data": map[string]interface{}{
						"document": map[string]interface{}{
							"id":               "child-id",
							"title":            "Child",
							"text":             "Child content",
							"parentDocumentId": "parent-id",
						},
					},
				})
			}
		}
	}))
	defer server.Close()

	// When: Clone collection
	tmpDir := t.TempDir()

	os.Setenv("OUTLINE_BASE_URL", server.URL)
	os.Setenv("OUTLINE_API_KEY", "test-key")
	defer func() {
		os.Unsetenv("OUTLINE_BASE_URL")
		os.Unsetenv("OUTLINE_API_KEY")
	}()

	err := runClone(nil, []string{"test-collection-id", tmpDir})
	require.NoError(t, err, "Clone should succeed")

	// Then: Manifest contains parent relationships
	manifestPath := filepath.Join(tmpDir, ".outline", "manifest.json")
	manifestData, err := os.ReadFile(manifestPath)
	require.NoError(t, err, "Manifest file should exist")

	var entries map[string]interface{}
	err = json.Unmarshal(manifestData, &entries)
	require.NoError(t, err, "Manifest should be valid JSON")

	// Root document should have no parent
	parentEntry := entries["parent.md"].(map[string]interface{})
	parentID, hasParent := parentEntry["parentId"]
	assert.False(t, hasParent || parentID == "", "Root document has no parent (field omitted or empty)")

	// Child document should have parent ID
	childEntry := entries["parent/child/child.md"].(map[string]interface{})
	assert.Equal(t, "parent-id", childEntry["parentId"], "Child has parent ID")
}

func TestCloneStoresParentInFrontmatter(t *testing.T) {
	// Given: Mock Outline API returning parent-child tree
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock collections.info endpoint
		if r.URL.Path == "/api/collections.info" {
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"ok": true,
				"data": map[string]interface{}{
					"id":   req["id"],
					"name": "Test Collection",
				},
			})
			return
		}

		// Mock collections.documents endpoint
		if r.URL.Path == "/api/collections.documents" {
			response := map[string]interface{}{
				"ok": true,
				"data": []map[string]interface{}{
					{
						"id":    "parent-id",
						"title": "Parent",
						"url":   "/doc/parent",
						"children": []map[string]interface{}{
							{
								"id":    "child-id",
								"title": "Child",
								"url":   "/doc/child",
								"children": []interface{}{},
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Mock documents.info endpoint
		if r.URL.Path == "/api/documents.info" {
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)

			if req["id"] == "parent-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"ok": true,
					"data": map[string]interface{}{
						"document": map[string]interface{}{
							"id":               "parent-id",
							"title":            "Parent",
							"text":             "Parent content",
							"parentDocumentId": "",
						},
					},
				})
			} else if req["id"] == "child-id" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"ok": true,
					"data": map[string]interface{}{
						"document": map[string]interface{}{
							"id":               "child-id",
							"title":            "Child",
							"text":             "Child content",
							"parentDocumentId": "parent-id",
						},
					},
				})
			}
		}
	}))
	defer server.Close()

	// When: Clone collection
	tmpDir := t.TempDir()

	os.Setenv("OUTLINE_BASE_URL", server.URL)
	os.Setenv("OUTLINE_API_KEY", "test-key")
	defer func() {
		os.Unsetenv("OUTLINE_BASE_URL")
		os.Unsetenv("OUTLINE_API_KEY")
	}()

	err := runClone(nil, []string{"test-collection-id", tmpDir})
	require.NoError(t, err, "Clone should succeed")

	// Then: Markdown frontmatter contains parent ID
	childPath := filepath.Join(tmpDir, "parent", "child", "child.md")
	childContent, err := os.ReadFile(childPath)
	require.NoError(t, err, "Child document should exist")

	// Parse YAML frontmatter
	contentStr := string(childContent)
	require.Contains(t, contentStr, "---", "Should have YAML frontmatter")

	// Extract frontmatter (between first two ---)
	parts := strings.SplitN(contentStr, "---", 3)
	require.Equal(t, 3, len(parts), "Should have frontmatter delimiters")

	frontmatterYAML := parts[1]

	var frontmatter map[string]interface{}
	err = yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter)
	require.NoError(t, err, "Frontmatter should be valid YAML")

	// Parent document should have no outline_parent_id
	parentPath := filepath.Join(tmpDir, "parent", "parent.md")
	parentContent, _ := os.ReadFile(parentPath)
	parentFrontmatter := parseFrontmatter(t, string(parentContent))
	parentID, hasParent := parentFrontmatter["outline_parent_id"]
	assert.False(t, hasParent || parentID == "", "Root document has no parent in frontmatter (field omitted or empty)")

	// Child document should have parent ID in frontmatter
	assert.Equal(t, "parent-id", frontmatter["outline_parent_id"], "Child has parent ID in frontmatter")
}

// parseFrontmatter extracts and parses YAML frontmatter from markdown content
func parseFrontmatter(t *testing.T, content string) map[string]interface{} {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return map[string]interface{}{}
	}

	var frontmatter map[string]interface{}
	err := yaml.Unmarshal([]byte(parts[1]), &frontmatter)
	require.NoError(t, err)
	return frontmatter
}
