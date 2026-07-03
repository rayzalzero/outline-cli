package sort

import (
	"path/filepath"
	"sort"
	"strings"
)

func SortPathsForUpload(paths []string) {
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
	isFolderIdx := IsFolderIndex(path)
	
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

