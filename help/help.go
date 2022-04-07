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
			list = append(list, other...)
		}
	}

	if len(args) == 0 {
		return list
	}

	return filt.HasPrefix(list, args[0])
}

// printIfHave takes any thing with a named field and a value, converts
// everything to string values (with to.String) and prints it with
// Print after passing it through Format. If the value is an empty
// string logs that the thing has no field of that name.
func printIfHave(thing, name, value any) {
	if len(to.String(value)) == 0 {
		log.Printf("%v has no %v\n", to.String(thing), to.String(name))
		return
	}
	Z.PrintEmph(to.String(value))
	fmt.Println()
}

// ForTerminal converts the collective help documentation of the given
// command into curses terminal-friendly output and prints the help for
// the specified section. If the special "all" section is passed all
// sections will be displayed. The style is similar to UNIX manual pages
// and supports terminal formatting including color.. Documentation must
// be in BonzaiMark markup (see Z.Format). Emphasis is omitted if the
// terminal is not interactive (see Z.Emph).
func ForTerminal(x *Z.Cmd, section string) {

	switch section {

	case "name":
		printIfHave("command", "name", x.Name)

	case "title":
		printIfHave(x.Name, "title", x.Title())

	case "summary":
		printIfHave(x.Name, "summary", x.Summary)

	case "params":
		printIfHave(x.Name, "params", x.UsageParams())

	case "commands":
		printIfHave(x.Name, "commands", x.UsageCmdTitles())

	case "description", "desc":
		printIfHave(x.Name, "description", Z.Mark(x.Description))

	case "examples":
		log.Printf("examples are planned but not yet implemented")

	case "legal":
		printIfHave(x.Name, "legal", x.Legal())

	case "copyright":
		printIfHave(x.Name, "copyright", x.Copyright)

	case "license":
		printIfHave(x.Name, "license", x.License)

	case "version":
		printIfHave(x.Name, "version", x.Version)

	case "all":

		Z.PrintEmph("**NAME**\n")
		Z.PrintMark(x.Title() + "\n\n")

		// always print a synopsis so we can communicate with command
		// developers about invalid field combinations through ERRORs

		Z.PrintEmph("**SYNOPSIS**\n")

		switch {

		case x.Usage != "":
			Z.PrintMarkf("%v %v", x.Name, x.Usage)

		case x.Call == nil && x.Params != nil:
			Z.PrintMarkf(
				"{ERROR: Params without Call: %v}\n\n",
				strings.Join(x.Params, ", "),
			)

		case len(x.Commands) == 0 && x.Call == nil:
			Z.PrintMark("{ERROR: neither Call nor Commands defined}")

		case len(x.Commands) > 0 && x.Call == nil:
			Z.PrintMarkf("%v COMMAND", x.Name)

		case len(x.Commands) > 0 && x.Call != nil && len(x.Params) > 0:
			Z.PrintMarkf("%v (COMMAND|%v)", x.Name, x.UsageParams())

		case len(x.Commands) == 0 && x.Call != nil && len(x.Params) > 0:
			Z.PrintMarkf("%v %v", x.Name, x.UsageParams())

		case len(x.Commands) > 0 && x.Call != nil:
			Z.PrintMarkf(`%v [COMMAND]`, x.Name)

		case len(x.Commands) == 0 && x.Call != nil:
			Z.PrintMarkf(`%v`, x.Name)

		case x.Call != nil:
			Z.PrintMarkf(`%v`, x.Name)

		default:
			Z.PrintMark("{ERROR: unknown synopsis combination}")
		}

		if len(x.Commands) > 0 {
			Z.PrintEmph("**COMMANDS**\n")
			Z.PrintIndent(x.UsageCmdTitles())
			fmt.Println()
		}

		if len(x.Description) > 0 {
			Z.PrintEmph("**DESCRIPTION**\n")
			Z.PrintMark(x.Description)
		}

		legal := x.Legal()
		if len(legal) > 0 {
			Z.PrintEmph("**LEGAL**\n")
			Z.PrintIndent(legal)
			fmt.Println()
		}

		if len(x.Other) > 0 {
			for _, s := range x.Other {
				Z.PrintEmphf("**%v**\n", strings.ToUpper(s.Title))
				Z.PrintMark(s.Body)
			}
		}

	default:
		for _, s := range x.Other {
			if s.Title == section {
				Z.PrintMark(s.Body)
			}
		}
	}
}
