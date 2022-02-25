// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"github.com/rwxrob/bonzai/filter"
	"github.com/rwxrob/bonzai/util"
)

// Files returns all file names for the directory and file prefix
// passed. If nothing is passed assumes the current working directory.
func Files(x Command, args ...string) []string {
	match := ""

	if args != nil && len(args) > 0 {
		match = args[0]
	}

	list := []string{}

	// TODO if file separators detected, drill down to the proper
	// directory then check the leaf as the prefix of files

	list = append(list, util.Files(match)...)

	return filter.HasPrefix(list, match)
}
