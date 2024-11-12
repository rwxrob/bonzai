package main

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/rwxrob/bonzai/cmds/help"
)

var cmd = &bonzai.Cmd{
	Name:  `help-test`,
	Alias: `h|ht`,
	Short: `just a help test`,
	Opts:  `some|-y|--yaml`,
	Cmds:  []*bonzai.Cmd{help.Cmd, fooCmd},
	// Cmds:  []*bonzai.Cmd{fooCmd},
	Comp: comp.CmdsOpts,
	Def:  help.Cmd,
	// Def:   fooCmd,
}

var fooCmd = &bonzai.Cmd{
	Name: `foo`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		return nil
	},
}

func main() {
	cmd.Exec()
}
