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
	Cmds:  []*bonzai.Cmd{help.Cmd, fooCmd},
	Comp:  comp.CmdsOpts,
	Def:   help.Cmd,
}

var fooCmd = &bonzai.Cmd{
	Name: `foo`,
	Do:   bonzai.Nothing,
}

func main() {
	cmd.Exec()
}
