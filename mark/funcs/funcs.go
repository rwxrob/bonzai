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
	"aka":          AKA,
	"code":         Code,
	"cmdtree":      CmdTree,
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

// CmdTree generates and returns a formatted string representation
// of the command tree for the [Cmd] instance and all its [Cmd].Cmds
// subcommands. It aligns [Cmd].Short summaries in the output for better
// readability, adjusting spaces based on the position of the dashes.
func CmdTree(x *bonzai.Cmd) string {
	tree := cmdTree(x, 2)
	lines := to.Lines(tree)
	var widest int
	for _, line := range lines {
		if length := countRunes(line, '←'); length > widest {
			widest = length
		}
	}
	for i, line := range lines {
		parts := strings.Split(line, "←")
		if len(parts) > 1 {
			lines[i] = fmt.Sprintf("    %-*v←%v", widest, parts[0], parts[1])
		} else {
			lines[i] = "    " + line
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func cmdTree(x *bonzai.Cmd, depth int) string {
	out := new(strings.Builder)
	addbranch := func(level int, c *bonzai.Cmd) error {
		if level > depth {
			return nil
		}
		if c.IsHidden() {
			return nil
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
