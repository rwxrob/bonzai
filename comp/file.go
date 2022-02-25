// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"github.com/rwxrob/bonzai/filter"
	"github.com/rwxrob/bonzai/util"
)

// File returns all file names for the directory and file prefix
// passed. If nothing is passed assumes the current working directory.
func File(x Command, args ...string) []string {
	match := ""
	dir := "."

	if len(args) > 0 {
		match = args[0]
	}

	list := []string{}
	list = append(list, util.Files(dir)...)
	list = filter.HasPrefix(list, match)
	return list

}
