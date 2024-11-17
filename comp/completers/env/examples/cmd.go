package main

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp/completers/env"
)

var XDGCompNames = env.NewCompNames("XDG_", false)
var NamesCmd = &bonzai.Cmd{
	Name:  "names",
	Alias: `n|-n|--names`,
	Comp:  XDGCompNames,
	Do: func(x *bonzai.Cmd, args ...string) error {
		fmt.Println("Names:", args)
		return nil
	},
}

var XDGCompVars = env.NewCompVars("XDG_", true)
var VarsCmd = &bonzai.Cmd{
	Name:  "vars",
	Alias: `v|-v|--vars`,
	Comp:  XDGCompVars,
	Do: func(x *bonzai.Cmd, args ...string) error {
		fmt.Println("Vars:", args)
		return nil
	},
}

var Cmd = &bonzai.Cmd{
	Name: `env`,
	Cmds: []*bonzai.Cmd{NamesCmd, VarsCmd, help.Cmd},
	Def:  help.Cmd,
}
