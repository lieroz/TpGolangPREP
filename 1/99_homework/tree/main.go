package main

import (
	"os"
	"fmt"
	"sort"
	"bytes"
	"strconv"
	"strings"
	"path/filepath"
)

type Node struct {
	Name  string
	Size  int64
	IsDir bool
	Nodes []Node
}

func (node Node) String() (result string) {
	var spaces []bool
	result += stringifyNodes(node.Nodes, spaces)
	return
}

func stringifyLine(node *Node, spaces []bool, last bool) (result string) {
	for _, space := range spaces {
		if space {
			result += "	"
		} else {
			result += "│	"
		}
	}
	indicator := "├───"
	if last {
		indicator = "└───"
	}
	if !node.IsDir {
		if strings.Compare(node.Name, "main.go") == 0 {
			result += indicator + node.Name + " (vary)\n"
		} else {
			if node.Size == 0 {
				result += indicator + node.Name + " (empty)\n"
			} else {
				size := strconv.FormatInt(node.Size, 10)
				result += indicator + node.Name + " (" + size + "b)\n"
			}
		}
	} else {
		result += indicator + node.Name + "\n"
	}
	return
}

func stringifyNodes(nodes []Node, spaces []bool) (result string) {
	for i, n := range nodes {
		last := i >= len(nodes)-1
		result += stringifyLine(&n, spaces, last)
		if len(n.Nodes) > 0 {
			spacesChild := append(spaces, last)
			result += stringifyNodes(n.Nodes, spacesChild)
		}
	}
	return
}

func PrintTree(buf *bytes.Buffer, node Node) {
	buf.WriteString(node.String())
}

func ReadFolder(dir string, printFiles bool) (root Node, err error) {
	root.Name = dir
	root.Nodes, err = readFolder(dir, printFiles)
	return
}

func readDir(dir string) ([]os.FileInfo, error) {
	f, err := os.Open(dir)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	return list, nil
}

func readFolder(dir string, printFiles bool) ([]Node, error) {
	var nodes []Node
	files, err := readDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		var child Node
		child.Name = f.Name()
		child.Size = f.Size()
		child.IsDir = f.IsDir()
		if f.IsDir() {
			newDir := filepath.Join(dir, f.Name())
			child.Nodes, err = readFolder(newDir, printFiles)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, child)
		}
		if !f.IsDir() && printFiles {
			nodes = append(nodes, child)
		}
	}
	return nodes, nil
}

func dirTree(buf *bytes.Buffer, path string, printFiles bool) error {
	obj, err := ReadFolder(path, printFiles)
	if err != nil {
		return err
	}
	PrintTree(buf, obj)
	return nil
}

func main() {
	buf := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	if err := dirTree(buf, path, printFiles); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Fprint(os.Stdout, buf)
}
