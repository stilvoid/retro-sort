package retrosort

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"os"
	"path/filepath"

	"github.com/gobwas/glob"
)

func copyFile(srcFn, dstFn string) error {
	src, err := os.Open(srcFn)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFn)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func findFiles(dir, pattern string) ([]string, error) {
	g, err := glob.Compile(strings.ToLower(pattern))
	if err != nil {
		return []string{}, err
	}

	files := make([]string, 0)

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			if pattern != "*" && !g.Match(strings.ToLower(filepath.Base(path))) {
				return err
			}

			files = append(files, path)
		}

		return err
	})

	return files, err
}

func getPrefix(fn string, prefixSize int) string {
	fn = strings.ToLower(filepath.Base(fn))

	if len(fn) < prefixSize {
		return fn
	}

	return fn[:prefixSize]
}

func getCategory(fn string) string {
	c := strings.ToLower(filepath.Base(fn))[0]

	if c >= 'a' && c <= 'z' {
		return string(c)
	}

	return "#"
}

func groupName(group []string, prefixSize int) string {
	a := getPrefix(group[0], prefixSize)
	b := getPrefix(group[len(group)-1], prefixSize)

	if prefixSize == 1 {
		a = getCategory(a)
		b = getCategory(b)
	}

	if a == b {
		return a
	}

	return fmt.Sprintf("%s-%s", a, b)
}
