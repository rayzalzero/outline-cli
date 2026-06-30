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
	
	// Check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("not an outline repository (no .outline/config found)")
	}
	
	// For now, return basic config with env var
	// TODO: Parse INI format config file
	cfg := &Config{
		APIKey: os.Getenv("OUTLINE_API_KEY"),
		ConflictStrategy: "prompt",
		APIDelay: "300ms",
	}
	
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("OUTLINE_API_KEY not set")
	}
	
	return cfg, nil
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
