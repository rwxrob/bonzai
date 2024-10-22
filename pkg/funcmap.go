package bonzai

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rwxrob/bonzai/pkg/core/mark"
	"github.com/rwxrob/bonzai/pkg/core/term"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

// ExePath holds the full path to the current running process executable
// which is determined at init() time by calling os.Executable and
// passing it to path/filepath.EvalSymlinks to ensure it is the actual
// binary executable file. Errors are reported to stderr, but there
// should never be an error logged unless something is wrong with the Go
// runtime environment.
var ExePath string

// ExeName holds just the base name of the executable without any suffix
// (ex: .exe) and is set at init() time (see ExePath).
var ExeName string

// ExeSymLink holds the name of the symbolic link pointing to the real
// [ExeName]. Note that hard links are indistinguishable.
var ExeSymLink string

func init() {
	var err error
	// get the full path to current running process executable
	ExePath, err = os.Executable()
	if err != nil {
		log.Print(err)
		return
	}
	ExeName = strings.TrimSuffix(
		filepath.Base(ExePath), filepath.Ext(ExePath))
	ExePath, err = filepath.EvalSymlinks(ExePath)
	if err != nil {
		log.Print(err)
	}
	realname := strings.TrimSuffix(
		filepath.Base(ExePath), filepath.Ext(ExePath))
	if realname != ExeName {
		ExeSymLink = ExeName
		ExeName = realname
	}
}

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
	"exepath":     exepath,
	"exename":     exename,
	"exesymlink":  exesymlink,
	"execachedir": execachedir,
	"execonfdir":  execonfdir,
	"cachedir":    cachedir,
	"confdir":     confdir,
	"homedir":     homedir,
	"pathsep":     func() string { return string(os.PathSeparator) },
	"pathjoin":    filepath.Join,
}

func indent(n int, in string) string {
	return to.Indented(in, mark.IndentBy+n)
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

func exepath() string    { return ExePath }
func exename() string    { return ExeName }
func exesymlink() string { return ExeSymLink }

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
