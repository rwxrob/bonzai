package bon

import (
	"text/template"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"
	"github.com/rwxrob/bonzai/pkg/core/opts"
	"github.com/rwxrob/bonzai/pkg/core/run"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

func init() {
	run.AllowPanic = false
}

var Cmd = &bonzai.Cmd{
	Name:  `bon`,
	Short: `manage bonzai composite command trees`,
	Comp:  comp.Opts,
	Opts:  opts.WeekDaysAbbr,
	Vers:  `v0.0.1`,
	Cmds:  []*bonzai.Cmd{barCmd, fooCmd},
	Def:   fooCmd,
	Long:  ``,

	//Comp:    comp.Cmds,
	//Comp:    comp.ThreeLetterEngWeekday,
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
	Name:    `foo`,
	Alias:   `f|something`,
	Opts:    `one|two|three`,
	Comp:    comp.Opts,
	FuncMap: template.FuncMap{},
	Call: func(x *bonzai.Cmd, args ...string) error {
		x.FuncMap[`args`] = func() string { return to.Human(args) }
		x.Println(`hello from {{pre .Name}} with args {{args}} in {{exepath}}`)
		return nil
	},
}

var barCmd = &bonzai.Cmd{
	Name:  `bar`,
	Alias: `whatever|b`,
	Comp:  comp.FileDirCmdsOpts,
	Cmds:  []*bonzai.Cmd{otherCmd},
	Vars: map[string]string{
		`some`: `thing`,
	},
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}
