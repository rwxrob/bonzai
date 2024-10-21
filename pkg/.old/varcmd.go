package bonzai

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/bonzai/pkg/core/term"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var vars Map

func init() {
	dir, _ := os.UserCacheDir()
	vars = New()
	vars.Id = ExeName
	vars.Dir = dir
	vars.File = `vars`
	Z.Vars = vars
}

//go:embed text/vars_help.md
var helpDoc string

var Cmd = &Z.Cmd{
	Name:        `var`,
	Summary:     `cache variables in {{ execachedir "vars"}}`,
	Description: helpDoc,
	Version:     `v0.6.4`,
	Copyright:   `Copyright 2021 Robert S Muhlestein`,
	License:     `Apache-2.0`,
	Source:      `git@github.com:rwxrob/vars.git`,
	Issues:      `https://github.com/rwxrob/vars/issues`,
	Commands: []*Z.Cmd{
		getCmd, // default
		help.Cmd, initCmd, setCmd, fileCmd, dataCmd, editCmd, deleteCmd,
	},
}

//go:embed get.md
var getDoc string

var getCmd = &Z.Cmd{
	Name:        `get`,
	Summary:     `print a cached variable with a new line`,
	Commands:    []*Z.Cmd{help.Cmd},
	Description: getDoc,
	NumArgs:     1,

	Call: func(x *Z.Cmd, args ...string) error {
		val, err := x.Caller.Caller.Get(args[0])
		if err != nil {
			return err
		}
		term.Print(val)
		return nil
	},
}

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

//go:embed text/vars_set.md
var setDoc string

var fileCmd = &Z.Cmd{
	Name:     `file`,
	Aliases:  []string{"f"},
	Summary:  `outputs full path to the cached vars file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		term.Print(vars.Path())
		return nil
	},
}

//go:embed text/vars_init.md
var initDoc string

var initCmd = &Z.Cmd{
	Name:        `init`,
	Aliases:     []string{"i"},
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

//go:embed text/vars_data.md
var dataDoc string

var dataCmd = &Z.Cmd{
	Name:        `data`,
	Aliases:     []string{"d"},
	Summary:     `outputs contents of the cached variables file`,
	Description: dataDoc,
	Commands:    []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Print(vars.Data())
		return nil
	},
}

//go:embed text/vars_edit.md
var editDoc string

var editCmd = &Z.Cmd{
	Name:        `edit`,
	Summary:     `edit variables file ({{execachedir "vars"}}) `,
	Description: editDoc,
	Aliases:     []string{"e"},
	Commands:    []*Z.Cmd{help.Cmd},
	Call:        func(x *Z.Cmd, _ ...string) error { return vars.Edit() },
}

var deleteCmd = &Z.Cmd{
	Name:        `delete`,
	Aliases:     []string{`d`, `del`, `unset`},
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
