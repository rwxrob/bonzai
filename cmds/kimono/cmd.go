package kimono

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/comp/completers/git"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/vars"
)

const (
	WorkScopeEnv   = `KIMONO_WORK_SCOPE`
	WorkScopeVar   = `work-scope`
	TagVerPartEnv  = `KIMONO_VERSION_PART`
	TagVerPartVar  = `tag-ver-part`
	TagShortenEnv  = `KIMONO_TAG_SHORTEN`
	TagShortenVar  = `tag-shorten`
	TagRmRemoteEnv = `KIMONO_TAG_RM_REMOTE`
	TagRmRemoteVar = `tag-rm-remote`
	TagPushEnv     = `KIMONO_PUSH_TAG`
	TagPushVar     = `tag-push`
	TidyScopeEnv   = `KIMONO_TIDY_SCOPE`
	TidyScopeVar   = `tidy-scope`
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Vers:  `v0.7.0`,
	Short: `manage golang monorepos`,
	Long: `
The kimono tool helps manage Go monorepos. It simplifies common monorepo
operations and workflow management.

# Features:
- Toggle go.work files on/off for local development
- Perform coordinated version tagging
- Keep go.mod files tidy across modules
- View dependency graphs and module information
- Track dependent modules and their relationships

# Commands:
- work:     Toggle go.work files for local development
- tidy:     run 'go get -u' and 'go mod tidy' across modules
- tag:      List and coordinate version tagging across modules
- deps:     List and manage module dependencies
- depsonme: List and manage module dependencies
- vars:     View and set configuration variables

Use 'kimono help <command> <subcommand>...' for detailed information
about each command.
`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{
		workCmd,
		tidyCmd,
		tagCmd,
		dependenciesCmd,
		dependentsCmd,
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
  Configures the scope in which to toggle.
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
	Cmds: []*bonzai.Cmd{
		workInitCmd,
		help.Cmd.AsHidden(),
		vars.Cmd.AsHidden(),
	},
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

Run "work init all" to include all dependencies from the monorepo in a 
new go.work file. Alternatively, provide specific module paths to 
initialize a workspace tailored to those dependencies.
`,
	MinArgs:  1,
	RegxArgs: `all`,
	Cmds: []*bonzai.Cmd{
		help.Cmd.AsHidden(),
		vars.Cmd.AsHidden(),
	},
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
	Long: `
The tag command helps with listing, smart tagging of modules in a
monorepo. This ensures that all modules are consistently tagged with the
appropriate module prefix and version numbers, facilitating error-free
version control and release management.
`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{
		tagBumpCmd,
		tagListCmd,
		tagDeleteCmd,
		help.Cmd.AsHidden(),
		vars.Cmd.AsHidden(),
	},
	Def: tagListCmd,
}

var tagListCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `l`,
	Short: `list the tags for the go module`,
	Long: `
The "list" subcommand displays the list of semantic version (semver)
tags for the current Go module. This is particularly useful for
inspecting version history or understanding the current state of version 
tags in your project.

# Behavior

By default, the command lists all tags that are valid semver tags and 
associated with the current module. The tags can be displayed in their 
full form or shortened by setting the KIMONO_TAG_SHORTEN env var.

# Environment Variables

- KIMONO_TAG_SHORTEN: (Defaults to "true")
  Determines whether to display tags in a shortened format, removing 
  the module prefix. It accepts any truthy value.

# Examples

List tags with the module prefix:

    $ export TAG_SHORTEN=false
    $ tag list

List tags in shortened form (default behavior):

    $ KIMONO_TAG_SHORTEN=1 tag list

The tags are automatically sorted in semantic version order.
`,
	Env: bonzai.VarMap{
		TagShortenEnv: bonzai.Var{
			Key: TagShortenEnv,
			Str: "true",
		},
	},
	Vars: bonzai.VarMap{
		TagShortenVar: bonzai.Var{
			Key:  TagShortenVar,
			Bool: true,
		},
	},
	Do: func(x *bonzai.Cmd, args ...string) error {
		shorten := vars.Fetch(
			TagShortenEnv,
			TagShortenVar,
			false,
		)
		each.Println(TagList(shorten))
		return nil
	},
}

var tagDeleteCmd = &bonzai.Cmd{
	Name:  `delete`,
	Alias: `d|del|rm`,
	Short: `delete the given semver tag for the go module`,
	Long: `
The "delete" subcommand removes a specified semantic version (semver) 
tag. This operation is useful for cleaning up incorrect, outdated, or
unnecessary version tags.
By default, the "delete" command only removes the tag locally. To 
delete a tag both locally and remotely, set the TAG_RM_REMOTE 
environment variable or variable to "true". For example:


# Arguments
  tag: The semver tag to be deleted.

# Environment Variables

- TAG_RM_REMOTE: (Defaults to "false")
  Configures whether the semver tag should also be deleted from the 
  remote repository. Set to "true" to enable remote deletion.

# Examples

    $ tag delete v1.2.3
    $ TAG_RM_REMOTE=true tag delete submodule/v1.2.3

This command integrates with Git to manage semver tags effectively.
`,
	Env: bonzai.VarMap{
		TagRmRemoteEnv: bonzai.Var{Key: TagRmRemoteEnv, Str: "false"},
	},
	Vars: bonzai.VarMap{
		TagRmRemoteVar: bonzai.Var{Key: TagRmRemoteVar, Bool: false},
	},
	NumArgs: 1,
	Comp:    comp.Combine{git.CompTags},
	Cmds:    []*bonzai.Cmd{help.Cmd.AsHidden(), vars.Cmd.AsHidden()},
	Do: func(x *bonzai.Cmd, args ...string) error {
		rmRemote := vars.Fetch(
			TagRmRemoteEnv,
			TagRmRemoteVar,
			false,
		)
		return TagDelete(args[0], rmRemote)
	},
}

var tagBumpCmd = &bonzai.Cmd{
	Name:  `bump`,
	Alias: `b|up|i|inc`,
	Short: `bumps semver tags based on given version part.`,
	Long: `
The "bump" subcommand increments the current semantic version (semver) 
tag of the Go module based on the specified version part. This command 
is ideal for managing versioning in a structured manner, following 
semver conventions.

# Arguments
  part: (Defaults to "patch") The version part to increment.
        Accepted values:
          - major (or M): Increments the major version (x.0.0).
          - minor (or m): Increments the minor version (a.x.0).
          - patch (or p): Increments the patch version (a.b.x).

# Environment Variables

- TAG_VER_PART: (Defaults to "patch")
  Specifies the default version part to increment when no argument is 
  passed.

- TAG_PUSH: (Defaults to "false")
  Configures whether the bumped tag should be pushed to the remote 
  repository after being created. Set to "true" to enable automatic 
  pushing. It accepts any truthy value.

# Examples

Increment the version tag locally:

    $ tag bump patch

Automatically push the incremented tag:

    $ TAG_PUSH=true tag bump minor
`,
	Env: bonzai.VarMap{
		TagVerPartEnv: bonzai.Var{Key: TagVerPartEnv, Str: `patch`},
		TagPushEnv:    bonzai.Var{Key: TagPushEnv, Str: `false`},
	},
	Vars: bonzai.VarMap{
		TagVerPartVar: bonzai.Var{Key: TagVerPartVar, Str: `patch`},
		TagPushVar:    bonzai.Var{Key: TagPushVar, Bool: false},
	},
	MaxArgs: 1,
	Opts:    `major|minor|patch|M|m|p`,
	Comp:    comp.CmdsOpts,
	Cmds:    []*bonzai.Cmd{help.Cmd.AsHidden(), vars.Cmd.AsHidden()},
	Do: func(x *bonzai.Cmd, args ...string) error {
		mustPush := vars.Fetch(TagPushEnv, TagPushVar, false)
		if len(args) == 0 {
			part := vars.Fetch(
				TagVerPartEnv,
				TagVerPartVar,
				`patch`,
			)
			return TagBump(optsToVerPart(part), mustPush)
		}
		part := optsToVerPart(args[0])
		return TagBump(part, mustPush)
	},
}

var tidyCmd = &bonzai.Cmd{
	Name:  `tidy`,
	Alias: `tidy|update`,
	Short: "tidy dependencies on all modules in repo",
	Long: `
The "tidy" command updates and tidies the Go module dependencies
across all modules in a monorepo or within a specific scope. This
is particularly useful for maintaining consistency and ensuring
that dependencies are up-to-date.

# Arguments:
  module|mod:          Tidy the current module only.
  repo:                Tidy all modules in the repository.
  deps|dependencies:   Tidy dependencies of the current module in the 
                       monorepo.
  depsonme|dependents: Tidy modules in the monorepo dependent on the 
                       current module.

# Environment Variables:

- KIMONO_TIDY_SCOPE: (Defaults to "module")
  Defines the scope of the tidy operation. Can be set to "module(mod)",
  "root", "dependencies(deps)", or "dependent(depsonme)".

The scope can also be configured using the "tidy-scope" variable or
the "KIMONO_TIDY_SCOPE" environment variable. If no argument is provided,
the default scope is "module".

# Examples:

    # Tidy all modules in the repository
    $ kimono tidy root

    # Tidy only dependencies of the current module in the monorepo
    $ kimono tidy deps

    # Tidy modules in the monorepo dependent on the current module
    $ kimono tidy depsonme

`,
	Env: bonzai.VarMap{
		TidyScopeEnv: bonzai.Var{Key: TidyScopeEnv, Str: `module`},
	},
	Vars: bonzai.VarMap{
		TidyScopeVar: bonzai.Var{Key: TidyScopeVar, Str: `module`},
	},
	MaxArgs: 1,
	Opts:    `module|mod|repo|deps|depsonme|dependencies|dependents`,
	Comp:    comp.Opts,
	Cmds:    []*bonzai.Cmd{help.Cmd.AsHidden(), vars.Cmd.AsHidden()},
	Do: func(x *bonzai.Cmd, args ...string) error {
		var scope string
		if len(args) == 0 {
			scope = vars.Fetch(
				TidyScopeEnv,
				TidyScopeVar,
				`module`,
			)
		} else {
			scope = args[0]
		}
		switch scope {
		case `module`:
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return TidyAll(pwd)
		case `repo`:
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

var dependenciesCmd = &bonzai.Cmd{
	Name:  `dependencies`,
	Alias: `deps`,
	Short: `list or update dependencies`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		help.Cmd.AsHidden(),
		vars.Cmd.AsHidden(),
		dependencyListCmd,
		// dependencyUpdateCmd,
	},
	Def: help.Cmd,
}

var dependencyListCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `on`,
	Short: `list the dependencies for the go module`,
	Long: `
The list subcommand provides a list of all dependencies for the Go
module. The scope of dependencies can be customized using the options
provided. By default, it lists all dependencies.
`,
	NoArgs: true,
	Cmds:   []*bonzai.Cmd{help.Cmd.AsHidden(), vars.Cmd.AsHidden()},
	Do: func(x *bonzai.Cmd, args ...string) error {
		deps, err := ListDependencies()
		if err != nil {
			return err
		}
		each.Println(deps)
		return nil
	},
}

var dependentsCmd = &bonzai.Cmd{
	Name:  `dependents`,
	Alias: `depsonme`,
	Short: `list or update dependents`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		help.Cmd.AsHidden(),
		vars.Cmd.AsHidden(),
		dependentListCmd,
		// dependentUpdateCmd,
	},
	Def: help.Cmd,
}

var dependentListCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `onme`,
	Short: `list the dependents of the go module`,
	Long: `
The list subcommand provides a list of all modules or packages that
depend on the current Go module. This is useful to determine the
downstream impact of changes made to the current module.
`,
	Comp: comp.Cmds,
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
