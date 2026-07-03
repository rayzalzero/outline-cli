package manifest

import (
	sortpkg "github.com/rayzalzero/outline-cli/pkg/sort"
)

func (m Manifest) Reindex() {
	paths := m.Paths()
	sortpkg.SortPathsByPriority(paths)
	
	for i, path := range paths {
		entry := m[path]
		entry.Index = i + 1
		m[path] = entry
	}
}
