package markdown

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter in a markdown file
type Frontmatter struct {
	OutlineID         string    `yaml:"outline_id"`
	OutlineCollection string    `yaml:"outline_collection"`
	OutlineURL        string    `yaml:"outline_url"`
	OutlineUpdated    time.Time `yaml:"outline_updated"`
	OutlineRevision   int       `yaml:"outline_revision"`
}

// Parse extracts frontmatter and content from markdown
func Parse(data []byte) (*Frontmatter, string, error) {
	content := string(data)

	// Check if file starts with ---
	if !strings.HasPrefix(content, "---\n") {
		return nil, content, nil
	}

	// Find end of frontmatter
	endIdx := strings.Index(content[4:], "\n---\n")
	if endIdx == -1 {
		return nil, content, nil
	}
	endIdx += 4

	// Extract frontmatter YAML
	fmYAML := content[4:endIdx]
	body := strings.TrimLeft(content[endIdx+5:], "\n")

	// Parse frontmatter
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(fmYAML), &fm); err != nil {
		return nil, "", fmt.Errorf("parse frontmatter: %w", err)
	}

	return &fm, body, nil
}

// Serialize combines frontmatter and content into markdown
func Serialize(fm *Frontmatter, content string) ([]byte, error) {
	var buf bytes.Buffer

	// Write frontmatter
	buf.WriteString("---\n")
	fmData, err := yaml.Marshal(fm)
	if err != nil {
		return nil, fmt.Errorf("marshal frontmatter: %w", err)
	}
	buf.Write(fmData)
	buf.WriteString("---\n\n")

	// Write content
	buf.WriteString(content)

	return buf.Bytes(), nil
}

// StripFrontmatter removes frontmatter and returns only content
func StripFrontmatter(data []byte) string {
	_, content, err := Parse(data)
	if err != nil {
		return string(data)
	}
	return content
}
