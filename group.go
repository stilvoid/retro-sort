package retrosort

import (
	"fmt"
	"slices"
	"strings"

	"path/filepath"
)

type group struct {
	entries    []entry
	prefixSize int
	path       string
}

func newGroup(sources []string) group {
	groupedSources := make(map[string][]string)
	for _, fn := range sources {
		name := filepath.Base(fn)

		if TosecMode {
			name = tosecName(name)
		}

		if _, ok := groupedSources[name]; !ok {
			groupedSources[name] = make([]string, 0)
		}

		groupedSources[name] = append(groupedSources[name], fn)
	}

	entries := make([]entry, 0)
	for name, sources := range groupedSources {
		entries = append(entries, newEntry(name, sources))
	}

	slices.SortStableFunc(entries, func(a, b entry) int {
		return strings.Compare(a.sortName, b.sortName)
	})

	return group{
		entries:    entries,
		prefixSize: 0,
	}
}

func (g group) Len() int {
	return len(g.entries)
}

func (g group) name() string {
	if g.prefixSize == 0 {
		return ""
	}

	a := g.entries[0].prefix(g.prefixSize)
	b := g.entries[g.Len()-1].prefix(g.prefixSize)

	if g.prefixSize == 1 {
		a = getCategory(a)
		b = getCategory(b)
	}

	if a == b {
		return a
	}

	return fmt.Sprintf("%s-%s", a, b)
}

func (g group) String() string {
	path := filepath.Join(g.path, g.name())

	return fmt.Sprintf("%s: %d", path, g.Len())
}

func (g group) fileMap() map[string]string {
	out := make(map[string]string)

	for _, entry := range g.entries {
		for src, dst := range entry.fileMap() {
			out[src] = filepath.Join(g.path, g.name(), dst)
		}
	}

	return out
}

type groups []group

func (gs groups) fileMap() map[string]string {
	out := make(map[string]string)

	dedup := make(map[string]bool)

	for _, g := range gs {
		for src, dst := range g.fileMap() {
			if dedup[dst] {
				dir, fn := filepath.Split(dst)
				ext := filepath.Ext(fn)
				for i := 2; ; i++ {
					dst = filepath.Join(dir, fmt.Sprintf("%s-%d%s", fn[:len(fn)-len(ext)], i, ext))
					if !dedup[dst] {
						break
					}
				}
			}

			out[src] = dst
			dedup[dst] = true
		}
	}

	return out
}
