package sort

import (
	"path/filepath"
	"sort"
	"strings"
)

func SortPathsByPriority(paths []string) {
	sort.Slice(paths, func(i, j int) bool {
		a, b := paths[i], paths[j]
		
		baseA := filepath.Base(a)
		baseB := filepath.Base(b)
		
		isRootIndexA := (baseA == "index.md" || baseA == "overview.md") && !strings.Contains(a, "/")
		isRootIndexB := (baseB == "index.md" || baseB == "overview.md") && !strings.Contains(b, "/")
		
		if isRootIndexA != isRootIndexB {
			return isRootIndexA
		}
		
		depthA := strings.Count(a, "/")
		depthB := strings.Count(b, "/")
		
		if depthA != depthB {
			return depthA < depthB
		}
		
		isFolderIdxA := IsFolderIndex(a)
		isFolderIdxB := IsFolderIndex(b)
		
		if isFolderIdxA != isFolderIdxB {
			return isFolderIdxA
		}
		
		return a < b
	})
}

func IsFolderIndex(path string) bool {
	dir := filepath.Dir(path)
	if dir == "." {
		return false
	}
	base := filepath.Base(path)
	folderName := filepath.Base(dir)
	return base == folderName+".md"
}
