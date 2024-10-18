package dir

import (
	"io/fs"
	"os"
	"path/filepath"

	_fs "github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
)

// DefaultPerms are defaults for new directory creation.
var DefaultPerms = 0700

// Create a new directory with the DefaultPerms creating any new
// directories as well (see os.MkdirAll)
func Create(path string) error {
	return os.MkdirAll(path, fs.FileMode(DefaultPerms))
}

// Entries returns a slice of strings with all the files in the directory
// at that path joined to their path (as is usually wanted). Returns an
// empty slice if empty or path doesn't point to a directory.
func Entries(path string) []string {
	var list []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return list
	}
	for _, f := range entries {
		list = append(list, filepath.Join(path, f.Name()))
	}
	return list
}

// EntriesWithSlash returns Entries passed to AddSlash so that all
// directories will have a trailing slash.
func EntriesWithSlash(path string) []string {
	return AddSlash(Entries(path))
}

// AddSlash adds a filepath.Separator to the end of all entries passed
// that are directories.
func AddSlash(entries []string) []string {
	var list []string
	for _, entry := range entries {
		if _fs.IsDir(entry) {
			entry += string(filepath.Separator)
		}
		list = append(list, entry)
	}
	return list
}

// Exists calls fs.Exists and further confirms that the path is
// a directory and not a file.
func Exists(path string) bool { return _fs.Exists(path) && _fs.IsDir(path) }

// IsEmpty returns true if the directory at path either contains no
// files or only files with zero length. Directories are recursively
// checked. Returns false, however, if the path does not exist.
func IsEmpty(path string) bool {
	if _fs.NotExists(path) {
		return false
	}
	for _, entry := range Entries(path) {

		if _fs.IsDir(entry) {
			return IsEmpty(entry)
		}

		if !file.IsEmpty(entry) {
			return false
		}

	}
	return true
}

// Name returns the current working directory name or an empty string.
func Name() string {
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}

// Abs returns the absolute path to the current working directory. If
// unable to determine, returns empty string.
func Abs() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// HereOrAbove returns the full path to the dir if the dir is found in
// the current working directory, or if not exists in any parent
// directory recursively.
func HereOrAbove(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for ; len(dir) > 0 && dir != "/"; dir = filepath.Dir(dir) {
		path := filepath.Join(dir, name)
		if Exists(path) {
			return path, nil
		}
	}
	return "", nil
}
