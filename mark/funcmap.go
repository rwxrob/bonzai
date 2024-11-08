package mark

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
)

// FuncMap contains the package the BonzaiMark standard template tag
// functions. Any [Renderer] must support them all (even if some may
// output an empty string).
var FuncMap = &template.FuncMap{
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
	// TODO add dummy formatting tags
}
