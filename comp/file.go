// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"path/filepath"

	"github.com/rwxrob/bonzai/check"
	"github.com/rwxrob/bonzai/filt"
	"github.com/rwxrob/bonzai/util"
)

// File returns all file names for the directory and file prefix
// passed. If nothing is passed assumes the current working directory.
// This completer is roughly based on the behavior and appearance of the
// bash shell with forward slashes as separators and escaped spaces. By
// using this completer (instead of the shell) the command line
// interface remains consistent across all runtimes.
func File(x Command, args ...string) []string {

	// no more completion if we already have one
	if len(args) > 1 {
		return []string{}
	}

	// FIXME NEED TO ESCAPE ALL SPACES!

	// catch edge case where no space follows the command
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// no prefix of any kind, just a space following command
	if args[0] == "" {
		return util.Files("")
	}

	dir, pre := filepath.Split(args[0])
	list := filt.BaseHasPrefix(util.Files(dir), pre)

	for {

		if len(list) > 1 {
			return list
		}

		if len(list) == 1 && check.IsDir(list[0]) {
			list = util.Files(list[0])
			continue
		}

		break
	}

	// just a single file left in the list
	return list
}
