package bonzai

import (
	"os"

	"github.com/rwxrob/bonzai/pkg/core/vars"
)

// VarsDriver specifies the package persistent variables driver interface. All
// implementations must assign themselves to [Vars] package-scoped
// variable during init.
//
// Implementations must persist (cache) simple string key/value
// variables. Implementations of Vars can persist in different ways, but
// most write to [os.UserCacheDir]. Files, network storage, or cloud
// databases, etc. are all allowed and expected.  However, each must
// always present the data in a .key=val format with \r and \n escaped
// and the key never must contain an equal (=). (Equal signs in the
// value are ignored.) This is the fastest format to read and parse.
//
// For maximum compatibility the following codes (which can be easily
// assigned to constants) are returned indicating status of requested
// operation rather than error codes:
//
//	 -1 - unable to successfully complete operation
//		0 - not found
//	 1 - successful
//
// Status codes are particular important when an empty string is
// a valid/preferred value or default.
type VarsDriver interface {
	Get(key string) (string, int) // accessor, "" if non-existent
	Set(key, val string) int      // mutator
	Del(key string) int           // destroyer
	Fetch() (string, int)         // k=v with \r and \n escaped in v
	Print()                       // (printed)
}

var Vars VarsDriver

func init() {
	v := vars.New()
	dir, _ := os.UserCacheDir()
	v.Id = ExeName
	v.Path = dir
	//	v.SoftInit() // FIXME
	Vars = v
}

const (
	FAILED   = -1
	NOTFOUND = 0
	SUCCESS  = 1
)

// Completer specifies a struct with a [Completer.Complete] function
// that will complete the first argument (usually a command of some kind)
// based on the remaining arguments. The [Complete] method must never
// panic and always return at least an empty slice of strings. By
// convention Completers that do not make use of or other arguments
// should use an underscore identifier since they are ignored.
type Completer interface {
	Complete(x any, args ...string) []string
}
