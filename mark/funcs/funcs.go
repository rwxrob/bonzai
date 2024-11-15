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
