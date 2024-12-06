package main

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
)

var cmd = &bonzai.Cmd{
	Name:  `help-test`,
	Alias: `h|ht`,
	Short: `just a help test`,
	Opts:  `some|-y|--yaml`,
	Cmds:  []*bonzai.Cmd{help.Cmd, fooCmd, hiddenCmd.AsHidden()},
	Comp:  comp.CmdsOpts,
	Def:   help.Cmd,
}

var fooCmd = &bonzai.Cmd{
	Name: `foo`,
	Cmds: []*bonzai.Cmd{underfooCmd},
}

var underfooCmd = &bonzai.Cmd{
	Name: `underfoo`,
	Do:   bonzai.Nothing,
}

var hiddenCmd = &bonzai.Cmd{
	Name: `imhidden`,
	Cmds: []*bonzai.Cmd{help.Cmd, barCmd},
}

var barCmd = &bonzai.Cmd{
	Name: `bar`,
	Do:   bonzai.Nothing,
}

func main() {
	cmd.Exec()
}
