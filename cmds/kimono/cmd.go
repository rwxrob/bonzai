package kimono

import (
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai"
	varc "github.com/rwxrob/bonzai/cmds/vars"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/vars"
)

var Cmd = &bonzai.Cmd{
	Name:  `kimono`,
	Alias: `kmono|km`,
	Short: `kimono is a tool for managing golang monorepos`,
	Vers:  `0.0.1`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{sanitizeCmd, workCmd, tagCmd},
}

var sanitizeCmd = &bonzai.Cmd{
	Name: `sanitize`,
	Short: `sanitize will run ` + "`go get -u` and `go mod tidy`\n" +
		`on all go modules in the current git repo`,
	Comp: comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return Sanitize()
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
	Cmds:    []*bonzai.Cmd{varc.Cmd},
	Opts:    `major|minor|patch|m|M|p`,
	MaxArgs: 1,
	Call: func(x *bonzai.Cmd, args ...string) error {
		mustPush := stateVar(`push-tags`, TagPushEnv, false)
		var part VerPart
		part = optsToVerPart(
			stateVar(`version-part`, TagVersionPartEnv, `patch`),
		)
		if len(args) == 0 {
			val, err := vars.Data.Get(`default-ver-part`)
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
	Short: `list the tags for the go module`,
	Comp:  comp.Cmds,
	Call: func(x *bonzai.Cmd, args ...string) error {
		shorten := stateVar(`shorten-tags`, TagShortenEnv, false)
		each.Println(TagList(shorten))
		return nil
	},
}

// stateVar retrieves a value by first checking an environment variable.
// If the environment variable does not exist, it checks bonzai.Vars. If
// neither contain a value, it returns the provided fallback.
func stateVar[T any](key, envVar string, fallback T) T {
	if val, exists := os.LookupEnv(envVar); exists {
		return convertValue(val, fallback)
	}
	if val, err := vars.Data.Get(key); err == nil {
		return convertValue(val, fallback)
	}
	return fallback
}

// convertValue attempts to convert a string to the same type as fallback.
func convertValue[T any](val string, fallback T) T {
	var result any = fallback

	switch any(fallback).(type) {
	case string:
		result = val
	case bool:
		result = isTruthy(val)
	case int:
		result, _ = strconv.Atoi(val)
	}

	return result.(T)
}

// isTruthy determines if a string represents a "truthy" value,
// interpreting "t", "true", and positive numbers as true; "f", "false",
// and zero or negative numbers as false.
func isTruthy(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if slices.Contains([]string{"t", "true"}, val) {
		return true
	}
	if slices.Contains([]string{"f", "false"}, val) {
		return false
	}
	if num, err := strconv.Atoi(val); err == nil {
		return num > 0
	}
	return false
}
