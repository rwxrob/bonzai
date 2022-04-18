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

var markFuncMap = template.FuncMap{

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
}

func indent(n int, in string) string {
	return to.Indented(in, IndentBy+n)
}

func cachedir() string {
	dir, _ := os.UserCacheDir()
	return dir
}

func confdir() string {
	dir, _ := os.UserCacheDir()
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
