package util

import (
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/filt"
	"github.com/rwxrob/bonzai/maps"
)

// Files returns a slice of strings matching the names of the files
// within the given directory adding a slash to the end of any
// directories.
func Files(dir string) []string {
	dir = filepath.Clean(dir)
	files := []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	return maps.Prefix(maps.MarkDirs(entries), dir+string(os.PathSeparator))
}

//FilesWith takes the path of a directory and returns the name of the
//files with the matching prefix.
func FilesWith(dir, pre string) []string {
	return filt.HasPrefix(Files(dir), filepath.Join(dir, pre))
}
