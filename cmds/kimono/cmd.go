package kimono

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name: `kimono`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{sanitizeCmd, workCmd},
}

var sanitizeCmd = &bonzai.Cmd{
	Name: `sanitize`,
	Comp: comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return Sanitize()
	},
}

var workCmd = &bonzai.Cmd{
	Name: `work`,
	Comp: comp.CmdsOpts,
	Opts: `on|off`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing on|off")
		}
		if args[0] == "on" {
			return WorkOn()
		}
		return WorkOff()
	},
}

