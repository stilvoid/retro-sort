package retrosort

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEntryPrefixLong(t *testing.T) {
	e := newEntry("a file name.foo", []string{"a file name.foo"})

	if d := cmp.Diff("a_file_name_foo", e.prefix(15)); d != "" {
		t.Error(d)
	}

	if d := cmp.Diff("a_file_name_foo", e.prefix(16)); d != "" {
		t.Error(d)
	}
}

func TestEntryPrefixSingle(t *testing.T) {
	es := []entry{
		newEntry("123", []string{}),
		newEntry(".-#", []string{}),
		newEntry(" _%3", []string{}),
	}

	for _, e := range es {
		if d := cmp.Diff("#", e.prefix(1)); d != "" {
			t.Error(e, d)
		}
	}
}

func TestEntryPrefixVarious(t *testing.T) {
	e := newEntry("my great(-)file.yes", []string{})

	testCases := []struct {
		expected string
		size     int
	}{
		{"m", 1},
		{"my_great", 8},
		{"my_great_file", 13},
	}

	for _, c := range testCases {

		if d := cmp.Diff(c.expected, e.prefix(c.size)); d != "" {
			t.Error(d)
		}
	}
}

func TestEntryFileMapSingle(t *testing.T) {
	e := newEntry("a", []string{"/path/to/a.file"})

	expected := map[string]string{
		"/path/to/a.file": "a.file",
	}

	actual := e.fileMap()

	if d := cmp.Diff(expected, actual); d != "" {
		t.Error(d)
	}
}

func TestEntryFileMapMultiple(t *testing.T) {
	e := newEntry("a", []string{"/path/to/a.file", "/other/path/for/a [test].file"})

	expected := map[string]string{
		"/path/to/a.file":               "a/a.file",
		"/other/path/for/a [test].file": "a/a [test].file",
	}

	actual := e.fileMap()

	if d := cmp.Diff(expected, actual); d != "" {
		t.Error(d)
	}
}
