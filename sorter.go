package retrosort

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	// Make groups, bro
	group := newGroup(files)

	groups := group.sort(s.size)

	counter := 0
	var div int = len(files) / 20

	for _, group := range groups {
		for in, out := range group.fileMap() {
			dir, fn := filepath.Split(out)

			if s.upperCase {
				dir = strings.ToUpper(dir)
			}

			dir = filepath.Join(s.dst, dir)

			dstFile := filepath.Join(dir, fn)

			if s.printOnly {
				fmt.Printf("%s -> %s\n", in, dstFile)
			} else {
				os.MkdirAll(dir, 0750)

				if err := copyFile(in, dstFile); err != nil {
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
