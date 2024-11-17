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
	WorkScopeVar = `work-scope`
	WorkScopeEnv = `KIMONO_WORK_SCOPE`

	TidyScopeEnv      = `KIMONO_TIDY_SCOPE`
	SanitizeScopeEnv  = `KIMONO_SANITIZE_SCOPE`
	TagPushEnv        = `KIMONO_PUSH_TAG`
	TagShortenEnv     = `KIMONO_SHORTEN_TAG`
	TagVersionPartEnv = `KIMONO_VERSION_PART`
	TagDeleteRemote   = `KIMONO_DELETE_REMOTE_TAG`
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Short: `manage golang monorepos`,
	Vers:  `v0.7.0`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		workCmd,
		sanitizeCmd,
		tagCmd,
		depsCmd,
		vars.Cmd,
		help.Cmd,
	},
	Def: help.Cmd,
}

var workCmd = &bonzai.Cmd{
	Name:  `work`,
	Alias: `w`,
	Short: `toggle go work files on or off`,
	Long: `
Work command toggles the state of Go workspace files (go.work) between
active (on) and inactive (off) modes. This is useful for managing
monorepo development by toggling Go workspace configurations. The scope
in which to toggle the work files can be configured using either the
'work-scope' variable or the 'KIMONO_WORK_SCOPE' environment variable.

# Arguments
  on  : Renames go.work.off to go.work, enabling the workspace.
  off : Renames go.work to go.work.off, disabling the workspace.

# Environment Variables

- KIMONO_WORK_SCOPE: module|repo|tree (Defaults to "module")
Configures the scope.
  - module: Toggles the go.work file in the current module.
  - repo: Toggles all go.work files in the monorepo.
  - tree: Toggles go.work files in the directory tree starting from pwd.
  `,
	Env: bonzai.VarMap{
		WorkScopeEnv: bonzai.Var{Key: WorkScopeEnv, Str: `module`},
	},
	Vars: bonzai.VarMap{
		WorkScopeVar: bonzai.Var{Key: WorkScopeVar, Str: `module`},
	},
	NumArgs:  1,
	RegxArgs: `on|off`,
	Opts:     `on|off`,
	Comp:     comp.CmdsOpts,
	Cmds:     []*bonzai.Cmd{workInitCmd},
	Do: func(x *bonzai.Cmd, args ...string) error {
		root := ``
		var err error
		var from, to string
		invArgsErr := fmt.Errorf("invalid arguments: %s", args[0])
		switch args[0] {
		case `on`:
			from = `go.work.off`
			to = `go.work`
		case `off`:
			from = `go.work`
			to = `go.work.off`
		default:
			return invArgsErr
		}
		// FIXME: the default here could come from Env or Vars.
		scope := vars.Fetch(WorkScopeEnv, WorkScopeVar, `module`)
		switch scope {
		case `module`:
			return WorkToggleModule(from, to)
		case `repo`:
			root, err = getGitRoot()
			if err != nil {
				return err
			}
		case `tree`:
			root, err = os.Getwd()
			if err != nil {
				return err
			}
		}
		return WorkToggleRecursive(root, from, to)
	},
}

var workInitCmd = &bonzai.Cmd{
	Name:  `init`,
	Alias: `i`,
	Short: `new go.work in module using dependencies from monorepo`,
	Long: `
The "init" subcommand initializes a new Go workspace file (go.work) 
for the current module. It helps automate the creation of a workspace
file that includes relevant dependencies, streamlining monorepo
development.

# Arguments
  all:     Automatically generates a go.work file with all module
           dependencies from the monorepo.
  modules: Relative path(s) to modules, same as used with 'go work use'.

# Usage

Run "work init all" to include all dependencies from the monorepo in a 
new go.work file. Alternatively, provide specific module paths to 
initialize a workspace tailored to those dependencies.
`,
	MinArgs:  1,
	RegxArgs: `all`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if args[0] == `all` {
			return WorkGenerate()
		}
		return WorkInit(args...)
	},
}

var sanitizeCmd = &bonzai.Cmd{
	Name:    `sanitize`,
	Alias:   `tidy|update`,
	Opts:    `all|a|deps|depsonme|dependencies|dependents`,
	Short:   "run `go get -u` and `go mod tidy` on all go modules in repo",
	Comp:    comp.Cmds,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return TidyAll(pwd)
		}
		scope := args[0]
		switch scope {
		case ``:
			scope = vars.Fetch(
				SanitizeScopeEnv,
				`sanitize-scope`,
				``,
			)
			fallthrough
		case `all`:
			root, err := futil.HereOrAbove(".git")
			if err != nil {
				return err
			}
			return TidyAll(filepath.Dir(root))
		case `deps`, `dependencies`:
			TidyDependencies()
		case `depsonme`, `dependents`, `deps-on-me`:
			TidyDependents()
		}
		return nil
	},
}

var tagCmd = &bonzai.Cmd{
	Name:  `tag`,
	Alias: `t`,
	Short: `manage or list tags for the go module`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		tagBumpCmd, tagListCmd, tagDeleteCmd, help.Cmd.AsHidden(),
	},
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

func getGitRoot() (string, error) {
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return "", err
	}
	return filepath.Dir(root), nil
}
