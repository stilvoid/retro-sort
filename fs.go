package retrosort

import (
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
