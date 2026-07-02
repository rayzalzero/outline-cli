package manifest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrateManifestInfersParents(t *testing.T) {
	// Given: Old manifest without ParentID fields
	tmpDir := t.TempDir()

	// Create workspace structure
	err := os.MkdirAll(filepath.Join(tmpDir, "parent", "child"), 0755)
	require.NoError(t, err)

	// Create old-style manifest (no parentId fields)
	manifestPath := filepath.Join(tmpDir, ".outline", "manifest.json")
	err = os.MkdirAll(filepath.Dir(manifestPath), 0755)
	require.NoError(t, err)

	oldManifest := map[string]interface{}{
		"entries": map[string]interface{}{
			"parent/parent.md": map[string]interface{}{
				"id":         "parent-id",
				"revision":   1,
				"collection": "test-collection",
				// No parentId field - old format
			},
			"parent/child/child.md": map[string]interface{}{
				"id":         "child-id",
				"revision":   1,
				"collection": "test-collection",
				// No parentId field - old format
			},
		},
	}

	manifestJSON, _ := json.MarshalIndent(oldManifest, "", "  ")
	err = ioutil.WriteFile(manifestPath, manifestJSON, 0644)
	require.NoError(t, err)

	// When: Migrate manifest with folder structure
	err = MigrateManifest(manifestPath, tmpDir)
	require.NoError(t, err)

	// Then: Parents inferred from folder paths
	manifestData, err := ioutil.ReadFile(manifestPath)
	require.NoError(t, err)

	var newManifest map[string]interface{}
	err = json.Unmarshal(manifestData, &newManifest)
	require.NoError(t, err)

	entries := newManifest["entries"].(map[string]interface{})

	// Root document should have empty parent
	parentEntry := entries["parent/parent.md"].(map[string]interface{})
	assert.Equal(t, "", parentEntry["parentId"], "Root has no parent")

	// Child document should have parent ID inferred
	childEntry := entries["parent/child/child.md"].(map[string]interface{})
	assert.Equal(t, "parent-id", childEntry["parentId"], "Child parent inferred from folder structure")
}

func TestMigrateManifestValidatesViaAPI(t *testing.T) {
	// Given: Manifest with inferred parents (unverified)
	tmpDir := t.TempDir()

	manifestPath := filepath.Join(tmpDir, ".outline", "manifest.json")
	err := os.MkdirAll(filepath.Dir(manifestPath), 0755)
	require.NoError(t, err)

	manifest := map[string]interface{}{
		"entries": map[string]interface{}{
			"parent/parent.md": map[string]interface{}{
				"id":             "parent-id",
				"parentId":       "",
				"parentVerified": false,
			},
			"parent/child/child.md": map[string]interface{}{
				"id":             "child-id",
				"parentId":       "parent-id", // Inferred
				"parentVerified": false,
			},
			"parent/grandchild/grandchild.md": map[string]interface{}{
				"id":             "grandchild-id",
				"parentId":       "wrong-parent-id", // Incorrectly inferred
				"parentVerified": false,
			},
		},
	}

	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	err = ioutil.WriteFile(manifestPath, manifestJSON, 0644)
	require.NoError(t, err)

	// Mock Outline API for validation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/documents.info" {
			var req map[string]string
			json.NewDecoder(r.Body).Decode(&req)

			// Return correct parent info from API
			responses := map[string]string{
				"parent-id":      "",          // Root
				"child-id":       "parent-id", // Correct
				"grandchild-id":  "child-id",  // Should correct wrong inference
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"id":               req["id"],
					"parentDocumentId": responses[req["id"]],
				},
			})
		}
	}))
	defer server.Close()

	// When: Validate inferred parents against API
	err = ValidateParentsViaAPI(manifestPath, server.URL)
	require.NoError(t, err)

	// Then: Corrects mismatches, marks verified
	manifestData, err := ioutil.ReadFile(manifestPath)
	require.NoError(t, err)

	var validated map[string]interface{}
	err = json.Unmarshal(manifestData, &validated)
	require.NoError(t, err)

	entries := validated["entries"].(map[string]interface{})

	// All entries should be verified
	parentEntry := entries["parent/parent.md"].(map[string]interface{})
	assert.True(t, parentEntry["parentVerified"].(bool), "Parent verified")

	childEntry := entries["parent/child/child.md"].(map[string]interface{})
	assert.True(t, childEntry["parentVerified"].(bool), "Child verified")
	assert.Equal(t, "parent-id", childEntry["parentId"], "Child parent correct")

	// Grandchild should be CORRECTED
	grandchildEntry := entries["parent/grandchild/grandchild.md"].(map[string]interface{})
	assert.True(t, grandchildEntry["parentVerified"].(bool), "Grandchild verified")
	assert.Equal(t, "child-id", grandchildEntry["parentId"], "Grandchild parent corrected from API")
}
