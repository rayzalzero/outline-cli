package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	var files []string
	
	// Collect all .md files
	filepath.Walk("testing-collection", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			relPath := strings.TrimPrefix(path, "testing-collection/")
			files = append(files, relPath)
		}
		return nil
	})
	
	// Sort files
	sortFilesForPush(files)
	
	// Debug: print cli-commands files
	fmt.Println("DEBUG: cli-commands files after sort:")
	for i, f := range files {
		if strings.Contains(f, "cli-commands") {
			p := getPriority(f, files)
			fmt.Printf("  [%2d] %s (priority=%d)\n", i+1, f, p)
		}
	}
	fmt.Println()
	
	// Create order map
	orderMap := make(map[string]int)
	for i, f := range files {
		orderMap[f] = i + 1
	}
	
	// Print tree with order
	fmt.Println("📁 testing-collection/")
	printTree("testing-collection", "", orderMap, true)
}

func printTree(dir string, prefix string, orderMap map[string]int, isRoot bool) {
	entries, _ := os.ReadDir(dir)
	
	// Separate dirs and files
	var dirs []os.DirEntry
	var mdFiles []os.DirEntry
	
	for _, entry := range entries {
		if entry.Name() == ".outline" {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else if strings.HasSuffix(entry.Name(), ".md") {
			mdFiles = append(mdFiles, entry)
		}
	}
	
	// Print directories first
	for i, d := range dirs {
		isLast := i == len(dirs)-1 && len(mdFiles) == 0
		connector := "├──"
		if isLast {
			connector = "└──"
		}
		
		fmt.Printf("%s%s 📁 %s/\n", prefix, connector, d.Name())
		
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		
		printTree(filepath.Join(dir, d.Name()), newPrefix, orderMap, false)
	}
	
	// Print markdown files
	for i, f := range mdFiles {
		isLast := i == len(mdFiles)-1
		connector := "├──"
		if isLast {
			connector = "└──"
		}
		
		relPath := strings.TrimPrefix(filepath.Join(dir, f.Name()), "testing-collection/")
		order := orderMap[relPath]
		
		fmt.Printf("%s%s [%2d] %s\n", prefix, connector, order, f.Name())
	}
}

func getPriority(path string, allFiles []string) int {
	depth := strings.Count(path, "/")
	parent := filepath.Dir(path)
	isIndex := strings.HasSuffix(path, "/index.md") || strings.HasSuffix(path, "/overview.md") || path == "index.md" || path == "overview.md"
	isFolderIdx := isFolderIndex(path)
	
	if depth == 0 {
		if isIndex {
			return 10000
		}
		return 9000
	}
	
	maxSubfolderDepth := depth
	for _, f := range allFiles {
		fParent := filepath.Dir(f)
		if strings.HasPrefix(fParent, parent+"/") {
			fDepth := strings.Count(f, "/")
			if fDepth > maxSubfolderDepth {
				maxSubfolderDepth = fDepth
			}
		}
	}
	
	if maxSubfolderDepth > depth {
		return maxSubfolderDepth*1000 + 800
	}
	
	if isIndex {
		return depth*1000 + 900
	}
	
	if isFolderIdx {
		return depth*1000 + 850
	}
	
	return depth*1000 + 800
}

func sortFilesForPush(files []string) {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		
		aPriority := getPriority(a, files)
		bPriority := getPriority(b, files)
		
		if aPriority != bPriority {
			return aPriority < bPriority
		}
		
		return a < b
	})
}

func isFolderIndex(path string) bool {
	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}
	
	base := filepath.Base(path)
	folderName := filepath.Base(dir)
	
	return base == folderName+".md"
}
