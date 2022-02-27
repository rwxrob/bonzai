package maps

import (
	"io/fs"
	"path/filepath"

	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/mapf"
)

// Note to maintainers: This file contains simple maps that are
// implemented in the mapf package. Please keep complex maps in
// complex.go instead.

// MarkDirs will add an os.PathSeparator to the end of the name if the
// fs.DirEntry is a directory.
func MarkDirs(s []fs.DirEntry) []string { return fn.Map(s, mapf.MarkDirs) }

// Base extracts the filepath.Base of each path.
func Base(s []string) []string { return fn.Map(s, filepath.Base) }

// HashComment add the "# " prefix to each.
func HashComment(s []string) []string { return fn.Map(s, mapf.HashComment) }

// EscSpace replaces all spaces with backslashed spaces.
func EscSpace(s []string) []string { return fn.Map(s, mapf.EscSpace) }
