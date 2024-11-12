package vars

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/term"
)

func cmdDefMap(x *bonzai.Cmd) (*Map, error) {
	file, err := futil.UserStateDir()
	if err != nil {
		return nil, err
	}
	file = filepath.Join(file, x.Root().Name, DefaultFileName)
	return NewMapFromInit(file)
}

var Cmd = &bonzai.Cmd{
	Name:  `var`,
	Alias: `vars`,
	Vers:  `v0.1.0`,
	Cmds: []*bonzai.Cmd{
		GetCmd, SetCmd, editCmd, initCmd, clearCmd,
		grepCmd, loadCmd, deleteCmd, dataCmd, help.Cmd,
	},
	Comp: comp.Cmds,
	Def:  GetCmd,
}

// can be imported directly into other Cmds without problem
var GetCmd = &bonzai.Cmd{
	Name:    `get`,
	Comp:    Comp,
	NumArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		value, err := m.Get(args[0])
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
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Set(args[0], strings.Join(args[1:], " "))
	},
}

var loadCmd = &bonzai.Cmd{
	Name:    `load`,
	Comp:    comp.FileDir,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		file := ""
		if len(args) > 0 {
			file = args[0]
		}
		data, err := run.FileOrIn(file)
		if err != nil {
			return err
		}
		return m.Load(data)
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
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		value, err := m.GrepK(args[0])
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
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		value, err := m.GrepV(args[0])
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
	Do: func(x *bonzai.Cmd, _ ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Init()
	},
}

var clearCmd = &bonzai.Cmd{
	Name:    `clear`,
	Alias:   `cl`,
	NumArgs: 0,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Clear()
	},
}

var editCmd = &bonzai.Cmd{
	Name:  `edit`,
	Alias: `e|ed`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Edit()
	},
}

var deleteCmd = &bonzai.Cmd{
	Name:    `delete`,
	Alias:   `d|del`,
	NumArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Delete(args[0])
	},
}

var dataCmd = &bonzai.Cmd{
	Name:    `data`,
	NumArgs: 0,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		m, err := cmdDefMap(x)
		if err != nil {
			return err
		}
		return m.Print()
	},
}
