package main

import (
	"fmt"

	"math"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/exp/maps"
)

/*
   Arguments:
       retro-sort ([source] ([dest])) (--move)

       [source] is a glob that matches files you wish to sort
       [dest] is an optional directory to place the sorted files into
       --move means that files will be moved rather than copied

       retro-sort will scan the files in [source]
       (or all files in the current directory if [source] is not supplied)
       new folders will then be created with a maximum of 100 files in each
       The short unique sequence of letters of the first file within it
       will be the folder name.

       For example:
           A/
               Aardvark
               ... 98 other files
               Azerty
           Azt/
               Aztec
               ... 98 other files
               Camel
           Cat/
               Cat
               ... 98 other files
               Zebra
*/

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s ([src] ([dst])) (--move)\n", os.Args[0])
	os.Exit(1)
}

func main() {
	var src, dst string
	move := false

	args := os.Args[1:]

	for _, arg := range args {
		if arg == "--move" {
			move = true
		} else {
			if src == "" {
				src = arg
			} else if dst == "" {
				dst = arg
			} else {
				usage()
			}
		}
	}

	if src == "" {
		src = "./*"
	}

	if dst == "" {
		dst = "./"
	}

	fmt.Printf("%s %s %s --move %v\n", os.Args[0], src, dst, move)

	ms, err := filepath.Glob(src)
	if err != nil {
		panic(err)
	}

	// Make groups, bro
	groups := makeGroups(ms)

	for group, fs := range groups {
		dstDir := filepath.Join(dst, group)
		os.MkdirAll(dstDir, 0750)
		for _, srcFile := range fs {
			dstFile := filepath.Join(dstDir, filepath.Base(srcFile))
			os.Link(srcFile, dstFile)
		}
	}

	fmt.Println("done")
}

func makeGroups(in []string) map[string][]string {
	slices.Sort(in)

	groups := make(map[string][]string)
	last := ""

	name := findPrefix(last, filepath.Base(in[0]))
	groups[name] = make([]string, 0)

	for _, m := range in {
		groups[name] = append(groups[name], m)

		// TODO: Make this more clever
		if len(groups[name]) > 95 {
			newName := findPrefix(last, filepath.Base(m))

			if len(newName) < 5 {
				name = newName
				last = filepath.Base(m)
			}
		}
	}

	// Recurse?
	if len(groups) > 100 {
		keys := maps.Keys(groups)

		topGroups := makeGroups(keys)

		newGroups := make(map[string][]string)

		for tg, gs := range topGroups {
			for _, g := range gs {
				newGroups[filepath.Join(tg, g)] = groups[g]
			}
		}

		groups = newGroups
	}

	return groups
}

func findPrefix(prev, next string) string {
	minLen := int(math.Min(float64(len(prev)), float64(len(next))))

	if minLen == 0 {
		return next[:1]
	}

	for i := 0; i < int(math.Min(float64(len(prev)), float64(len(next)))); i++ {
		a := prev[:i]
		b := next[:i]

		if a != b {
			return b
		}
	}

	return next
}
