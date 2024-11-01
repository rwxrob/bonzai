package futil

import (
	"io/fs"
	"os"
	"path/filepath"
)

// CreateDir a new directory with the DefaultDirPerms creating any new
// directories as well (see os.MkdirAll)
func CreateDir(path string) error {
	return os.MkdirAll(path, fs.FileMode(DefaultDirPerms))
}

// DirEntries returns a slice of strings with all the files in the directory
// at that path joined to their path (as is usually wanted). Returns an
// empty slice if empty or path doesn't point to a directory.
func DirEntries(path string) []string {
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

// DirEntriesAddSlashPath returns Entries passed to AddSlash so that all
// directories will have a trailing slash.
func DirEntriesAddSlashPath(path string) []string {
	return DirEntriesAddSlash(DirEntries(path))
}

// DirEntriesAddSlash adds a [filepath.Separator] to the end of all
// entries passed that are directories.
func DirEntriesAddSlash(entries []string) []string {
	var list []string
	for _, entry := range entries {
		if IsDir(entry) {
			entry += string(filepath.Separator)
		}
		list = append(list, entry)
	}
	return list
}

// DirIsEmpty returns true if the directory at path either contains no
// files or only files with zero length. Directories are recursively
// checked. Returns false, however, if the path does not exist.
func DirIsEmpty(path string) bool {
	if NotExists(path) {
		return false
	}
	for _, entry := range DirEntries(path) {

		if IsDir(entry) {
			return DirIsEmpty(entry)
		}

		if !FileIsEmpty(entry) {
			return false
		}

	}
	return true
}

// DirName returns the current working directory name or an empty string.
func DirName() string {
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}

// AbsDir returns the absolute path to the current working directory. If
// unable to determine, returns empty string.
func AbsDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}
