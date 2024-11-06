package kimono

import (
	"os"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/vars"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{sanitizeCmd, workCmd, tagCmd},
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
	Alias:     `w`,
	Comp:      comp.CmdsOpts,
	Opts:      `on|off`,
	MinArgs:   1,
	MaxArgs:   1,
	MatchArgs: `on|off`,
	Call: func(x *bonzai.Cmd, args ...string) error {
		if args[0] == `on` {
			return WorkOn()
		}
		return WorkOff()
	},
}

const (
	TagPushEnv        = `KIMONO_PUSH_TAG`
	TagShortenEnv     = `KIMONO_SHORTEN_TAG`
	TagVersionPartEnv = `KIMONO_VERSION_PART`
)

var tagCmd = &bonzai.Cmd{
	Name:  `tag`,
	Alias: `t`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{tagListCmd, tagBumpCmd},
	Def:   tagListCmd,
}

var tagBumpCmd = &bonzai.Cmd{
	Name:    `bump`,
	Alias:   `b|up|i|inc`,
	Comp:    comp.CmdsOpts,
	Cmds:    []*bonzai.Cmd{vars.Cmd},
	Opts:    `major|minor|patch|m|M|p`,
	MaxArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		mustPush := getValueFor(
			`push-tags`, TagPushEnv, `false`,
		) == `true`
		var part VerPart
		part = optsToVerPart(
			getValueFor(
				`version-part`,
				TagVersionPartEnv,
				`patch`,
			),
		)
		if len(args) == 0 {
			val, err := bonzai.Vars.Get(`default-ver-part`)
			if err != nil {
				return err
			}
			part = optsToVerPart(val)
		} else {
			part = optsToVerPart(args[0])
		}
		return TagBump(part, mustPush)
	},
}

func optsToVerPart(x string) VerPart {
	switch x {
	case `major`, `M`:
		return Major
	case `minor`, `m`:
		return Minor
	case `patch`, `p`:
		return Patch
	}
	return Minor
}

var tagListCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `l`,
	Comp:  comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		shorten := getValueFor(
			`shorten-tags`,
			TagShortenEnv,
			`false`,
		) == `true`
		each.Println(TagList(shorten))
		return nil
	},
}

func getValueFor(key, envVar string, defaultValue string) string {
	if val, exists := os.LookupEnv(envVar); exists {
		return val
	}
	if val, err := bonzai.Vars.Get(key); err != nil {
		return val
	}
	return defaultValue
}
