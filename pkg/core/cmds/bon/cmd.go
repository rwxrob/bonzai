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
	Name:  `bon`,
	Short: `manage bonzai composite command trees`,
	//Comp:  comp.Cmds,
	Comp: comp.ThreeLetterEngWeekday,
	Vers: `v0.0.1`,
	Cmds: []*bonzai.Cmd{barCmd, fooCmd},
}

var otherCmd = &bonzai.Cmd{
	Name:  `other`,
	Alias: `o`,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var fooCmd = &bonzai.Cmd{
	Name:   `foo`,
	Alias:  `f|something`,
	Params: `one|two|three`,
	Comp:   comp.Combine{comp.CmdsParams, comp.ThreeLetterEngWeekday},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var barCmd = &bonzai.Cmd{
	Name:  `bar`,
	Alias: `whatever|b`,
	Comp:  comp.FileDirCmdsParams,
	Cmds:  []*bonzai.Cmd{otherCmd},
	Vars: map[string]string{
		`some`: `thing`,
	},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}
