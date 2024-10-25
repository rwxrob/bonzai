// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package bonzai provides command tree composition and singular parameters
making tab completion a breeze and complicated applications much easier
to intuit without reading all the docs. Documentation is embedded with
each command removing the need for separate man pages and such and can
be viewed as text or a locally served web page.

# Rooted Node Tree

Commands and parameters are linked to create a rooted node tree of the
following types of nodes:

  - Leaves with a method and optional parameters
  - Branches with leaves, other branches, and a optional method
  - Parameters, single words that are passed to a leaf command
*/
package bonzai

import (
	"os"
	"strings"
)

// Vars may be optionally assigned any implementation of a bonzai.Vars
// but this is normally assigned at init() time by a bonzai.Vars driver
// module. Once assigned it should not be reassigned
// at any later time during runtime. Certain Bonzai branches and
// commands may require Z.Vars to be defined and those that do generally
// require the same implementation throughout all of runtime. Commands
// that require Z.Vars should set ReqVars to true. Other than the
// exceptional case of configuration commands that fulfill bonzai.Vars
// (and usually assign themselves to Z.Vars at init() time), commands
// must never require a specific implementation of bonzai.Vars.  This
// encourages command creators and Bonzai tree composers to centralize
// on a single form of caching without creating brittle dependencies and
// tight coupling. Caching persistence can be implemented in any number
// of ways without a problem and Bonzai trees simply need to be
// recompiled with a different bonzai.Vars implementation to switch
// everything that depends on cached variables.
// FIXME replace with embedded vars command
//var Vars bonzai.Vars

// UsageFunc is the default first-class function called if a Cmd that
// does not already define its own when usage information is needed (see
// bonzai.UsageFunc and Cmd.UsageError for more). By default,
// InferredUsage is assigned.
//
// It is used to return a usage summary. Generally, it should only
// return a single line (even if that line is very long).  Developers
// are encouraged to refer users to their chosen help command rather
// than producing usually long usage lines.
var UsageFunc = InferredUsage

// InferredUsage returns a single line of text summarizing only the
// Commands (less any Hide commands), Params, and Alias. If a Cmd
// is currently in an invalid state (Params without Call, no Call and no
// Commands) a string beginning with ERROR and wrapped in braces ({}) is
// returned instead. The string depends on the current language (see
// lang.go). Note that aliases does not include package bonzai.Alias.
func InferredUsage(x *Cmd) string {

	if x.Call == nil && x.Commands == nil {
		return "{ERROR: neither Call nor Commands defined}"
	}

	if x.Call == nil && x.Params != nil {
		return "{ERROR: Params without Call: " + strings.Join(x.Params, ", ") + "}"
	}

	params := UsageGroup(x.Params, x.MinParm, x.MaxParm)

	var names string
	if x.Commands != nil {
		var snames []string
		for _, x := range x.Commands {
			snames = append(snames, x.UsageNames())
		}
		if len(snames) > 0 {
			names = UsageGroup(snames, 1, 1)
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
// in the Commands and delegates accordingly, prepending any arguments
// provided. This allows for BusyBox-like (https://www.busybox.net)
// multicall binaries to be used for such things as very light-weight
// Linux distributions when used "FROM SCRATCH" in containers.  Although
// it shares the same name Z.Run should not confused with Cmd.Run. In
// general, Z.Run is for "multicall" and Cmd.Run is for "monoliths".
// Run may exit with the following errors:
//
// * MultiCallCmdNotFound
// * MultiCallCmdNotCmd
// * MultiCallCmdNotCmd
// * MultiCallCmdArgNotString
func Run() {
	if v, has := Commands[ExeName]; has {
		if len(v) < 1 {
			ExitError(MultiCallCmdNotFound{ExeName})
			return
		}
		cmd, iscmd := v[0].(*Cmd)
		if !iscmd {
			ExitError(MultiCallCmdNotCmd{ExeName, v[0]})
			return
		}
		args := []string{cmd.Name}
		if len(v) > 1 {
			rest := os.Args[1:]
			for _, a := range v[1:] {
				s, isstring := a.(string)
				if !isstring {
					ExitError(MultiCallCmdArgNotString{ExeName, a})
					return
				}
				args = append(args, s)
			}
			args = append(args, rest...)
		}
		os.Args = args
		cmd.Run()
		Exit()
		return
	}
	ExitError(MultiCallCmdNotFound{ExeName})
}
