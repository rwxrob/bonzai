package vars

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/pkg/core/futil"
)

// Driver specifies anything that implements a persistence layer
// for high-speed caching of key/value combinations. Implementations may
// use whatever technology for storing the cache but must represent the
// data in traditional key=value pairs with carriage return (\r) and
// line feed (\n) being escaped in the value portion of each line.
// Equals signs may be included in the value without escaping. Blank
// lines and comments are not allowed and must always produce an error.
//
// # Init
//
// Initialize a new cache if one does not yet exist. Must never clear or
// delete a previously initialized cache.
//
// # Clear
type Driver interface {
	Init() error                       // initialize (not clear)
	Clear() error                      // delete ever key=value pair
	Get(key string) (string, error)    // accessor, "" if non-existent
	Set(key, val string) error         // mutator
	Match(regx string) (string, error) // query
	Load(keyvals string) error         // multiple pairs
	Delete(key string) error           // destroyer
	Fetch() (string, error)            // k=v with \r and \n escaped in v
}

// Get returns the value from a key/value pair stored in one of the
// following locations:
//
// * [os.Getenv("VARSFILE")] - full path to properties file
// * [os.UserCacheDir] + `vars.properties` if exists
// * [os.Environ] - all of current running environment variables
//
// If the VARSFILE is set but the file does not exist an error is logged
// but searching continues to other sources.
//
// If VARSFILE is present but no key is found search continues.
//
// If vars.properties is found to exist but no key is found search
// continues.
//
// Note that if a key is found and its value is an empty string there is
// no distinction, the empty string is enough to satisfy the lookup and
// the default value is not returned.
//
// The second argument is returned as a default value if a matching key
// is not found.
//
// When several alternative keys are required use [Match] instead.
func Get(key, def string) string {

	// VARSFILE
	file := os.Getenv(`VARSFILE`)
	if len(file) > 0 {
		if !futil.Exists(file) {
			log.Print(`VARSFILE defined but file missing: %v`, file)
		}
		m, err := NewMapFrom(file)
		if err != nil {
			log.Print(err)
		}
		if val, has := m.M[key]; has {
			return val
		}
	}

	// vars.properties
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Print(err)
	}
	file = filepath.Join(dir, PropsFileName)
	if !futil.Exists(file) {
		m, err := NewMapFrom(file)
		if err != nil {
			log.Print(err)
		}
		if val, has := m.M[key]; has {
			return val
		}
	}

	// env
	if val, has := os.LookupEnv(key); has {
		return val
	}

	return def
}

// Match (Q) is like [Get] but dynamically compiles a regular expression
// to lookup the first matching key instead of a direct match.
func Match(regx, def string) string {
	// TODO
	return ""
}

/*
// NewFrom creates a new [Map] and calls [LoadFrom] on it.
func NewFrom() {
}
*/

func LoadFrom() {

}
