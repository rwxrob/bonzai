// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package help

import (
	"log"

	Z "github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/fn/filt"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/term"
)

var OBracketed = term.Under
var CBracketed = term.Reset

// Cmd provides help documentation for the caller allowing the specific
// section of help wanted to be passed as a tab-completable parameter.
var Cmd = &Z.Cmd{
	Name: `help`,
	Params: []string{
		"name", "title", "summary", "params", "commands", "description",
		"examples", "legal", "copyright", "license", "version",
	},
	Completer: helpCompleter,
	Call: func(caller *Z.Cmd, args ...string) error {
		section := "all"
		if len(args) > 0 {
			section = args[0]
		}
		log.Printf("would show help about %v %v", caller.Name, section)
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

// Render renders the incoming Help markup string for a curses terminal
// detecting if the terminal is interactive and if not rendering as
// plain text instead.
func Render(in string) string {
	out := ""
	for i := 0; i < len([]rune(in)); i++ {
		cur := in[i]
		switch cur {

		// <bracketed>
		case '<':
			out += OBracketed
			for {

			}
			out += CBracketed
		}
	}
	return string(out)
}

func parseBracketed() {
}
