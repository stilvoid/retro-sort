package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	retrosort "github.com/stilvoid/retro-sort"
)

var src, dst string
var size int
var pattern string
var upperCase bool
var printOnly bool

var rootCmd = &cobra.Command{
	Use:   "retro-sort [src] [dst]",
	Short: "retro-sort sorts your files into a folder structure suitable for use with retro hardware",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		src = args[0]
		dst = args[1]

		retrosort.Execute(src, dst, size, pattern, upperCase, printOnly)
	},
}

func init() {
	rootCmd.Flags().IntVarP(&size, "size", "s", 100, "Maximum number of directory entries")
	rootCmd.Flags().StringVarP(&pattern, "glob", "g", "*", "Only include files matching this glob")
	rootCmd.Flags().BoolVarP(&upperCase, "upper", "u", false, "Make upper-case directory names")
	rootCmd.Flags().BoolVarP(&printOnly, "dry-run", "n", false, "Dry run. Print the file names and exit")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
