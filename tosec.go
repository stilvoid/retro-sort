package retrosort

import (
	"fmt"
	"path/filepath"
	"regexp"
)

var TosecMode = false

var tosecFiles = make(map[string][]string)

func addTosec(tosecName, fileName string) {
	if _, ok := tosecFiles[tosecName]; !ok {
		tosecFiles[tosecName] = make([]string, 0)
	}

	tosecFiles[tosecName] = append(tosecFiles[tosecName], fileName)
}

// TOSEC name should begin with the game title,
// then an optional demo flag,
// then the publish date,
// then the publisher,
// then optional stuff up to the file extension
// We'll just capture up to the end of the publisher
var tosecRe = regexp.MustCompile(`^(.+?)\s*(?:\([^\)]+\)\s*)?\(\s*([0-9x]{4}(?:-[0-9x]{2}-[0-9x]{2})?)\s*\)\s*\(\s*([^\)]+\s*)\)`)

// doTosec checks if fn represents a valid tosec title
// if so, it stores in the registry and return a canonical title string
// If not, it return the original filename
func doTosec(fn string) string {
	match := tosecRe.FindStringSubmatch(filepath.Base(fn))

	if match == nil {
		return fn
	}

	name := fmt.Sprintf("%s (%s)(%s)", match[1], match[2], match[3])

	addTosec(name, fn)

	return name
}
