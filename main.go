package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dirOnly := flag.Bool("d", false, "Only list directories")
	maxDepth := flag.Int("l", -1, "Display depth of directory tree")
	flag.Parse()

	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	fmt.Println(root)
	err := printTree(root, "", 0, *dirOnly, *maxDepth)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
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
