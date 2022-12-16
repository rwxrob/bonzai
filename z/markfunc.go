package Z

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// This file contains the BonzaiMark builtins that Cmd authors can use
// in their Description and other places where templated BonzaiMark is
// allowed.

// Dynamic contains the package global default template language which
// can be supplemented or overridden by Bonzai developers. Note this is in
// addition to any specific syntax added specifically to a Cmd with
// Cmd.Dynamic (which takes higher priority). Use Z.Dynamic when
// a shared template language is to be used across all Cmds within
// a single Bonzai tree or branch. This allows powerful, template-driven
// applications to work well with one another.
var Dynamic = template.FuncMap{

	// semantic
	"exe": func(a string) string { return term.Under + a + term.Reset },
	"pkg": func(a string) string { return term.Bold + a + term.Reset },
	"cmd": func(a string) string { return term.Bold + a + term.Reset },

	// stylistic
	"indent": indent,
	"pre":    func(a string) string { return term.Under + a + term.Reset },

	// host system information
	"exepath":     exepath,
	"exename":     exename,
	"execachedir": execachedir,
	"execonfdir":  execonfdir,
	"cachedir":    cachedir,
	"confdir":     confdir,
	"homedir":     homedir,
	"pathsep":     func() string { return string(os.PathSeparator) },
	"pathjoin":    filepath.Join,
}

func indent(n int, in string) string {
	return to.Indented(in, IndentBy+n)
}

func cachedir() string {
	dir, _ := os.UserCacheDir()
	return dir
}

func confdir() string {
	dir, _ := os.UserConfigDir()
	return dir
}

func homedir(a ...string) string {
	dir, _ := os.UserHomeDir()
	extra := filepath.Join(a...)
	path := filepath.Join(dir, extra)
	return path
}

func exepath() string { return ExePath }

func exename() string { return ExeName }

func execachedir(a ...string) string {
	path := filepath.Join(cachedir(), ExeName)
	extra := filepath.Join(a...)
	return filepath.Join(path, extra)
}

func execonfdir(a ...string) string {
	path := filepath.Join(confdir(), ExeName)
	extra := filepath.Join(a...)
	return filepath.Join(path, extra)
}
