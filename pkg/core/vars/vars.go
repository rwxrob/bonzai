package vars

// Driver specifies anything that implements a persistence layer for
// high-speed caching of key/value combinations. Implementations may use
// whatever technology for storing the cache but must represent the data
// in traditional key=value pairs with carriage return (\r) and line
// feed (\n) being escaped in the value portion of each line. Equals
// signs may be included in the value without escaping. Blank lines and
// comments are not allowed and must always produce an error.
//
// # Init
//
// Initialize a new cache if one does not yet exist. Must never clear or
// delete a previously initialized cache. Usually this is called within
// an init() function after the other specific configurations of the
// driver have been set (much like database or other driver).
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
// # Match
//
// Match is the same as [Set] except key is matched against a regular
// expression instead of the exact key.
//
// # Load
//
// # Delete
//
// # Fetch
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
	Get(key string) (string, error)    // accessor, "" if non-existent
	Set(key, val string) error         // mutator
	Match(regx string) (string, error) // query
	Load(keyvals string) error         // multiple pairs
	Delete(key string) error           // destroyer
	Data() (string, error)             // k=v with \r and \n escaped in v
	Print() error                      // prints Data to os.Stdout
	Edit() error                       // open default editor, then load
}
