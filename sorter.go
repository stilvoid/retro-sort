package retrosort

import (
	"fmt"
	"slices"
	"strings"

	"os"
	"path/filepath"
)

/*
TODO:
* Break into a-z folders if initial sort would make too many groups (recursive)
*/

func Execute(src, dst string, size int, pattern string, upperCase bool, printOnly bool) {
	s := sorter{src, dst, size, pattern, upperCase, printOnly}
	s.execute()
}

type sorter struct {
	src       string
	dst       string
	size      int
	pattern   string
	upperCase bool
	printOnly bool
}

func (s sorter) execute() {
	files, err := findFiles(s.src, s.pattern)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d files\n", len(files))

	slices.SortStableFunc(files, func(a, b string) int {
		a = filepath.Base(a)
		b = filepath.Base(b)

		return strings.Compare(a, b)
	})

	prefixSize := s.findMinPrefix(files)

	// Make groups, bro
	groups := s.makeGroups(files, prefixSize)

	counter := 0
	var div int = len(files) / 20

	for _, group := range groups {
		gName := groupName(group, prefixSize)

		if s.upperCase {
			gName = strings.ToUpper(gName)
		}

		dstDir := filepath.Join(s.dst, gName)

		if !s.printOnly {
			os.MkdirAll(dstDir, 0750)
		}

		for _, fn := range group {
			dstFile := filepath.Join(dstDir, filepath.Base(fn))

			if s.printOnly {
				fmt.Printf("%s -> %s\n", fn, dstFile)
			} else {
				if err := copyFile(fn, dstFile); err != nil {
					panic(err)
				}

				counter++
				if counter%div == 0 {
					fmt.Print(".")
				}
			}
		}
	}

	fmt.Println("done")
}

func (s sorter) findMinPrefix(in []string) int {
	for prefixSize := 1; ; prefixSize++ {
		if s.tryPrefix(in, prefixSize) {
			return prefixSize
		}
	}
}

func (s sorter) tryPrefix(in []string, prefixSize int) bool {
	seen := make(map[string]int)

	for _, fn := range in {
		prefix := getPrefix(fn, prefixSize)
		seen[prefix]++

		if seen[prefix] > s.size {
			return false
		}
	}

	return true
}

func (s sorter) makeGroups(in []string, prefixSize int) [][]string {
	remaining := in

	groups := make([][]string, 0)

	for len(remaining) > 0 {
		//if len(remaining) < s.size && getCategory(remaining[0]) == getCategory(remaining[len(remaining)-1]) {
		//	groups = append(groups, remaining)
		//	remaining = []string{}
		//} else {
		// Find the break point
		var point int

		for point = 1; point < len(remaining); point++ {
			if getCategory(remaining[point-1]) != getCategory(remaining[point]) {
				break
			}
		}

		if point > s.size {
			for point = s.size; point > 0; point-- {
				if getPrefix(remaining[point-1], prefixSize) != getPrefix(remaining[point], prefixSize) {
					break
				}
			}
		}

		groups = append(groups, remaining[:point])
		remaining = remaining[point:]
		//}
	}

	return groups
}
