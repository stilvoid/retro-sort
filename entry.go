package retrosort

import (
	"regexp"
	"strings"

	"path/filepath"
)

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

// entry represents a single entry in the output tree
// name is the file name for a single source
// or the group name for multiple sources
// sources is a list of paths to copy from
type entry struct {
	name     string
	sortName string
	sources  []string
}

var sortNameRe = regexp.MustCompile(`[^a-z0-9]+`)

func newEntry(name string, sources []string) entry {
	sortName := strings.ToLower(name)
	sortName = sortNameRe.ReplaceAllString(sortName, "_")

	return entry{
		name:     name,
		sortName: sortName,
		sources:  sources,
	}
}

func (e entry) prefix(size int) string {
	if size == 1 {
		return getCategory(e.sortName[:size])
	}

	if size >= len(e.sortName) {
		return e.sortName
	}

	return e.sortName[:size]
}

func (e entry) fileMap() map[string]string {
	out := make(map[string]string)

	for _, source := range e.sources {
		path := filepath.Base(source)

		if len(e.sources) > 1 {
			path = filepath.Join(e.name, path)
		}

		out[source] = path
	}

	return out
}
