// originally from github.com/inconshreveable/mousetrap (Apache license)

//go:build !windows
// +build !windows

package is

// StartedByExplorer returns true if the program was invoked by the user
// double-clicking on the executable from MS Explorer (explorer.exe). It
// is conservative and returns false if any of the internal calls fail
// or the operating system is not Windows.
func StartedByExplorer() bool {
	return false
}
