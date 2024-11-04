package main

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/term"
)

var Cmd = &bonzai.Cmd{
	Name: `multicall`,
	Cmds: []*bonzai.Cmd{fooCmd, barCmd},
	Comp: comp.Cmds,
	Def:  fooCmd,

	/*
		Call: func(_ *bonzai.Cmd, _ ...string) error {
			term.Print(`multicall`)
			return nil
		},
	*/

}

var fooCmd = &bonzai.Cmd{
	Name: `foo`,
	Cmds: []*bonzai.Cmd{fooHelpCmd},
	Comp: comp.Cmds,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		term.Print(`fooCmd`)
		return nil
	},
}

var fooHelpCmd = &bonzai.Cmd{
	Name: `help`,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		term.Print(`fooHelp`)
		return nil
	},
}

var barCmd = &bonzai.Cmd{
	Name: `bar`,
	Cmds: []*bonzai.Cmd{barHelpCmd},
	Comp: comp.Cmds,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		term.Print(`barCmd`)
		return nil
	},
}

var barHelpCmd = &bonzai.Cmd{
	Name: `help`,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		term.Print(`barHelp`)
		return nil
	},
}
