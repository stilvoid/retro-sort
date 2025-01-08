package retrosort

import (
	"maps"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testSortOutput(expected map[string]string) string {
	files := slices.Collect(maps.Keys(expected))

	actual := Sort(files, len(files))

	return cmp.Diff(expected, actual)
}

func TestGroupTosec(t *testing.T) {
	TosecMode = true

	d := testSortOutput(map[string]string{
		"/some/path/a game (1985)(madeup games).dsk":                         "a game (1985)(madeup games)/a game (1985)(madeup games).dsk",
		"/some/path/a game (demo) (1985)(madeup games)[a].dsk":               "a game (1985)(madeup games)/a game (demo) (1985)(madeup games)[a].dsk",
		"/somewhere/amazing thing (demo) (1995)(giant flop).dsk":             "amazing thing (1995)(giant flop)/amazing thing (demo) (1995)(giant flop).dsk",
		"/other/path/different game (1972) (unknown software).dsk":           "different game (1972) (unknown software).dsk",
		"/some/other/path/a game ( 1985) (madeup games)[a2][cr Quartex].dsk": "a game (1985)(madeup games)/a game ( 1985) (madeup games)[a2][cr Quartex].dsk",
		"/elsewhere/not a tosec title (potato) (lemon).dsk":                  "not a tosec title (potato) (lemon).dsk",
		"/somewhere/amazing thing (1995) (giant flop) (retail) [b].dsk":      "amazing thing (1995)(giant flop)/amazing thing (1995) (giant flop) (retail) [b].dsk",
	})

	if d != "" {
		t.Error(d)
	}
}

func TestGroupNoTosec(t *testing.T) {
	TosecMode = false

	d := testSortOutput(map[string]string{
		"/some/path/a game (1985)(madeup games).dsk":                         "a game (1985)(madeup games).dsk",
		"/some/path/a game (demo) (1985)(madeup games)[a].dsk":               "a game (demo) (1985)(madeup games)[a].dsk",
		"/somewhere/amazing thing (demo) (1995)(giant flop).dsk":             "amazing thing (demo) (1995)(giant flop).dsk",
		"/other/path/different game (1972) (unknown software).dsk":           "different game (1972) (unknown software).dsk",
		"/some/other/path/a game ( 1985) (madeup games)[a2][cr Quartex].dsk": "a game ( 1985) (madeup games)[a2][cr Quartex].dsk",
		"/elsewhere/not a tosec title (potato) (lemon).dsk":                  "not a tosec title (potato) (lemon).dsk",
		"/somewhere/amazing thing (1995) (giant flop) (retail) [b].dsk":      "amazing thing (1995) (giant flop) (retail) [b].dsk",
	})

	if d != "" {
		t.Error(d)
	}
}

func TestDosecName(t *testing.T) {
	testCases := map[string]string{
		"/some/path/a game (1985)(madeup games)[a].dsk":           "a game (1985)(madeup games)",
		"/some/path/different game (1972) (unknown software).dsk": "different game (1972)(unknown software)",
		"/somewhere/amazing thing (demo) (1995)(giant flop).dsk":  "amazing thing (1995)(giant flop)",
		"/other/not a tosec title (potato) (lemons).dsk":          "not a tosec title (potato) (lemons).dsk",
	}

	for testCase, expected := range testCases {
		if d := cmp.Diff(expected, tosecName(testCase)); d != "" {
			t.Error(d)
		}
	}
}
