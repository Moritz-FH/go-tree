package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Node struct {
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Target   string  `json:"target,omitempty"`
	Contents []*Node `json:"contents,omitempty"`
}

type Report struct {
	Type        string `json:"type"`
	Directories int    `json:"directories"`
	Files       *int   `json:"files,omitempty"`
}

func main() {
	dirOnly := flag.Bool("d", false, "Only list directories")
	maxDepth := flag.Int("l", -1, "Display depth of directory tree")
	var jsonOutput bool
	flag.BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	flag.BoolVar(&jsonOutput, "J", false, "Output results as JSON (like tree -J)")
	flag.BoolVar(&jsonOutput, "j", false, "Output results as JSON (shorthand)")
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	if jsonOutput {
		data, err := generateTreeJSON(root, *dirOnly, *maxDepth)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(string(data))
		return
	}

	fmt.Println(root)
	err := printTree(root, "", 0, *dirOnly, *maxDepth)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func buildNode(path string, displayName string, currentDepth int, dirOnly bool, maxDepth int, dirCount *int, fileCount *int) (*Node, error) {
	node := &Node{
		Type: "directory",
		Name: displayName,
	}

	if maxDepth != -1 && currentDepth >= maxDepth {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var targets []os.DirEntry
	for _, e := range entries {
		if !dirOnly || e.IsDir() {
			targets = append(targets, e)
		}
	}

	for _, entry := range targets {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			*dirCount++
			childNode, err := buildNode(entryPath, entry.Name(), currentDepth+1, dirOnly, maxDepth, dirCount, fileCount)
			if err != nil {
				return nil, err
			}
			node.Contents = append(node.Contents, childNode)
		} else if entry.Type()&os.ModeSymlink != 0 {
			*fileCount++
			childNode := &Node{
				Type: "link",
				Name: entry.Name(),
			}
			if target, err := os.Readlink(entryPath); err == nil {
				childNode.Target = target
			}
			node.Contents = append(node.Contents, childNode)
		} else {
			*fileCount++
			childNode := &Node{
				Type: "file",
				Name: entry.Name(),
			}
			node.Contents = append(node.Contents, childNode)
		}
	}

	return node, nil
}

func generateTreeJSON(root string, dirOnly bool, maxDepth int) ([]byte, error) {
	var dirCount, fileCount int
	rootNode, err := buildNode(root, root, 0, dirOnly, maxDepth, &dirCount, &fileCount)
	if err != nil {
		return nil, err
	}

	if dirCount > 0 || fileCount > 0 {
		dirCount++
	}

	report := Report{
		Type:        "report",
		Directories: dirCount,
	}
	if !dirOnly {
		report.Files = &fileCount
	}

	return json.MarshalIndent([]any{rootNode, report}, "", "  ")
}

func printTree(path string, indent string, currentDepth int, dirOnly bool, maxDepth int) error {
	if maxDepth != -1 && currentDepth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var targets []os.DirEntry
	for _, e := range entries {
		if !dirOnly || e.IsDir() {
			targets = append(targets, e)
		}
	}

	for i, entry := range targets {
		isLast := i == len(targets)-1
		connector := "├── "
		newIndent := indent + "│   "
		if isLast {
			connector = "└── "
			newIndent = indent + "    "
		}

		fmt.Print(indent + connector + entry.Name())

		if entry.IsDir() {
			fmt.Println("/")
			if err := printTree(filepath.Join(path, entry.Name()), newIndent, currentDepth+1, dirOnly, maxDepth); err != nil {
				return err
			}
		} else {
			fmt.Println()
		}
	}
	return nil
}
