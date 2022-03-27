package cache

// Cacher specifies high-performance caching for Bonzai commands with
// guaranteed system-wide concurrency safety (across process
// boundaries). Safe and strong caching support is a key Bonzai design
// principle since much of what might have been put into getopt dashed
// options and keys before is now done with contextual session cache to
// preserve state between Bonzai command calls. It is not uncommon for
// a given task to span multiple command calls. This also provides
// support for operating-system independent command invocation history.
type Cacher interface {
	// TODO
}
