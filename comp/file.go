// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"path/filepath"

	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/dir"
)

// File returns all file names for the directory and file prefix
// passed. If nothing is passed assumes the current working directory.
// This completer is roughly based on the behavior and appearance of the
// bash shell with forward slashes as separators and escaped spaces. By
// using this completer (instead of the shell) the command line
// interface remains consistent across all runtimes.
func File(x Command, args ...string) []string {

	// no completion if we already have one
	if len(args) > 1 {
		return []string{}
	}

	// catch edge cases
	if len(args) == 0 {
		if x != nil {
			return []string{x.GetName()} // will add tailing space
		}
		return dir.Entries("")
	}

	// no prefix of any kind, just a space following command
	if args[0] == "" {
		return dir.Entries("")
	}

	d, pre := filepath.Split(args[0])
	list := filt.BaseHasPrefix(dir.Entries(d), pre)

	for {

		if len(list) > 1 {
			return list
		}

		if len(list) == 1 && fs.IsDir(list[0]) {
			list = dir.Entries(list[0])
			continue
		}

		break
	}

	// just a single file left in the list
	return list
}
