package retrosort

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGroupTosecNoMoveSingle(t *testing.T) {
	fileMap := map[string]string{
		"a": "/some/path/a game (1985)(madeup games)[a].dsk",
		"b": "/other/path/different game (1972) (unknown software).dsk",
		"c": "/somewhere/amazing thing (demo) (1995)(giant flop).dsk",
	}

	// Expecting no change
	if d := cmp.Diff(fileMap, groupTosec(fileMap)); d != "" {
		t.Error(d)
	}
}

func TestGroupTosec(t *testing.T) {
	input := map[string]string{
		"a": "/some/path/a game (1985)(madeup games).dsk",
		"b": "/some/path/a game (demo) (1985)(madeup games)[a].dsk",
		"c": "/some/path/a game ( 1985) (madeup games)[a2][cr Quartex].dsk",
		"d": "/some/path/different game (1972) (unknown software).dsk",
		"e": "/somewhere/amazing thing (demo) (1995)(giant flop).dsk",
		"f": "/somewhere/amazing thing (1995) (giant flop) (retail) [b].dsk",
	}

	expected := map[string]string{
		"a": "/some/path/a game (1985)(madeup games)/a game (1985)(madeup games).dsk",
		"b": "/some/path/a game (1985)(madeup games)/a game (demo) (1985)(madeup games)[a].dsk",
		"c": "/some/path/a game (1985)(madeup games)/a game ( 1985) (madeup games)[a2][cr Quartex].dsk",
		"d": "/some/path/different game (1972) (unknown software).dsk",
		"e": "/somewhere/amazing thing (1995)(giant flop)/amazing thing (demo) (1995)(giant flop).dsk",
		"f": "/somewhere/amazing thing (1995)(giant flop)/amazing thing (1995) (giant flop) (retail) [b].dsk",
	}

	// Expecting no change
	if d := cmp.Diff(expected, groupTosec(input)); d != "" {
		t.Error(d)
	}
}

func TestTosecTitle(t *testing.T) {
	testCases := map[string]string{
		"/some/path/a game (1985)(madeup games)[a].dsk":           "a game (1985)(madeup games)",
		"/some/path/different game (1972) (unknown software).dsk": "different game (1972)(unknown software)",
		"/somewhere/amazing thing (demo) (1995)(giant flop).dsk":  "amazing thing (1995)(giant flop)",
		"/other/not a tosec title (potato) (lemons).dsk":          "",
	}

	for testCase, expected := range testCases {
		if d := cmp.Diff(expected, tosecTitle(testCase)); d != "" {
			t.Error(d)
		}
	}
}
