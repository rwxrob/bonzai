// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
	"github.com/rwxrob/bonzai/futil"
)

type _FileDir struct{}

var FileDir = _FileDir{}

// Complete returns all file names for the directory and file prefix passed.
// If nothing is passed assumes the current working directory.  This
// [Completer] is roughly based on the behavior of the bash shell with
// forward slashes as separators and escaped spaces. By using this
// completer (instead of the shell) the command line interface remains
// consistent across all runtimes. Note that unlike bash completion no
// indication of the type of file is provided (i.e. dircolors support).
// The first argument, if passed, must be a [*bonzai.Cmd].
func (_FileDir) Complete(an any, args ...string) []string {

	if len(args) > 1 {
		return []string{}
	}

	if args == nil || (len(args) > 0 && args[0] == "") {
		return futil.DirEntriesAddSlashPath(".")
	}

	// catch edge cases
	if len(args) == 0 {
		if an != nil {
			x, is := an.(*bonzai.Cmd)
			if !is {
				log.Printf(`%v is a %T not *bonzai.Cmd`, an, an)
				return []string{}
			}
			return []string{x.Name} // will add tailing space
		}
		return futil.DirEntriesAddSlashPath("")
	}

	first := strings.TrimRight(args[0], string(filepath.Separator))
	d, pre := filepath.Split(first)

	if d == "" {
		list := filt.HasPrefix(futil.DirEntries("."), pre)
		if len(list) == 1 && futil.IsDir(list[0]) {
			return futil.DirEntriesAddSlashPath(list[0])
		}
		return futil.DirEntriesAddSlash(list)
	}

	for {
		list := filt.BaseHasPrefix(futil.DirEntries(d), pre)
		if len(list) > 1 {
			return futil.DirEntriesAddSlash(list)
		}
		if futil.IsDir(list[0]) {
			d = list[0]
			continue
		}
		return futil.DirEntriesAddSlash(list)
	}

	return []string{}
}
