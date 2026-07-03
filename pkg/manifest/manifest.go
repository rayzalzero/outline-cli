package manifest

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	sortpkg "github.com/rayzalzero/outline-cli/pkg/sort"
)

// Entry represents a single document in the manifest
type Entry struct {
	ID         string    `json:"id"`
	Revision   int       `json:"revision"`
	Hash       string    `json:"hash"`
	Updated    time.Time `json:"updated"`
	Collection string    `json:"collection"`
	ParentID   string    `json:"parentId,omitempty"`
	Index      int       `json:"index"`
}

// Manifest tracks sync state for all documents
type Manifest map[string]Entry

// Load reads manifest from file
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(Manifest), nil
		}
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// Save writes manifest to file with sorted keys
func (m Manifest) Save(path string) error {
	paths := m.Paths()
	sortpkg.SortPaths(paths)
	
	var buf strings.Builder
	buf.WriteString("{\n")
	for i, p := range paths {
		entry := m[p]
		entryJSON, err := json.MarshalIndent(entry, "  ", "  ")
		if err != nil {
			return err
		}
		buf.WriteString("  ")
		buf.WriteString(jsonQuote(p))
		buf.WriteString(": ")
		buf.WriteString(string(entryJSON))
		if i < len(paths)-1 {
			buf.WriteString(",")
		}
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	
	return os.WriteFile(path, []byte(buf.String()), 0644)
}

// Get retrieves an entry by path
func (m Manifest) Get(path string) (Entry, bool) {
	entry, ok := m[path]
	return entry, ok
}

// Set updates or adds an entry
func (m Manifest) Set(path string, entry Entry) {
	m[path] = entry
}

// Delete removes an entry
func (m Manifest) Delete(path string) {
	delete(m, path)
}

// Paths returns all tracked file paths
func (m Manifest) Paths() []string {
	paths := make([]string, 0, len(m))
	for path := range m {
		paths = append(paths, path)
	}
	return paths
}

func jsonQuote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
