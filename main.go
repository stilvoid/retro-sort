package main

import (
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"os"
	"path/filepath"
)

/*
TODO:
* Argument for source and dest
* Argument for max size
* Argument for glob
*/

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s ([src] ([dst])) (--move)\n", os.Args[0])
	os.Exit(1)
}

var size = 100

func main() {
	if len(os.Args) != 2 {
		panic("You must supply a source directory")
	}

	src := os.Args[1]

	files, err := findFiles(src)
	if err != nil {
		panic(err)
	}

	slices.SortStableFunc(files, func(a, b string) int {
		a = filepath.Base(a)
		b = filepath.Base(b)

		return strings.Compare(a, b)
	})

	prefixSize := findMinPrefix(files)

	// Make groups, bro
	groups := makeGroups(files, prefixSize)

	for _, group := range groups {
		dstDir := filepath.Join("./dst", groupName(group, prefixSize))
		os.MkdirAll(dstDir, 0750)
		for _, fn := range group {
			dstFile := filepath.Join(dstDir, filepath.Base(fn))
			fmt.Println(dstFile)
			os.Link(fn, dstFile)
		}
	}

	fmt.Println("done")
}

func findFiles(dir string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}

		return err
	})

	return files, err
}

func findMinPrefix(in []string) int {
	for prefixSize := 1; ; prefixSize++ {
		if tryPrefix(in, prefixSize) {
			return prefixSize
		}
	}
}

func getPrefix(fn string, prefixSize int) string {
	return strings.ToLower(filepath.Base(fn)[:prefixSize])
}

func tryPrefix(in []string, prefixSize int) bool {
	seen := make(map[string]int)

	for _, fn := range in {
		prefix := getPrefix(fn, prefixSize)
		seen[prefix]++

		if seen[prefix] > size {
			return false
		}
	}

	return true
}

func getCategory(fn string) byte {
	c := strings.ToLower(filepath.Base(fn))[0]

	if c >= 'a' && c <= 'z' {
		return c
	}

	return '#'
}

func makeGroups(in []string, prefixSize int) [][]string {
	remaining := in

	groups := make([][]string, 0)

	for len(remaining) > 0 {
		if len(remaining) < size {
			groups = append(groups, remaining)
			remaining = []string{}
		} else {
			// Find the break point
			var point int

			for point = 1; point < len(remaining); point++ {
				if getCategory(remaining[point-1]) != getCategory(remaining[point]) {
					break
				}
			}

			if point > size {
				for point = size; point > 0; point-- {
					if getPrefix(remaining[point-1], prefixSize) != getPrefix(remaining[point], prefixSize) {
						break
					}
				}
			}

			groups = append(groups, remaining[:point])
			remaining = remaining[point:]
		}
	}

	return groups
}

func groupName(group []string, prefixSize int) string {
	return fmt.Sprintf("%s-%s", getPrefix(group[0], prefixSize), getPrefix(group[len(group)-1], prefixSize))
}
