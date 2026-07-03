package sort

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetSortUpload(path string) map[string]int {
	var files []string
	
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".md") {
			relPath := strings.TrimPrefix(p, path+"/")
			files = append(files, relPath)
		}
		return nil
	})
	
	SortPaths(files)
	
	result := make(map[string]int)
	for i, file := range files {
		result[file] = i + 1
	}
	return result
}

func SortPaths(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		a, b := paths[i], paths[j]
		aPriority := getPushPriority(a, paths)
		bPriority := getPushPriority(b, paths)
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
