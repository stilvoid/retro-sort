package retrosort

import (
	"fmt"
	"regexp"
	"slices"
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

type file struct {
	name     string
	sortName string
}

var sortNameRe = regexp.MustCompile(`[^a-z0-9]+`)

func newFile(fn string) file {
	sortName := filepath.Base(fn)
	sortName = strings.ToLower(sortName)
	sortName = sortNameRe.ReplaceAllString(sortName, "_")

	return file{
		name:     fn,
		sortName: sortName,
	}
}

func (f file) prefix(size int) string {
	if size >= len(f.sortName) {
		return f.sortName
	}

	if size == 1 {
		return getCategory(f.sortName[:size])
	}

	return f.sortName[:size]
}

type group struct {
	files      []file
	prefixSize int
	path       string
}

func newGroup(names []string) group {
	files := make([]file, len(names))
	for i, name := range names {
		files[i] = newFile(name)
	}

	slices.SortStableFunc(files, func(a, b file) int {
		return strings.Compare(a.sortName, b.sortName)
	})

	return group{
		files:      files,
		prefixSize: 0,
	}
}

func (g group) Len() int {
	return len(g.files)
}

func (g group) name() string {
	a := g.files[0].prefix(g.prefixSize)
	b := g.files[g.Len()-1].prefix(g.prefixSize)

	if g.prefixSize == 1 {
		a = getCategory(a)
		b = getCategory(b)
	}

	if a == b {
		return a
	}

	return fmt.Sprintf("%s-%s", a, b)
}

// split returns g split into sub-groups using the specific prefixSize.
// the second return value indicates whether the group was able to meet the size constraint
func (g group) split(prefixSize, size int) ([]group, bool) {
	success := true

	counts := make(map[string]int)
	prefixes := make([]string, 0)

	// Fail if any individual prefix is too big
	for _, file := range g.files {
		prefix := file.prefix(prefixSize)

		if counts[prefix] == 0 {
			prefixes = append(prefixes, prefix)
		}

		counts[prefix]++

		if counts[prefix] > size {
			success = false
		}
	}

	groups := make([]group, 0)
	cur := newGroup([]string{})
	cur.prefixSize = prefixSize
	cur.path = g.path
	pos := 0

	// Consolidate
	for i, prefix := range prefixes {
		// Copy
		for j := 0; j < counts[prefix]; j++ {
			cur.files = append(cur.files, g.files[pos])
			pos++
		}

		// Check if we split here
		if i == len(prefixes)-1 || cur.Len()+counts[prefixes[i+1]] > size {
			groups = append(groups, cur)
			cur = newGroup([]string{})
			cur.prefixSize = prefixSize
			cur.path = g.path
		}
	}

	if len(cur.files) > 0 {
		groups = append(groups, cur)
	}

	if len(groups) > size {
		success = false
	}

	return groups, success
}

func (g group) sort(size int) []group {
	groups, _ := g.split(g.prefixSize+1, size)

	out := make([]group, 0)

	for _, sub := range groups {
		if sub.Len() <= size {
			out = append(out, sub)
		} else {
			if len(groups) > 1 {
				sub.path = filepath.Join(sub.path, sub.name())
			}

			for _, part := range sub.sort(size) {
				out = append(out, part)
			}
		}
	}

	return out
}

func (g group) String() string {
	return fmt.Sprintf("%s / %s\t%d\t%d", g.path, g.name(), g.Len(), g.prefixSize)
}

func (g group) fileMap() map[string]string {
	out := make(map[string]string)

	for _, file := range g.files {
		out[file.name] = filepath.Join(g.path, g.name(), filepath.Base(file.name))
	}

	return out
}
