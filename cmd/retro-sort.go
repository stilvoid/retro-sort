package main

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"github.com/stilvoid/retrosort"
)

var src, dst string
var size int
var pattern string
var upperCase bool
var printOnly bool
var quiet bool
var tosec bool

func init() {
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "Maximum number of directory entries")
	rootCmd.Flags().StringVarP(&pattern, "glob", "g", "*", "Only include files matching this glob")
	rootCmd.Flags().BoolVarP(&upperCase, "upper", "u", false, "Make upper-case directory names")
	rootCmd.Flags().BoolVarP(&printOnly, "dry-run", "n", false, "Dry run. Print the file names and exit")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Don't print anything, just do it")
	rootCmd.Flags().BoolVar(&tosec, "tosec", false, "Experimental: Detect TOSEC filenames and group related files")
}

var rootCmd = &cobra.Command{
	Use:   "retro-sort [src] [dst]",
	Short: "retro-sort sorts your files into a folder structure suitable for use with retro hardware",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src = args[0]
		dst = args[1]

		// Exit if dst exists
		if _, err := os.Stat(dst); !printOnly && !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "You must specify a destination folder that does not exist yet")
			os.Exit(1)
		}

		files, err := retrosort.FindFiles(src, pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding files: %s\n", err)
			os.Exit(1)
		}

		if !quiet && !printOnly {
			fmt.Printf("Found %d files\n", len(files))
		}

		retrosort.TosecMode = tosec

		fileMap := retrosort.Sort(files, size)

		if len(files) != len(fileMap) {
			fmt.Fprintln(os.Stderr, "Duplicate filenames detected. This operation will result in files missing.")
			// TODO: Print dupes
			os.Exit(1)
		}

		if printOnly {
			// Guarantee order
			sources := slices.Collect(maps.Keys(fileMap))
			slices.Sort(sources)
			for _, src := range sources {
				fmt.Printf("%s\t->\t%s\n", src, fileMap[src])
			}
		} else {
			if err := retrosort.CopyFiles(dst, fileMap, upperCase, false); err != nil {
				fmt.Fprintf(os.Stderr, "Error copying files: %s", err)
				os.Exit(1)
			}
			fmt.Println()
		}

		if !quiet && !printOnly {
			fmt.Println("done")
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
