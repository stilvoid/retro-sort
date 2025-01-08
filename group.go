package retrosort

import (
	"fmt"
	"slices"
	"strings"

	"path/filepath"
)

// Sort converts a list of paths to files into a mapping from source paths
// to destination paths where no directory in the destinations
// contains any more than size files
func Sort(sources []string, size int) map[string]string {
	group := newGroup(sources)

	groups := group.sort(size)

	return groups.fileMap()
}

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

// split returns g split into sub-groups using the specific prefixSize.
// the second return value indicates whether the group was able to meet the size constraint
func (g group) split(prefixSize, size int) (groups, bool) {
	success := true

	counts := make(map[string]int)
	prefixes := make([]string, 0)

	// Fail if any individual prefix is too big
	for _, entry := range g.entries {
		prefix := entry.prefix(prefixSize)

		if counts[prefix] == 0 {
			prefixes = append(prefixes, prefix)
		}

		counts[prefix]++

		if counts[prefix] > size {
			success = false
		}
	}

	groups := make(groups, 0)
	cur := newGroup([]string{})
	cur.prefixSize = prefixSize
	cur.path = g.path
	pos := 0

	// Consolidate
	for i, prefix := range prefixes {
		// Copy
		for j := 0; j < counts[prefix]; j++ {
			cur.entries = append(cur.entries, g.entries[pos])
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

	if len(cur.entries) > 0 {
		groups = append(groups, cur)
	}

	if len(groups) > size {
		success = false
	}

	return groups, success
}

func (g group) sort(size int) groups {
	if g.Len() <= size {
		return groups{g}
	}

	gs, _ := g.split(g.prefixSize+1, size)

	out := make(groups, 0)

	for _, sub := range gs {
		if sub.Len() <= size {
			out = append(out, sub)
		} else {
			if len(gs) > 1 {
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

	for _, g := range gs {
		for src, dst := range g.fileMap() {
			out[src] = dst
		}
	}

	return out
}
