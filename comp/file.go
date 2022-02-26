// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"strings"

	"github.com/rwxrob/bonzai/filt"
	"github.com/rwxrob/bonzai/maps"
	"github.com/rwxrob/bonzai/util"
)

// File returns all file names for the directory and file prefix
// passed. If nothing is passed assumes the current working directory.
func File(x Command, args ...string) []string {
	match := ""
	dir := ""

	if len(args) > 0 {
		// FIXME if there is an unescaped "/" at all, truncate the directory
		// and keep the rest to add on later for match
		if strings.HasSuffix(args[0], "/") {
			dir = args[0]
			match = ""
		} else {
			match = args[0]
		}
	}

	list := []string{}
	list = append(list, maps.Prefix(util.Files(dir), dir)...)
	list = filt.HasPrefix(list, match)
	list = maps.CleanPaths(list)
	return list

}
