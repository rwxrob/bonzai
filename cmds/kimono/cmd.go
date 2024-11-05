package kimono

import (
	"os"

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
	Cmds: []*bonzai.Cmd{tagListCmd, tagBumpCmd},
	Def:  tagListCmd,
}

var tagBumpCmd = &bonzai.Cmd{
	Name:    `bump`,
	Comp:    comp.CmdsOpts,
	Cmds:    []*bonzai.Cmd{tagListCmd},
	Opts:    `major|minor|patch|m|M|p`,
	MaxArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		mustPush := len(os.Getenv("KIMONO_TAG_PUSH")) > 0
		var part VerPart
		if len(args) == 0 {
			part = Minor
		} else {
			switch args[0] {
			case "major", "M":
				part = Major
			case "minor", "m":
				part = Minor
			case "patch", "p":
				part = Patch
			}
		}
		TagBump(part, mustPush)
		return nil
	},
}

var tagListCmd = &bonzai.Cmd{
	Name: `list`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		each.Println(TagList())
		return nil
	},
}
