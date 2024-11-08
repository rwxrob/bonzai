package vars

import (
	"log"
	"os"

	"github.com/rwxrob/bonzai/to"
)

// Initialized with a new [Map] for sharing across packages.
var Data Driver

func init() {
	m := NewMap()
	if err := m.Init(); err != nil {
		log.Print(err)
		return
	}
	Data = m
}

// Driver specifies anything that implements a persistence layer for
// high-speed storage and retrieval of key/value combinations.
// Implementations may use whatever technology for persisting the data but
// must represent the data in traditional key=value pairs with carriage
// return (\r) and line feed (\n) being escaped in the value portion of
// each line. Equals signs may be included in the value without
// escaping. Blank lines and comments are not allowed and must always
// produce an error. All data is expected to be in UTF-8. Keys may be
// any valid UTF-8 character except the equals sign (but most
// implementations will limit to keys that will work well with tab
// completion libraries (see [bonzai/comp]).
//
// # Init
//
// Initialize a new persistence store if one does not yet exist. Must
// never clear or delete one that has been previously initialized.
// Usually this is called within an init() function after the other
// specific configurations of the driver have been set (much like
// database or other drivers).
//
// # Clear
//
// Clear removes all persistence (without disposing) restoring the
// persistent state to the same as before it was ever used historically.
// Most implementations will hold on to the same references created
// during [Init].
//
// # Get
//
// Get retrieves a value for a specific key in a case-sensitive way or
// returns an empty string and the [NotFound] error if not found. This
// allows distinguishing an empty string as a value versus having it not
// there. This is somewhat similar to [os.LookupEnv] except that
// a [NotFound] error is used instead.
//
// # Set
//
// Set sets value for a given key. If the key does not exist, creates
// it and then sets. Otherwise, updates the value. For this reason
// [NotFound] is never returned. Implements should return an error if
// unable to persist the new value.
//
// # GrepK
//
// GrepK is the same as [Set] but uses a regular expression (PCRE) to
// return all the k=v pairs that have keys that match.
//
// # GrepV
//
// GrepV is the same as [GrepK] but uses a regular expression (PCRE) to
// return all the k=v pairs that have values that match.
//
// # Load
//
// Load takes a UTF-8 string, parses it, and sets all the key value
// pairs as if individually assigned with [Set].
//
// # Delete
//
// Delete must delete a key from the persistent storage.
//
// # All
//
// Must output all the data in k=v pairs, one per line.
//
// # Print
//
// Must print the [Data] to [os.Stdout].
//
// # Edit
//
// Must open a file containing [Data] for safely editing that will be
// saved when the editor is closed by effectively passing it to [Load].
//
// # KeysWithPrefix
//
// KeysWithPrefix returns a slice of keys from the map that start with the
// specified prefix (pre), refreshing the map first. If an error occurs during
// refresh, it returns an empty slice and the error. This is mandatory
// for any kind of completion implementation.
//
// # Best practices
//
// Usually it is best to keep driver implementation small and
// self-contained within its own package to enable the use of
// packaged-scoped variables during [init] when setting up the specifics
// of a particular driver implementation.
//
// Although the least amount of latency is preferred (generally
// sub-second) a driver makes no guarantees about the time it takes to
// perform any of the interface operations or even if they will ever
// complete at all. Latency requirements and timeouts must be managed
// externally.
type Driver interface {
	Init() error                       // initialize (not clear)
	Clear() error                      // delete ever key=value pair
	Has(key string) bool               // exists?
	Get(key string) (string, error)    // accessor, "" if non-existent
	Set(key, val string) error         // mutator
	GrepK(regx string) (string, error) // returns k=v combos of matches
	GrepV(regx string) (string, error) // returns k=v combos of matches
	Load(keyvals string) error         // multiple pairs
	Delete(key string) error           // destroyer
	All() (string, error)              // k=v with \r and \n escaped
	Print() error                      // prints Data to os.Stdout
	Edit() error                       // open default editor, then load
	KeysWithPrefix(pre string) ([]string, error)
}

// Set saves a key/value pair to the specified file. Returning and error
// if any args are missing or file does not exist or could not be
// created .
func Set(key, value, file string) error {
	m, err := NewMapFromInit(file)
	if err != nil {
		return err
	}
	return m.Set(key, value)
}

// Get retrieves the value associated with the key from the [Map]
// initialized from the specified file. It returns an error if the map
// cannot be initialized.
func Get(key, file string) (string, error) {
	m, err := NewMapFromInit(file)
	if err != nil {
		return "", err
	}
	return m.Get(key)
}

// Fetch retrieves a value of type [T] based on a prioritized lookup
// sequence. It first checks if an environment variable [envVar] is set,
// and if so, attempts to convert the value to type [T] using [fallback]
// as a reference type. If the environment variable is not set, it then
// checks [Data] for a value corresponding to [key]. If found, the value
// is converted to type [T] and returned. If neither lookup succeeds,
// [fallback] is returned as a default.
func Fetch[T any](envVar, key string, fallback T) T {
	if val, exists := os.LookupEnv(envVar); exists {
		return to.Type(val, fallback)
	}
	if val, err := Data.Get(key); err == nil {
		return to.Type(val, fallback)
	}
	return fallback
}
