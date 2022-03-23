/*
Package trace provides a tracing framework for use by Bonzai applications. Application developers can use the existing flags or create their own. Developers should start their own flags above 100.
*/
package trace

// Flags must always be 0 in production code. Developers may alter it
// during development but no released code should ever make permanent
// changes since collisions between Bonzai command packages would be
// likely.
var Flags = 0

const (
	All = 1 << (iota + 1)

	// like syslog

	Emerg
	Alert
	Error
	Warn
	Notice
	Info
	Debug

	// bonzai specific

	AutoUpdate
)
