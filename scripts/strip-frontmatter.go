package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rayzalzero/outline-cli/pkg/markdown"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run strip-frontmatter.go <directory>")
		os.Exit(1)
	}

	dir := os.Args[1]

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("✗ %s: read error: %v\n", path, err)
			return nil
		}

		content := markdown.StripFrontmatter(data)

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("✗ %s: write error: %v\n", path, err)
			return nil
		}

		fmt.Printf("✓ %s\n", path)
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
