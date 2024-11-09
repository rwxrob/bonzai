package funcs

import (
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/run"
)

var Map = &mark.Funcs{
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
}
