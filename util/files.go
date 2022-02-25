package util

import (
	"os"

	"github.com/rwxrob/bonzai/filter"
)

// Files returns a slice of strings matching the names of the files
// within the given directory adding a slash to the end of any
// directories.
func Files(dir string) []string {
	if dir == "" {
		dir = "."
	}
	files := []string{}
	finfo, _ := os.ReadDir(dir)
	for _, f := range finfo {
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}
		files = append(files, name)
	}
	return files
}

func FilesWith(dir, pre string) []string {
	return filter.HasPrefix(Files(dir), pre)
}
