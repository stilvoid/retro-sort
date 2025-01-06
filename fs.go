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

// FindFiles returns a list of filepaths to files with dir
// where the base filename matches the glob in pattern
func FindFiles(dir, pattern string) ([]string, error) {
	g, err := glob.Compile(strings.ToLower(pattern))
	if err != nil {
		return []string{}, err
	}

	files := make([]string, 0)

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if pattern != "*" && !g.Match(strings.ToLower(filepath.Base(path))) {
				return err
			}

			files = append(files, path)
		}

		return err
	})

	return files, err
}

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

// CopyFiles files from the path in the keys of fileMap to the path in the values.
// If upperCase is true, directory names will be upper-cased.
// If quiet is false, CopyFiles will output progress information to stdout.
func CopyFiles(root string, fileMap map[string]string, upperCase, quiet bool) error {
	counter := 0
	var div int = len(fileMap) / 20
	if div == 0 {
		div = 1
	}

	for src, dst := range fileMap {
		dir, fn := filepath.Split(dst)

		if upperCase {
			dir = strings.ToUpper(dir)
		}

		dir = filepath.Join(root, dir)

		os.MkdirAll(dir, 0750)

		dst = filepath.Join(dir, fn)

		if err := copyFile(src, dst); err != nil {
			return err
		}

		counter++
		if !quiet && counter%div == 0 {
			fmt.Print(".")
		}
	}

	return nil
}
