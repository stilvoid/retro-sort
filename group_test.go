package retrosort

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestGroup() group {
	return newGroup([]string{
		"/some/path/first",
		"/another/path/2nd.file",
		"/does/not/matter/and a third",
		"/somewhere/a 4th file",
	})
}

func TestGroupPrefixZero(t *testing.T) {
	g := getTestGroup()

	if g.Len() != 4 {
		t.Error("Length isn't 4")
	}

	if g.name() != "" {
		t.Error("Name isn't blank")
	}

	if d := cmp.Diff(": 4", g.String()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first":             "first",
		"/another/path/2nd.file":       "2nd.file",
		"/does/not/matter/and a third": "and a third",
		"/somewhere/a 4th file":        "a 4th file",
	}, g.fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroupPrefixOne(t *testing.T) {
	g := getTestGroup()

	g.prefixSize = 1

	if g.Len() != 4 {
		t.Error("Length is incorrect")
	}

	if g.name() != "#-f" {
		t.Error("Name is incorrect")
	}

	if d := cmp.Diff("#-f: 4", g.String()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first":             "#-f/first",
		"/another/path/2nd.file":       "#-f/2nd.file",
		"/does/not/matter/and a third": "#-f/and a third",
		"/somewhere/a 4th file":        "#-f/a 4th file",
	}, g.fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroupPrefixTwo(t *testing.T) {
	g := getTestGroup()

	g.prefixSize = 2

	if g.Len() != 4 {
		t.Error("Length is incorrect")
	}

	if g.name() != "2n-fi" {
		t.Error("Name is incorrect")
	}

	if d := cmp.Diff("2n-fi: 4", g.String()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first":             "2n-fi/first",
		"/another/path/2nd.file":       "2n-fi/2nd.file",
		"/does/not/matter/and a third": "2n-fi/and a third",
		"/somewhere/a 4th file":        "2n-fi/a 4th file",
	}, g.fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroupSortNone(t *testing.T) {
	g := getTestGroup()

	groups := g.sort(100)

	if len(groups) != 1 {
		t.Error("Should be just one group")
	}

	g = groups[0]

	if g.prefixSize != 0 {
		t.Errorf("Should be a 0 prefix, got %d", g.prefixSize)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first":             "first",
		"/another/path/2nd.file":       "2nd.file",
		"/does/not/matter/and a third": "and a third",
		"/somewhere/a 4th file":        "a 4th file",
	}, g.fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroupSortBasic(t *testing.T) {
	g := getTestGroup()

	groups := g.sort(2)

	if len(groups) != 3 {
		t.Errorf("Expected 3 groups, got %d", len(groups))
	}

	if d := cmp.Diff(map[string]string{
		"/another/path/2nd.file": "#/2nd.file",
	}, groups[0].fileMap()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/does/not/matter/and a third": "a/and a third",
		"/somewhere/a 4th file":        "a/a 4th file",
	}, groups[1].fileMap()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first": "f/first",
	}, groups[2].fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroupSortRecurse(t *testing.T) {
	g := getTestGroup()

	groups := g.sort(1)

	if len(groups) != 4 {
		t.Errorf("Expected 4 groups, got %d", len(groups))
	}

	if d := cmp.Diff(map[string]string{
		"/another/path/2nd.file": "#/2nd.file",
	}, groups[0].fileMap()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/somewhere/a 4th file": "a/a_/a 4th file",
	}, groups[1].fileMap()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/does/not/matter/and a third": "a/an/and a third",
	}, groups[2].fileMap()); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff(map[string]string{
		"/some/path/first": "f/first",
	}, groups[3].fileMap()); d != "" {
		t.Error(d)
	}
}

func TestGroups(t *testing.T) {
	g := getTestGroup()

	groups := g.sort(2)

	if d := cmp.Diff(map[string]string{
		"/another/path/2nd.file":       "#/2nd.file",
		"/somewhere/a 4th file":        "a/a 4th file",
		"/does/not/matter/and a third": "a/and a third",
		"/some/path/first":             "f/first",
	}, groups.fileMap()); d != "" {
		t.Error(d)
	}
}
