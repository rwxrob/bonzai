// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package vars contains a [bonzai.Cmd] branch that can be grafted (imported) into any other command but is usually reserved for the root level or even as a standalone `var` application since no matter where it is grafted, it always deals with the same [bonzai.Vars] driver and only one [bonzai.VarsDriver] may exist in any executable (given its bonzai package scope). Subcommands do, however, recognize their place in the command tree.
*/
package vars

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/term"
)

var Cmd = &bonzai.Cmd{
	Name:  `var`,
	Vers:  `v0.8.0`,
	Short: `persistent variables in {{ exestatedir "vars.properties"}}`,
	Cmds:  []*bonzai.Cmd{getCmd},
}

var getCmd = &bonzai.Cmd{
	Name:    `get`,
	Short:   `print a variable`,
	NumArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		val, err := x.Caller.Caller.Get(args[0])
		if err != nil {
			return err
		}
		term.Print(val)
		return nil
	},
}

/*

var setCmd = &Z.Cmd{
	Name:        `set`,
	Summary:     `safely sets (persists) a cached variable`,
	Usage:       `(help|<name>) [<args>...]`,
	Description: setDoc,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,

	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) > 1 {
			val := strings.Join(args[1:], " ")
			if err := x.Caller.Caller.Set(args[0], val); err != nil {
				return err
			}
		}
		return getCmd.Call(x, args[0])
	},
}

//go:embed set.md
var setDoc string

var fileCmd = &Z.Cmd{
	Name:     `file`,
	Alias:  []string{"f"},
	Summary:  `outputs full path to the cached vars file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		term.Print(vars.Path())
		return nil
	},
}

//go:embed init.md
var initDoc string

var initCmd = &Z.Cmd{
	Name:        `init`,
	Alias:     []string{"i"},
	Summary:     `(re)initializes current variable cache`,
	Commands:    []*Z.Cmd{help.Cmd},
	UseVars:     true, // but fulfills at init() above
	Description: initDoc,
	Call: func(x *Z.Cmd, _ ...string) error {
		if term.IsInteractive() {
			r := term.Prompt(`Really initialize %v? (y/N) `, vars.DirPath())
			if r != "y" {
				return nil
			}
		}
		return Z.Vars.Init()
	},
}

//go:embed data.md
var dataDoc string

var dataCmd = &Z.Cmd{
	Name:        `data`,
	Alias:     []string{"d"},
	Summary:     `outputs contents of the cached variables file`,
	Description: dataDoc,
	Commands:    []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Print(vars.Data())
		return nil
	},
}

//go:embed edit.md
var editDoc string

var editCmd = &Z.Cmd{
	Name:        `edit`,
	Summary:     `edit variables file ({{execachedir "vars"}}) `,
	Description: editDoc,
	Alias:     []string{"e"},
	Commands:    []*Z.Cmd{help.Cmd},
	Call:        func(x *Z.Cmd, _ ...string) error { return vars.Edit() },
}

var deleteCmd = &Z.Cmd{
	Name:        `delete`,
	Alias:     []string{`d`, `del`, `unset`},
	Summary:     `delete variable(s) from cache`,
	Usage:       `(help|<name>...)`,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,
	Description: ` The {{aka}} command deletes the specified variable from cache.`,

	Call: func(x *Z.Cmd, args ...string) error {
		for _, i := range args {
			vars.Del(x.Caller.Caller.Path(i))
		}
		return nil
	},
}
*/
