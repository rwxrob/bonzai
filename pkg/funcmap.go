package bonzai

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/mark"
	"github.com/rwxrob/bonzai/pkg/core/run"
	"github.com/rwxrob/bonzai/pkg/core/term"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

// This file contains the BonzaiMark builtins that Cmd authors can use
// in their Description and other places where templated BonzaiMark is
// allowed.

// FuncMap contains the package global default template domain specific
// language implemented as a collection of functions in
// a template.FuncMap which can be supplemented or overridden by Bonzai
// developers. Note this is in addition to any specific syntax added
// specifically to a Cmd with Cmd.FuncMap (which takes higher priority).
// Note that there is no protection against any [Cmd.Call] function
// changing one or all of these entries for every other Cmd within the
// same executable. This flexibility is by design but must be taken into
// careful consideration when deciding to alter this package-scoped
// variable. It is almost always preferable to use the [Cmd.FuncMap]
// instead.
var FuncMap = template.FuncMap{

	// semantic
	"exe": func(a string) string { return term.Under + a + term.Reset },
	"pkg": func(a string) string { return term.Bold + a + term.Reset },
	"cmd": func(a string) string { return term.Bold + a + term.Reset },

	// stylistic
	"indent": indent,
	"pre":    func(a string) string { return term.Under + a + term.Reset },

	// host system information
	"exepath":      run.ExePath,
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
}

func indent(n int, in string) string {
	return to.Indented(in, mark.IndentBy+n)
}
