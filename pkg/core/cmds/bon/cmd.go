package bon

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"
	"github.com/rwxrob/bonzai/pkg/core/run"
)

func init() {
	run.AllowPanic = false
}

var Cmd = &bonzai.Cmd{
	Name:     `bon`,
	Summary:  `manage bonzai composite command trees`,
	Comp:     comp.CmdsParams,
	Version:  `v0.0.1`,
	Commands: []*bonzai.Cmd{barCmd, fooCmd},
}

var otherCmd = &bonzai.Cmd{
	Name:    `other`,
	Aliases: `o`,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var fooCmd = &bonzai.Cmd{
	Name:    `foo`,
	Aliases: `f|something`,
	Params:  `one|two|three`,
	Comp:    comp.Combine{comp.CmdsParams, comp.ThreeLetterEngWeekday},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var barCmd = &bonzai.Cmd{
	Name:     `bar`,
	Aliases:  `whatever|b`,
	Comp:     comp.FileDirCmdsParams,
	Commands: []*bonzai.Cmd{otherCmd},
	InitVars: map[string]string{
		`some`: `thing`,
	},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}
