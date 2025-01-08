package retrosort

import (
	"fmt"
	"path/filepath"
	"regexp"
)

var TosecMode = false

// TOSEC name should begin with the game title,
// then an optional demo flag,
// then the publish date,
// then the publisher,
// then optional stuff up to the file extension
// We'll just capture up to the end of the publisher
var tosecRe = regexp.MustCompile(`^(.+?)\s*(?:\([^\)]+\)\s*)?\(\s*([0-9x]{4}(?:-[0-9x]{2}-[0-9x]{2})?)\s*\)\s*\(\s*([^\)]+\s*)\)`)

// tosecName checks if fn represents a valid tosec title
// if so, it stores in the registry and return a canonical title string
// If not, it return the original filename
func tosecName(fn string) string {
	fn = filepath.Base(fn)

	match := tosecRe.FindStringSubmatch(fn)

	if match == nil {
		return fn
	}

	return fmt.Sprintf("%s (%s)(%s)", match[1], match[2], match[3])
}
