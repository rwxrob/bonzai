// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package help

import (
	"fmt"
	"log"
	"strings"

	"github.com/rwxrob/bonzai/comp"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/to"
)

// Cmd provides help documentation for the caller allowing the specific
// section of help wanted to be passed as a tab-completable parameter.
var Cmd = &Z.Cmd{
	Name:    `help`,
	Summary: `display help similar to man page format`,
	Params: []string{
		"name", "title", "summary", "params", "commands", "description",
		"examples", "legal", "copyright", "license", "version",
	},
	Completer: helpCompleter,

	Description: `
		The *help* command provide generic help documentation by looking at
		the different fields of the given command associated with it. To get
		specific help provide the command for which help is wanted before
		the help command. The exact section of help can also be specified as
		an parameter after the help command itself. `,

	Call: func(x *Z.Cmd, args ...string) error {
		// TODO detect if local web help is preferred over terminal
		if len(args) == 0 {
			args = append(args, "all")
		}
		ForTerminal(x.Caller, args[0])
		return nil
	},
}

func helpCompleter(x comp.Command, args ...string) []string {

	// not sure we've completed the command name itself yet
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, x.GetParams()...)

	// if the caller has other sections get those
	caller := x.GetCaller()
	if caller != nil {
		other := caller.GetOther()
		if other != nil {
			list = append(list, maps.Keys(other)...)
		}
	}

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}

// ForTerminal converts the collective help documentation of the given
// command into curses terminal-friendly output which minics UNIX man
// pages as much as possible. Documentation text is expected to be in
// standard BonzaiMark markup (see Z.Format).
//
// If the special "all" section is passed all sections will be
// displayed.
//
// If the "less" pager is detected and the terminal is interactive
// (stdout is to a terminal/tty) will call Z.SysExec to transfer control
// to it and send the output to it making it virtual indistinguishable
// from "man" page output.
//
// If the terminal is non-interactive, simple prints as
// 80-column-wrapped plain text.
func ForTerminal(x *Z.Cmd, section string) {
	switch section {
	case "name":
		Z.PrintIfHave("command", "name", x.Name)
	case "title":
		Z.PrintIfHave(x.Name, "title", x.Title())
	case "summary":
		Z.PrintIfHave(x.Name, "summary", x.Summary)
	case "params":
		Z.PrintIfHave(x.Name, "params", x.UsageParams())
	case "commands":
		Z.PrintIfHave(x.Name, "commands", x.UsageCmdTitles())
	case "description", "desc":
		Z.PrintIfHave(x.Name, "description", Z.Format(x.Description))
	case "examples":
		log.Printf("examples are planned but not yet implemented")
	case "legal":
		Z.PrintIfHave(x.Name, "legal", x.Legal())
	case "copyright":
		Z.PrintIfHave(x.Name, "copyright", x.Copyright)
	case "license":
		Z.PrintIfHave(x.Name, "license", x.License)
	case "version":
		Z.PrintIfHave(x.Name, "version", x.Version)
	case "all":
		fmt.Println(Z.Format("**NAME**"))
		fmt.Println(to.IndentWrapped(x.Title(), 7, 80) + "\n")
		// always print a synopsis so we can communicate with command
		// developers about invalid field combinations through ERRORs
		fmt.Println(Z.Format("**SYNOPSIS**"))
		switch {
		case x.Call == nil && x.Params != nil:
			// FIXME: replace with string var from lang.go
			fmt.Println(to.IndentWrapped("{ERROR: Params without Call: "+
				strings.Join(x.Params, ", ")+"}", 7, 80) + "\n")
		case len(x.Commands) == 0 && x.Call == nil:
			// FIXME: replace with string var from lang.go
			fmt.Println(to.IndentWrapped(
				"{ERROR: neither Call nor Commands defined}", 7, 80) + "\n")
		case len(x.Commands) > 0 && x.Call == nil:
			fmt.Println(
				to.IndentWrapped(Z.Emphasize("**"+x.Name+"** COMMAND"), 7, 80))
		case len(x.Commands) > 0 && x.Call != nil && len(x.Params) == 0:
			fmt.Println(
				to.IndentWrapped(Z.Emphasize("**"+x.Name+"** COMMAND"), 7, 80),
			)
		case len(x.Commands) > 0 && x.Call != nil && len(x.Params) > 0:
			fmt.Println(
				to.IndentWrapped(Z.Emphasize("**"+x.Name+
					"** (COMMAND|"+x.UsageParams()+")"), 7, 80),
			)
		case len(x.Commands) == 0 && x.Call != nil && len(x.Params) > 0:
			fmt.Println(Z.Emphasize("       **" + x.Name + "** " + x.UsageParams()))
		case len(x.Commands) == 0 && x.Call != nil:
			fmt.Println(to.IndentWrapped(Z.Emphasize("**"+x.Name+"**"), 7, 80))
		}
		fmt.Println()

		if len(x.Commands) > 0 {
			fmt.Println(Z.Format("**COMMANDS**"))
			fmt.Println(to.IndentWrapped(x.UsageCmdTitles(), 7, 80) + "\n")
		}
		if len(x.Description) > 0 {
			fmt.Println(Z.Format("**DESCRIPTION**"))
			body := to.Dedented(x.Description)
			fmt.Println(to.IndentWrapped(Z.Format(body), 7, 80) + "\n")
		}
		legal := x.Legal()
		if len(legal) > 0 {
			fmt.Println(Z.Format("**LEGAL**"))
			fmt.Println(to.IndentWrapped(legal, 7, 80) + "\n")
		}
		if len(x.Other) > 0 {
			for section, text := range x.Other {
				fmt.Println(Z.Format("**" + strings.ToUpper(section) + "**"))
				fmt.Println(to.IndentWrapped(text, 7, 80) + "\n")
			}
		}
	default:
		v, has := x.Other[section]
		if !has {
			Z.PrintIfHave(x.Name, section, "")
		}
		Z.PrintIfHave(x.Name, section, Z.Format(v))
	}
}
