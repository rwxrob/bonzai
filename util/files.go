package util

import (
	"os"

	"github.com/rwxrob/fn/maps"
)

// Files returns a slice of strings matching the names of the files
// within the given directory adding a slash to the end of any
// directories and escaping any spaces by adding backslash. Note that
// this (and all functions of the bonzai package) assume forward slash
// path separators because no path argument should ever be passed to any
// bonzai command or high-level library that does not use forward slash
// paths. Commands should always use the comp.Files completer instead of
// host shell completion.
func Files(dir string) []string {
	if dir == "" {
		dir = "."
	}
	files := []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return files
	}
	names := maps.MarkDirs(entries)
	if dir == "." {
		return names
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return maps.EscSpace(maps.Prefix(names, dir))
}
