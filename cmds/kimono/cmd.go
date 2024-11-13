package kimono

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/vars"
)

const (
	SanitizeAllEnv    = `KIMONO_SANITIZE_ALL`
	TagPushEnv        = `KIMONO_PUSH_TAG`
	TagShortenEnv     = `KIMONO_SHORTEN_TAG`
	TagVersionPartEnv = `KIMONO_VERSION_PART`
	TagDeleteRemote   = `KIMONO_DELETE_REMOTE_TAG`
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Short: `manage golang monorepos`,
	Vers:  `v0.2.1`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		sanitizeCmd,
		workCmd,
		tagCmd,
		depsCmd,
		vars.Cmd,
		help.Cmd,
	},
	Def: help.Cmd,
}

var sanitizeCmd = &bonzai.Cmd{
	Name:    `sanitize`,
	Alias:   `tidy|update`,
	Short:   "run `go get -u` and `go mod tidy` on all go modules in repo",
	Comp:    comp.Cmds,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if argIsOr(
			args,
			`all`,
			vars.Fetch(
				SanitizeAllEnv,
				`sanitize-all`,
				false,
			),
		) {
			root, err := futil.HereOrAbove(".git")
			if err != nil {
				return err
			}
			return Tidy(filepath.Dir(root))
		}
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		return Tidy(pwd)
	},
}

var workCmd = &bonzai.Cmd{
	Name:      `work`,
	Alias:     `w`,
	Short:     `toggle go work files on or off`,
	Comp:      comp.CmdsOpts,
	Opts:      `on|off`,
	MinArgs:   1,
	MaxArgs:   1,
	MatchArgs: `on|off`,
	Cmds:      []*bonzai.Cmd{workInitCmd},
	Do: func(x *bonzai.Cmd, args ...string) error {
		switch args[0] {
		case `on`:
			return WorkOn()
		case `off`:
			return WorkOff()
		default:
			return fmt.Errorf(
				"invalid argument: %s",
				args[0],
			)
		}
	},
}

var workInitCmd = &bonzai.Cmd{
	Name:    `init`,
	Alias:   `i`,
	Short:   `new go.work in module for dependencies in monorepo`,
	MinArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if args[0] == `all` {
			return WorkGenerate()
		}
		return WorkInit(args...)
	},
}

var tagCmd = &bonzai.Cmd{
	Name:  `tag`,
	Alias: `t`,
	Short: `manage or list tags for the go module`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		tagBumpCmd, tagListCmd, tagDeleteCmd, help.Cmd.AsHidden()},
	Def: tagListCmd,
}

var tagBumpCmd = &bonzai.Cmd{
	Name:    `bump`,
	Alias:   `b|up|i|inc`,
	Short:   `bumps semver tags. based on given version part.`,
	Comp:    comp.CmdsOpts,
	Cmds:    []*bonzai.Cmd{vars.Cmd.AsHidden()},
	Opts:    `major|minor|patch|M|m|p`,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		mustPush := vars.Fetch(TagPushEnv, `push-tags`, false)
		if len(args) == 0 {
			part := optsToVerPart(
				vars.Fetch(
					TagVersionPartEnv,
					`version-part`,
					`patch`,
				),
			)
			return TagBump(part, mustPush)
		}
		part := optsToVerPart(args[0])
		return TagBump(part, mustPush)
	},
}

var tagDeleteCmd = &bonzai.Cmd{
	Name:    `delete`,
	Alias:   `d|del|rm`,
	Short:   `delete the given semver tag for the go module`,
	Comp:    comp.Cmds,
	MinArgs: 1,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
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

var tagListCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `l`,
	Short: `list the tags for the go module`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		shorten := vars.Fetch(
			TagShortenEnv,
			`shorten-tags`,
			false,
		)
		each.Println(TagList(shorten))
		return nil
	},
}

var depsCmd = &bonzai.Cmd{
	Name:  `dependencies`,
	Alias: `deps|dep`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{dependsOnCmd, usedByCmd},
	Def:   dependsOnCmd,
}

var dependsOnCmd = &bonzai.Cmd{
	Name:  `depends-on`,
	Alias: `on|uses`,
	Short: `list the dependencies for the go module`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		deps, err := ListDependencies()
		if err != nil {
			return err
		}
		each.Println(deps)
		return nil
	},
}

var usedByCmd = &bonzai.Cmd{
	Name:  `depends-on-me`,
	Alias: `onme|usedby|me`,
	Short: `list the dependents of the go module`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		deps, err := ListDependents()
		if err != nil {
			return err
		}
		if len(deps) == 0 {
			fmt.Println(`None`)
			return nil
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

func argIsOr(args []string, is string, fallback bool) bool {
	if len(args) == 0 {
		return fallback
	}
	return args[0] == is
}
