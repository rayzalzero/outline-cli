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
	
	sortFilesForPush(files)
	
	fmt.Println("📄 Upload Order (simulating Outline tree - newest at top):\n")
	
	type Node struct {
		path     string
		order    int
		children []*Node
	}
	
	root := &Node{path: "", children: []*Node{}}
	nodeMap := make(map[string]*Node)
	nodeMap[""] = root
	
	for i, file := range files {
		parts := strings.Split(file, "/")
		currentPath := ""
		
		for j, part := range parts {
			if j > 0 {
				currentPath += "/"
			}
			currentPath += part
			
			if _, exists := nodeMap[currentPath]; !exists {
				node := &Node{path: currentPath, order: i + 1, children: []*Node{}}
				nodeMap[currentPath] = node
				
				parentPath := filepath.Dir(currentPath)
				if parentPath == "." {
					parentPath = ""
				}
				if parent, ok := nodeMap[parentPath]; ok {
					parent.children = append(parent.children, node)
				}
			}
		}
	}
	
	var printTree func(node *Node, prefix string, isLast bool)
	printTree = func(node *Node, prefix string, isLast bool) {
		if node.path != "" {
			connector := "└── "
			if !isLast {
				connector = "├── "
			}
			
			name := filepath.Base(node.path)
			reverseOrder := len(files) - node.order + 1
			fmt.Printf("%s%s[%2d] %s\n", prefix, connector, reverseOrder, name)
			
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			
			for i := len(node.children) - 1; i >= 0; i-- {
				printTree(node.children[i], newPrefix, i == 0)
			}
		} else {
			for i := len(node.children) - 1; i >= 0; i-- {
				printTree(node.children[i], "", i == 0)
			}
		}
	}
	
	printTree(root, "", true)
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
