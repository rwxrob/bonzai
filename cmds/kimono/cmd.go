package kimono

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
)

var Cmd = &bonzai.Cmd{
	Name: `kimono`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{sanitizeCmd, workCmd, tagCmd},
}

var sanitizeCmd = &bonzai.Cmd{
	Name: `sanitize`,
	Comp: comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return Sanitize()
	},
}

var workCmd = &bonzai.Cmd{
	Name:      `work`,
	Comp:      comp.CmdsOpts,
	Opts:      `on|off`,
	MinArgs:   1,
	MaxArgs:   1,
	MatchArgs: `on|off`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		if args[0] == "on" {
			return WorkOn()
		}
		return WorkOff()
	},
}

var tagCmd = &bonzai.Cmd{
	Name: `tag`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{tagListCmd},
}

var tagListCmd = &bonzai.Cmd{
	Name: `list`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		each.Println(TagList())
		return nil
	},
}
