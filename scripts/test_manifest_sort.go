package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rayzalzero/outline-cli/pkg/manifest"
)

func main() {
	m := make(manifest.Manifest)
	
	filepath.Walk("../testing-collection", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			relPath := strings.TrimPrefix(path, "../testing-collection/")
			m.Set(relPath, manifest.Entry{
				ID:       "test-id-" + relPath,
				Revision: 1,
				Hash:     "hash123",
				Updated:  time.Now(),
			})
		}
		return nil
	})
	
	testFile := "/tmp/test_manifest_sorted.json"
	if err := m.Save(testFile); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("✅ Manifest saved to:", testFile)
	fmt.Println("\n📄 First 10 keys:")
	
	data, _ := os.ReadFile(testFile)
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, ".md") {
			fmt.Println(line)
			count++
			if count >= 10 {
				break
			}
		}
	}
}
