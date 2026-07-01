package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds repository configuration
type Config struct {
	// Remote settings
	RemoteURL      string
	CollectionID   string
	
	// Auth
	APIKey         string
	
	// Sync settings
	AutoPull       bool
	ConflictStrategy string // prompt, force-local, force-remote
	APIDelay       string   // e.g., "300ms"
}

// Load reads config from repository .outline/config
func Load(repoPath string) (*Config, error) {
	configPath := filepath.Join(repoPath, ".outline", "config")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not an outline repository (no .outline/config found)")
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	
	apiKey := os.Getenv("OUTLINE_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OUTLINE_TOKEN")
	}
	
	cfg := &Config{
		APIKey: apiKey,
		ConflictStrategy: "prompt",
		APIDelay: "300ms",
	}
	
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OUTLINE_API_KEY or OUTLINE_TOKEN not set")
	}
	
	lines := string(data)
	for _, line := range splitLines(lines) {
		line = trimSpace(line)
		if line == "" || line[0] == '#' || line[0] == '[' {
			continue
		}
		
		if contains(line, "collection = ") {
			cfg.CollectionID = extractValue(line, "collection = ")
		} else if contains(line, "url = ") {
			cfg.RemoteURL = extractValue(line, "url = ")
		}
	}
	
	return cfg, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func extractValue(line, key string) string {
	idx := findSubstring(line, key)
	if idx < 0 {
		return ""
	}
	return trimSpace(line[idx+len(key):])
}

// Save writes config to .outline/config
func (c *Config) Save(repoPath string) error {
	configPath := filepath.Join(repoPath, ".outline", "config")
	
	// Ensure .outline directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	
	// TODO: Write proper INI format
	content := fmt.Sprintf(`[remote "origin"]
	url = %s
	collection = %s

[auth]
	token = ${OUTLINE_API_KEY}

[sync]
	auto_pull = %v
	conflict_strategy = %s
	api_delay = %s
`, c.RemoteURL, c.CollectionID, c.AutoPull, c.ConflictStrategy, c.APIDelay)
	
	return os.WriteFile(configPath, []byte(content), 0644)
}

// IsRepository checks if directory is an outline repository
func IsRepository(path string) bool {
	configPath := filepath.Join(path, ".outline", "config")
	_, err := os.Stat(configPath)
	return err == nil
}

// FindRepository searches for .outline directory in current and parent dirs
func FindRepository() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	dir := cwd
	for {
		if IsRepository(dir) {
			return dir, nil
		}
		
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root
			break
		}
		dir = parent
	}
	
	return "", fmt.Errorf("not an outline repository (or any of the parent directories)")
}
