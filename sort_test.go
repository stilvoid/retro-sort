package retrosort

import (
	_ "embed"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed data/tosec-spectrum.txt
var tosecSpectrumString string

//go:embed data/tosec-spectrum-sorted-100.txt
var tosecSpectrumSorted100String string

//go:embed data/tosec-spectrum-sorted-255.txt
var tosecSpectrumSorted255String string

//go:embed data/tosec-spectrum-sorted-1000.txt
var tosecSpectrumSorted1000String string

var tosecSpectrum []string
var tosecSpectrumSorted100 map[string]string
var tosecSpectrumSorted255 map[string]string
var tosecSpectrumSorted1000 map[string]string

func init() {
	tosecSpectrum = strings.Split(strings.TrimSpace(tosecSpectrumString), "\n")

	tosecSpectrumSorted100 = make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(tosecSpectrumSorted100String), "\n") {
		parts := strings.Split(line, "\t->\t")
		tosecSpectrumSorted100[parts[0]] = parts[1]
	}

	tosecSpectrumSorted255 = make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(tosecSpectrumSorted255String), "\n") {
		parts := strings.Split(line, "\t->\t")
		tosecSpectrumSorted255[parts[0]] = parts[1]
	}

	tosecSpectrumSorted1000 = make(map[string]string)
	for _, line := range strings.Split(strings.TrimSpace(tosecSpectrumSorted1000String), "\n") {
		parts := strings.Split(line, "\t->\t")
		tosecSpectrumSorted1000[parts[0]] = parts[1]
	}
}

func testSort(t *testing.T, files []string, size int) {
	output := Sort(files, size)

	dedupOutputFiles := make(map[string]bool)
	for _, dst := range output {
		dedupOutputFiles[dst] = true
	}

	if len(files) != len(dedupOutputFiles) {
		t.Errorf("Expected %d entries, got %d of which %d unique", len(files), len(output), len(dedupOutputFiles))
	}

	counts := make(map[string]map[string]bool)

	for _, dst := range output {
		for {
			dir := filepath.Dir(dst)

			if _, ok := counts[dir]; !ok {
				counts[dir] = make(map[string]bool)
			}

			counts[dir][dst] = true

			if dir == "." {
				break
			}

			dst = dir
		}
	}

	for dir, count := range counts {
		if len(count) > size {
			t.Errorf("Too many entries in '%s': %d > %d", dir, len(count), size)
		}
	}
}

func TestSort100(t *testing.T) {
	testSort(t, tosecSpectrum, 100)

	if d := cmp.Diff(tosecSpectrumSorted100, Sort(tosecSpectrum, 100)); d != "" {
		t.Error(d)
	}
}

func TestSort255(t *testing.T) {
	testSort(t, tosecSpectrum, 255)

	if d := cmp.Diff(tosecSpectrumSorted255, Sort(tosecSpectrum, 255)); d != "" {
		t.Error(d)
	}
}

func TestSort1000(t *testing.T) {
	testSort(t, tosecSpectrum, 1000)

	if d := cmp.Diff(tosecSpectrumSorted1000, Sort(tosecSpectrum, 1000)); d != "" {
		t.Error(d)
	}
}

func TestSort32(t *testing.T) {
	testSort(t, tosecSpectrum, 32)
}

func TestSort10(t *testing.T) {
	testSort(t, tosecSpectrum, 10)
}

func TestSortClashingFilenames(t *testing.T) {
	files := []string{
		"/a/abc.dsk",
		"/b/abc.dsk",
		"/a/b/abc.dsk",
		"/a/def.dsk",
		"/c/b/def.dsk",
	}

	testSort(t, files, len(files))
}
