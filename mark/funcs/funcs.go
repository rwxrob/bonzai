package funcs

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

var Map = template.FuncMap{
	"exepath":      run.Executable,
	"exename":      run.ExeName,
	"execachedir":  run.ExeCacheDir,
	"exestatedir":  run.ExeStateDir,
	"execonfigdir": run.ExeConfigDir,
	"cachedir":     futil.UserCacheDir,
	"confdir":      futil.UserConfigDir,
	"homedir":      futil.UserHomeDir,
	"statedir":     futil.UserStateDir,
	"pathsep":      func() string { return string(os.PathSeparator) },
	"pathjoin":     filepath.Join,
	"indent":       to.Indented,
	"aka":          AKA,
	"code":         Code,
	"commands":     Commands,
	"command":      Command,
	"summary":      Summary,
	"usage":        Usage,
	"hasenv":       HasEnv,
	"long":         Long,
}

// Long returns the Long description of the command if found dedented so
// that it is left justified completely. Combine with indent when
// needed.
func Long(x *bonzai.Cmd) string {
	out := new(strings.Builder)
	if len(x.Long) > 0 {
		out.WriteString("\n" + to.Dedented(x.Long))
		if x.Long[len(x.Long)-1] != '\n' {
			out.WriteString("\n")
		}
	}
	return out.String()
}

// Summary returns the AKA joined by a long dash with the commands Short
// description if it has one.
func Summary(x *bonzai.Cmd) string {
	if x.Short == "" {
		return AKA(x)
	}
	return AKA(x) + " — " + x.Short
}

// Command returns the name of the command joined to any aliases at the
// end.
func Command(x *bonzai.Cmd) string {
	if x.Alias == "" {
		return x.Name
	}
	return x.Alias + `|` + x.Name
}

// HasEnv returns true if command has declared any environment variables
// in its Vars. Note that inherited vars are not resolved to see if they
// resolved to environment variables (Var.E).
func HasEnv(x *bonzai.Cmd) bool {
	vars := x.VarsSlice()
	if len(vars) > 0 {
		for _, v := range vars {
			if v.E != "" {
				return true
			}
		}
	}
	return false
}

// Usage return the [Command] plus any available Usage information. If
// there is no usage, it is inferred.
func Usage(x *bonzai.Cmd) string {
	usage := x.Usage
	if usage == "" {
		if len(x.Cmds) > 0 {
			usage = `COMMAND`
		}
		if x.Alias != "" {
			usage += `|` + x.Alias
		}
	}
	return Command(x) + " " + usage
}

// AKA returns the name followed by all aliases in parenthesis joined
// with a forward bar (|) suitable for inlining within help
// documentation. It is available as aka in [Map] as well.
func AKA(x *bonzai.Cmd) string {
	aliases := x.Aliases()

	switch len(aliases) {
	case 0:
		return ""
	case 1:
		return "`" + aliases[0] + "`"
	default:
		aliases = aliases[:len(aliases)-1]
	}

	for n, a := range aliases {
		aliases[n] = "`" + a + "`"
	}

	return "`" + x.Name + "`" + " (" + strings.Join(aliases, "|") + ")"
}

// Code returns a string with Markdown backticks surrounding it after
// converting it to a string with [fmt.Printf]. This is also available
// as "code" in [Map]. This fulfills a  specific use case when
// a developer would like to use backticks in a [bonzai.Cmd].Long or
// [bonzai.Cmd].Short but cannot because backticks are already used to
// contain the multi-line text itself.
func Code(it any) string { return fmt.Sprintf("`%v`", it) }

// Commands generates and returns a formatted string representation
// of the commands and subcommands for the [Cmd] instance.
// It aligns [Cmd].Short summaries in the output for better
// readability, adjusting spaces based on the position of the dashes.
func Commands(x *bonzai.Cmd) string {
	tree := CmdTree(x, 2)
	tree = to.PrefixTrimmed(`  `, tree)
	lines := to.Lines(tree)[1:]
	var widest int
	for _, line := range lines {
		if length := countRunes(line, '←'); length > widest {
			widest = length
		}
	}
	for i, line := range lines {
		parts := strings.Split(line, "←")
		if len(parts) > 1 {
			lines[i] = fmt.Sprintf("%-*v←%v", widest, parts[0], parts[1])
		} else {
			lines[i] = line
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func CmdTree(x *bonzai.Cmd, depth int) string {
	out := new(strings.Builder)
	hideunder := -1
	addbranch := func(level int, c *bonzai.Cmd) error {
		if level > depth {
			return nil
		}
		if c.IsHidden() {
			hideunder = level
			return nil
		}
		if hideunder >= 0 && level > hideunder {
			return nil
		}
		if hideunder >= 0 && level < hideunder {
			hideunder = -1
		}
		for range level {
			out.WriteString(`  `)
		}
		name := c.Name
		if len(name) == 0 {
			name = `noname`
		}
		out.WriteString(name)
		if len(c.Short) > 0 {
			out.WriteString(" ← " + c.Short)
		}
		caller := c.Caller()
		if caller != nil && caller.Def == c {
			if len(c.Short) == 0 {
				out.WriteString(" ←")
			}
			out.WriteString(" (default)")
		}
		out.WriteString("\n")
		return nil
	}
	x.WalkDeep(addbranch, nil)
	return out.String()
}

// strings.Index and IndexRune don't do what you think
func countRunes(in string, it rune) int {
	var i int
	for _, r := range []rune(in) {
		if r == it {
			return i
		}
		i++
	}
	return i
}
