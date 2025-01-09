package retrosort

import (
	"math"
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

	// Group folders that got too big
	if len(groups) > size {
		parts := int(math.Ceil(float64(len(groups)) / float64(size)))

		for i := 0; i < parts; i++ {
			start := i * len(groups) / parts
			end := (i+1)*len(groups)/parts - 1
			if end > len(groups)-1 {
				end = len(groups) - 1
			}

			if end-start > 0 {
				a := groups[start].entries[0].prefix(prefixSize)
				b := groups[end].entries[len(groups[end].entries)-1].prefix(prefixSize)
				prefix := a + "-" + b

				for j := start; j <= end; j++ {
					groups[j].path = filepath.Join(groups[j].path, prefix)
				}
			}
		}
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
