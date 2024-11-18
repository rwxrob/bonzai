package main

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/comp/completers/env"
	"github.com/rwxrob/bonzai/fn/tr"
)

var XDG = env.Env{"XDG_", false}
var NamesCmd = &bonzai.Cmd{
	Name:  "names",
	Alias: `n|-n|--names`,
	Comp:  XDG,
	Do: func(x *bonzai.Cmd, args ...string) error {
		fmt.Println("Names:", args)
		return nil
	},
}

var XDGCaseInsensitive = comp.Combine{env.Env{"XDG_", true}, tr.Prefix{`$`}}
var VarsCmd = &bonzai.Cmd{
	Name:  "vars",
	Alias: `v|-v|--vars`,
	Comp:  XDGCaseInsensitive,
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
