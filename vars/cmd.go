package vars

import (
	"fmt"
	"strings"

	"github.com/rwxrob/term"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"
)

var Cmd = &bonzai.Cmd{
	Name:  `var`,
	Alias: `vars`,
	Vers:  `v0.1.0`,
	Cmds: []*bonzai.Cmd{
		GetCmd, SetCmd, editCmd, initCmd, clearCmd,
		grepCmd, loadCmd, deleteCmd, dataCmd,
	},
	Comp: comp.Cmds,
	Def:  GetCmd,
}

// can be imported directly into other Cmds without problem
var GetCmd = &bonzai.Cmd{
	Name:    `get`,
	Comp:    Comp,
	NumArgs: 1,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		value, err := Data.Get(args[0])
		if err != nil {
			return err
		}
		term.Print(value)
		return nil
	},
}

// can be imported directly into other Cmds without problem
var SetCmd = &bonzai.Cmd{
	Name:    `set`,
	Comp:    Comp,
	MinArgs: 2,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		return Data.Set(args[0], strings.Join(args[1:], " "))
	},
}

var loadCmd = &bonzai.Cmd{
	Name:    `load`,
	Comp:    comp.FileDir,
	MaxArgs: 1,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		file := ""
		if len(args) > 0 {
			file = args[0]
		}
		data, err := run.FileOrIn(file)
		if err != nil {
			return err
		}
		return Data.Load(data)
	},
}

var grepCmd = &bonzai.Cmd{
	Name:    `grep`,
	NumArgs: 1,
	Cmds:    []*bonzai.Cmd{grepvCmd, grepkCmd},
	Comp:    comp.Cmds,
	Def:     grepkCmd,
}

var grepkCmd = &bonzai.Cmd{
	Name:    `keys`,
	Alias:   `k`,
	NumArgs: 1,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		value, err := Data.GrepK(args[0])
		if err != nil {
			return err
		}
		fmt.Print(value)
		return nil
	},
}

var grepvCmd = &bonzai.Cmd{
	Name:    `values`,
	Alias:   `v|val|vals`,
	NumArgs: 1,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		value, err := Data.GrepV(args[0])
		if err != nil {
			return err
		}
		fmt.Print(value)
		return nil
	},
}

var initCmd = &bonzai.Cmd{
	Name:    `init`,
	Alias:   `i`,
	NumArgs: 0,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		return Data.Init()
	},
}

var clearCmd = &bonzai.Cmd{
	Name:    `clear`,
	Alias:   `cl`,
	NumArgs: 0,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		return Data.Clear()
	},
}

var editCmd = &bonzai.Cmd{
	Name:  `edit`,
	Alias: `e|ed`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return Data.Edit()
	},
}

var deleteCmd = &bonzai.Cmd{
	Name:    `delete`,
	Alias:   `d|del`,
	NumArgs: 1,
	Call: func(_ *bonzai.Cmd, args ...string) error {
		return Data.Delete(args[0])
	},
}

var dataCmd = &bonzai.Cmd{
	Name:    `data`,
	NumArgs: 0,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		return Data.Print()
	},
}
