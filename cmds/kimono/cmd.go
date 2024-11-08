package kimono

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/vars"
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Short: `kimono is a tool for managing golang monorepos`,
	Vers:  `0.0.1`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{sanitizeCmd, workCmd, tagCmd, listCmd},
}

var sanitizeCmd = &bonzai.Cmd{
	Name: `sanitize`,
	Short: `sanitize will run ` + "`go get -u` and `go mod tidy`\n" +
		`on all go modules in the current git repo`,
	Comp: comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		root, err := futil.HereOrAbove(".git")
		if err != nil {
			return err
		}
		return Tidy(root)
	},
}

var workCmd = &bonzai.Cmd{
	Name:      `work`,
	Alias:     `w`,
	Short:     `work allows you to toggle go work files on or off`,
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
	TagDeleteRemote   = `KIMONO_DELETE_REMOTE_TAG`
)

var tagCmd = &bonzai.Cmd{
	Name:  `tag`,
	Alias: `t`,
	Short: `tag allows to bump version tags and list the tags for the go module`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{tagListCmd, tagBumpCmd},
	Def:   tagListCmd,
}

var tagBumpCmd = &bonzai.Cmd{
	Name:    `bump`,
	Alias:   `b|up|i|inc`,
	Short:   `bump bumps version tags subject to the given version part.`,
	Comp:    comp.CmdsOpts,
	Cmds:    []*bonzai.Cmd{vars.Cmd},
	Opts:    `major|minor|patch|m|M|p`,
	MaxArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		mustPush := vars.Fetch(TagPushEnv, `push-tags`, false)
		part := optsToVerPart(
			vars.Fetch(
				TagVersionPartEnv,
				`version-part`,
				`patch`,
			),
		)
		return TagBump(part, mustPush)
	},
}

var tagDeleteCmd = &bonzai.Cmd{
	Name:    `delete`,
	Alias:   `d|del|rm`,
	Short:   `delete the given tag from the go module`,
	Comp:    comp.Cmds,
	MinArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return TagDelete(
			args[0],
			vars.Fetch(
				TagDeleteRemote,
				`delete-remote-tag`,
				false,
			),
		)
	},
}

var listCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `l`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{tagListCmd, depListCmd, depDependentsCmd},
	Def:   depListCmd,
}

var depListCmd = &bonzai.Cmd{
	Name:  `dependencies`,
	Alias: `deps|dps`,
	Short: `list the dependencies of the go module`,
	Comp:  comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		deps, err := ListDependencies()
		if err != nil {
			return err
		}
		each.Println(deps)
		return nil
	},
}

var depDependentsCmd = &bonzai.Cmd{
	Name:  `dependents`,
	Alias: `depts|dpt`,
	Short: `list the dependents of the go module`,
	Comp:  comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		deps, err := ListDependents()
		if err != nil {
			return err
		}
		each.Println(deps)
		return nil
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
	Short: `list the tags for the go module`,
	Comp:  comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		shorten := vars.Fetch(
			TagShortenEnv,
			`shorten-tags`,
			false,
		)
		each.Println(TagList(shorten))
		return nil
	},
}
