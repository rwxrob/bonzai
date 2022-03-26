// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package bonzai provides a rooted node tree of commands and singular
parameters making tab completion a breeze and complicated applications
much easier to intuit without reading all the docs. Documentation is
embedded with each command removing the need for separate man pages and
such and can be viewed as text or a locally served web page.

Rooted Node Tree

Commands and parameters are linked to create a rooted node tree of the
following types of nodes:

    * Leaves with a method and optional parameters
		* Branches with leaves, other branches, and a optional method
		* Parameters, single words that are passed to a leaf command

*/
package bonzai

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/fs/file"
)

func init() {
	var err error
	// get the full path to current running process executable
	ExePath, err = os.Executable()
	if err != nil {
		log.Print(err)
		return
	}
	ExePath, err = filepath.EvalSymlinks(ExePath)
	if err != nil {
		log.Print(err)
	}
	ExeName = strings.TrimSuffix(
		filepath.Base(ExePath), filepath.Ext(ExePath))
}

// ExePath holds the full path to the current running process executable
// which is determined at init() time by calling os.Executable and
// passing it to path/filepath.EvalSymlinks to ensure it is the actual
// binary executable file. Errors are reported to stderr, but there
// should never be an error logged unless something is wrong with the Go
// runtime environment.
var ExePath string

// ExeName holds just the base name of the executable without any suffix
// (ex: .exe).
var ExeName string

// ReplaceSelf replaces the current running executable at its current
// location with the successfully retrieved file at the specified URL or
// file path and duplicates the original files permissions. Only http
// and https URLs are currently supported. If not empty, a checksum file
// will be fetched from sumurl and used to validate the download before
// making the replacement. For security reasons, no backup copy of the
// replaced executable is kept. Also see AutoUpdate.
func ReplaceSelf(url, sumurl string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}
	// TODO validate sum
	return file.Replace(exe, url)
}

// AutoUpdate automatically updates the current process executable (see
// Exe) by starting a goroutine that checks the current version against
// a remote one and calling ReplaceSelf if needed.
//
// If cur is an int assumes it is an isosec (see uniq.IsoSecond) and
// that newURL will return just a single line with an isosec in the body
// (usually a file named UPDATED).
//
// If cur is a string assumes it is a valid semantic version (beginning
// with a 'v') and will expect a single JSON string (don't forget the
// wrapping double-quotes) from newURL (usually in a file named VERSION)
// which will be compared using the Compare function from
// golang.org/x/mod/semver.
//
// If a URL to a checksum file (sum) is not empty will optionally
// validate the new version downloaded against the checksum before
// replace the currently running process executable with the new one.
// The format of the checksum file is the same as that output by any of
// the major checksum commands (sha512sum, for example) with one or more
// lines beginning with the checksum, whitespace, and then the name of
// the file. A single checksum file can be used for multiple versions
// but the most recent should always be at the top. When the update
// completes a message notifying of the update is logged to stderr.
//
// Since AutoUpdate happens in the background no return value is
// provided. This means that failed Internet connections and other
// common reasons blocking the update silently fail. To enable logging
// of the update process (presumably for debugging) add the AutoUpdate
// flag to the Trace flags (trace.Flags|=trace.AutoUpdate).
func AutoUpdate[T int | string](cur T, newURL, sum, exe string) {
	// TODO
}

// Method defines the main code to execute for a command (Cmd). By
// convention the parameter list should be named "args" if there are
// args expected and underscore (_) if not. Methods must never write
// error output to anything but standard error and should almost always
// use the log package to do so.
type Method func(caller *Cmd, args ...string) error

// ----------------------- errors, exit, debug -----------------------

// DoNotExit effectively disables Exit and ExitError allowing the
// program to continue running, usually for test evaluation.
var DoNotExit bool

// ExitOff sets DoNotExit to false.
func ExitOff() { DoNotExit = true }

// ExitOn sets DoNotExit to true.
func ExitOn() { DoNotExit = false }

// Exit calls os.Exit(0) unless DoNotExit has been set to true. Cmds
// should never call Exit themselves returning a nil error from their
// Methods instead.
func Exit() {
	if !DoNotExit {
		os.Exit(0)
	}
}

// ExitError prints err and exits with 1 return value unless DoNotExit
// has been set to true. Commands should usually never call ExitError
// themselves returning an error from their Method instead.
func ExitError(err ...interface{}) {
	switch e := err[0].(type) {
	case string:
		if len(e) > 1 {
			log.Printf(e+"\n", err[1:]...)
		} else {
			log.Println(e)
		}
	case error:
		out := fmt.Sprintf("%v", e)
		if len(out) > 0 {
			log.Println(out)
		}
	}
	if !DoNotExit {
		os.Exit(1)
	}
}

// ArgsFrom returns a list of field strings split on space with an extra
// trailing special space item appended if the line has any trailing
// spaces at all signifying a definite word boundary and not a potential
// prefix.
func ArgsFrom(line string) []string {
	args := []string{}
	if line == "" {
		return args
	}
	args = strings.Fields(line)
	if line[len(line)-1] == ' ' {
		args = append(args, "")
	}
	return args
}

// Files returns a slice of strings matching the names of the files
// within the given directory adding a slash to the end of any
// directories and escaping any spaces by adding backslash.  Note that
// this function assumes forward slash path separators since completion
// is only supported on operating systems where such is the case.
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
