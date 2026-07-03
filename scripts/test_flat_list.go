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
	
	filepath.Walk("../testing-collection", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			relPath := strings.TrimPrefix(path, "../testing-collection/")
			files = append(files, relPath)
		}
		return nil
	})
	
	sortFilesForPush(files)
	
	fmt.Println("📄 Sorted Files (flat list):\n")
	for i, file := range files {
		fmt.Printf("%3d. %s\n", i+1, file)
	}
}

func sortFilesForPush(files []string) {
	sort.Slice(files, func(i, j int) bool {
		a, b := files[i], files[j]
		
		aPriority := getPushPriority(a, files)
		bPriority := getPushPriority(b, files)
		
		if aPriority != bPriority {
			return aPriority < bPriority
		}
		
		return a > b
	})
}

func getPushPriority(path string, allFiles []string) int {
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
	
	if isFolderIdx {
		return depth * 1000
	}
	
	if maxSubfolderDepth > depth {
		return maxSubfolderDepth*1000 + 850
	}
	
	if isIndex {
		return depth*1000 + 900
	}
	
	return depth*1000 + 800
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
