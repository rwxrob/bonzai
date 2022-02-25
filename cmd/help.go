// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"log"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

// Help provides help documentation for the caller allowing the specific
// section of help wanted to be passed as a tab-completable parameter.
var Help = &bonzai.Cmd{
	Name: `help`,
	Params: []string{
		"name", "title", "summary", "params", "commands", "description",
		"examples", "legal", "copyright", "license", "version",
	},
	Completer: comp.Help,
	Call: func(caller *bonzai.Cmd, args ...string) error {
		section := "all"
		if len(args) > 0 {
			section = args[0]
		}
		log.Printf("would show help about %v %v", caller.Name, section)
		return nil
	},
}
