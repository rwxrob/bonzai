// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package Z (bonzai) provides a rooted node tree of commands and singular
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
package Z

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	config "github.com/rwxrob/config/pkg"
	"github.com/rwxrob/fn"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
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
// (ex: .exe) and is set at init() time (see ExePath).
var ExeName string

// Commands contains the commands to lookup when Run-ing an executable
// in "multicall" mode. Each value must begin with a *Cmd and the rest
// will be assumed to be string arguments to prepend. See Run.
var Commands map[string][]any

// UsageText is used for one-line UsageErrors. It's exported to allow
// for different languages.
var UsageText = `usage`

// DefaultUsageFunc is the default first-class function assigned to
// every Cmd that does not already define one. It is used to return
// a usage summary. Generally, it should only return a single line (even
// if that line is very long). Developers are encouraged to refer users
// to their chosen help command rather than producing usually long usage
// lines. If only the word "usage" needs to be changed (for a given
// language) consider UsageText instead. Note that most developers will
// simply change the Usage string when they do not want the default
// inferred usage string.
var DefaultUsageFunc = InferredUsage

// InferredUsage returns a single line of text summarizing only the
// Commands (less any Hidden commands), Params, and Aliases. If a Cmd
// is currently in an invalid state (Params without Call, no Call and no
// Commands) a string beginning with ERROR and wrapped in braces ({}) is
// returned instead. The string depends on the current language (see
// lang.go). Note that aliases does not include package Z.Aliases.
func InferredUsage(x *Cmd) string {

	if x.Call == nil && x.Commands == nil {
		// FIXME: replace with string var from lang.go
		return "{ERROR: neither Call nor Commands defined}"
	}

	if x.Call == nil && x.Params != nil {
		// FIXME: replace with string var from lang.go
		return "{ERROR: Params without Call: " + strings.Join(x.Params, ", ") + "}"
	}

	var params string
	if x.Params != nil {
		params = to.UsageGroup(x.Params)
		switch x.MinParm {
		case 0:
			params = params + "?"
		case 1:
			params = params
		default:
			params = fmt.Sprintf("%v{%v,}", params, x.MinParm)
		}
	}

	var names string
	if x.Commands != nil {
		var snames []string
		for _, x := range x.Commands {
			snames = append(snames, x.UsageNames())
		}
		if len(snames) > 0 {
			names = to.UsageGroup(snames)
		}
	}

	if params != "" && names != "" {
		return "(" + params + "|" + names + ")"
	}

	if params != "" {
		return params
	}

	return names
}

// Run infers the name of the command to run from the ExeName looked up
// in the Commands delegates accordingly, prepending any arguments
// provided in the CmdRun. Run produces an "unmapped multicall command"
// error if no match is found. This is an alternative to the simpler,
// direct Cmd.Run method from main where only one possible Cmd will ever
// be the root and allows for BusyBox (https://www.busybox.net)
// multicall binaries to be used for such things as very light-weight
// Linux distributions when used "FROM SCRATCH" in containers.
func Run() {
	if v, has := Commands[ExeName]; has {
		if len(v) < 1 {
			ExitError(fmt.Errorf("multicall command missing"))
		}
		cmd, iscmd := v[0].(*Cmd)
		if !iscmd {
			ExitError(fmt.Errorf("first value must be *Cmd"))
		}
		args := []string{cmd.Name}
		if len(v) > 1 {
			rest := os.Args[1:]
			for _, a := range v[1:] {
				s, isstring := a.(string)
				if !isstring {
					ExitError(fmt.Errorf("only string arguments allowed"))
				}
				args = append(args, s)
			}
			args = append(args, rest...)
		}
		os.Args = args
		cmd.Run()
		Exit()
	}
	ExitError(fmt.Errorf("unmapped multicall command: %v", ExeName))
}

// DefaultConfigurer is assigned to the Cmd.Root.Config during Cmd.Run.
// It is conventional for only Cmd.Root to have a Configurer defined.
var DefaultConfigurer = new(config.Configurer)

// ReplaceSelf replaces the current running executable at its current
// location with the successfully retrieved file at the specified URL or
// file path and duplicates the original files permissions. Only http
// and https URLs are currently supported. If not empty, a checksum file
// will be fetched from sumurl and used to validate the download before
// making the replacement. For security reasons, no backup copy of the
// replaced executable is kept. Also see AutoUpdate.
func ReplaceSelf(url, checksum string) error {
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

// AutoUpdate automatically updates the current process executable
// version by starting a goroutine that checks the current semantic
// version (version) against a remote one (versurl) and calling
// ReplaceSelf with the URL of the binary (binurl) and checksum (sumurl)
// if and update is needed. Note that the binary will often be named
// quite differently than the name of the currently running executable
// (ex: foo-mac -> foo, foo-linux.
//
// If a URL to a checksum file (sumurl) is not empty will optionally
// validate the new version downloaded against the checksum before
// replacing the currently running process executable with the new one.
// The format of the checksum file is the same as that output by any of
// the major checksum commands (sha512sum, for example) with one or more
// lines beginning with the checksum, whitespace, and then the name of
// the file. A single checksum file can be used for multiple versions
// but the most recent should always be at the top. When the update
// completes a message notifying of the update is logged to stderr.
//
// The function will fail silently under any of the following
// conditions:
//
//     * current user does not have write access to executable
//     * unable to establish a network connection
//     * checksum provided does not match
//
// Since AutoUpdate happens in the background no return value is
// provided. To enable logging of the update process (presumably for
// debugging) add the AutoUpdate flag to the Trace flags
// (trace.Flags|=trace.AutoUpdate). Also see Cmd.AutoUpdate.
func AutoUpdate(version, versurl, binurl, sumurl string) {
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

// ArgsOrIn takes an slice or nil as argument and if the slice has any
// length greater than 0 returns all the argument joined together with
// a single space between them. Otherwise, will read standard input
// until end of file reached (Cntl-D).
func ArgsOrIn(args []string) string {
	if args == nil || len(args) == 0 {
		return term.Read()
	}
	return strings.Join(args, " ")
}

// Aliases allows Bonzai tree developers to create aliases (similar to
// shell aliases) that are directly translated into arguments to the
// Bonzai tree executable by overriding the os.Args in a controlled way.
// The value of an alias is always a slice of strings that will replace
// the os.Args[2:]. A slice is used (instead of a string parsed with
// strings.Fields) to ensure that hard-coded arguments containing
// whitespace are properly handled.
var Aliases = make(map[string][]string)

// EscThese is set to the default UNIX shell characters which require
// escaping to be used safely on the terminal. It can be changed to suit
// the needs of different host shell environments.
var EscThese = " \r\t\n|&;()<>![]"

// Esc returns a shell-escaped version of the string s. The returned value
// is a string that can safely be used as one token in a shell command line.
func Esc(s string) string {
	var buf []rune
	for _, r := range s {
		for _, esc := range EscThese {
			if r == esc {
				buf = append(buf, '\\')
			}
		}
		buf = append(buf, r)
	}
	return string(buf)
}

// EscAll calls Esc on all passed strings.
func EscAll(args []string) []string { return fn.Map(args, Esc) }
