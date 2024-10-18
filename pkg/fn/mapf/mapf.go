// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package mapf contains nothing but map functions suitable for use with
the fn.Map generic function or the equivalent fn.A method. See the maps
package for functions that accept entire generic slices to be
transformed with mapf (or other) functions.

Note that any of the functions in this package can easily be added to a template.FuncMap for use in custom text|html/templates.
*/
package mapf

import (
	"io/fs"
	"strings"
)

// MarkDirs will add a slash (/) to the end of the name if the
// fs.DirEntry is a directory and return it as a string.
func MarkDirs(f fs.DirEntry) string {
	if f.IsDir() {
		return f.Name() + "/"
	}
	return f.Name()
}

// HashComment adds a "# " prefix.
func HashComment(line string) string { return "# " + line }

// EscSpace puts backslash in front of any space.
func EscSpace(s string) string { return strings.ReplaceAll(s, ` `, `\ `) }
