package bon

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/run"
)

func init() {
	run.AllowPanic = true
}

var fooCmd = &bonzai.Cmd{
	Name:    `foo`,
	Aliases: `f|F|something`,
	Hidden:  `something`,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var barCmd = &bonzai.Cmd{
	Name:    `bar`,
	Aliases: `whatever|b`,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		x.Println(`hello from {{.Name}} in {{exepath}}`)
		return nil
	},
}

var Cmd = &bonzai.Cmd{
	Name: `bon`,
	//Aliases:  `bon|bonzaicli`,
	Version:  `v0.0.1`,
	Commands: []*bonzai.Cmd{barCmd, fooCmd},
}
