package futil

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Tilde2Home expands a Tilde (~) prefix into a proper os.UserHomeDir path. If
// it cannot find os.UserHomeDir simple returns unchanged path. Will not
// expand for specific users (~username).
func Tilde2Home(dir string) string {
	if !strings.HasPrefix(dir, `~`) {
		return dir
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		return path.Join(home, dir[1:])
	}
	return dir
}

// IsDir returns true if the indicated path is a directory. Returns
// false and logs if unable to determine.
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Print(err)
		return false
	}
	return info.IsDir()
}

// DupPerms duplicates the perms of one file or directory onto the second.
func DupPerms(orig, clone string) error {
	ostats, err := os.Stat(orig)
	if err != nil {
		return err
	}
	_, err = os.Stat(clone)
	if err != nil {
		return nil
	}
	return os.Chmod(clone, ostats.Mode())
}

// ModTime returns the [FileInfo.ModTime] from [Stat] as a convenience
// or returns a zero time if not. See [time.IsZero].
func ModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err == nil {
		return info.ModTime()
	}
	var t time.Time
	return t
}

// Exists returns true if the given path was absolutely found to exist
// on the system. A false return value means either the file does not
// exists or it was not able to determine if it exists or not. Use
// [NotExists] instead.
//
// WARNING: do not use this function if a definitive check for the
// non-existence of a file is required since the possible indeterminate
// error state is a possibility. These checks are also not atomic on
// many file systems so avoid this usage for pseudo-semaphore designs
// and depend on file locks.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// NotExists definitively returns true if the given path does not exist.
// See Exists as well.
func NotExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

type ErrNotExist struct {
	N string
}

func (e ErrNotExist) Error() string {
	return fmt.Sprintf("file or directory does not exist: %v", e.N)
}

type ErrExist struct {
	N string
}

func (e ErrExist) Error() string {
	return fmt.Sprintf("file or directory already exists: %v", e.N)
}

// HereOrAbove returns the full path to the file or directory if it is
// found in the current working directory, or if not exists in any
// parent directory recursively. Returns ErrNotExist error with the name
// if not found.
func HereOrAbove(name string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for ; len(dir) > 0 && dir != string(os.PathSeparator); dir = filepath.Dir(dir) {
		path := filepath.Join(dir, name)
		if Exists(path) {
			return path, nil
		}
	}
	return "", ErrNotExist{name}
}

// IsDirFS simply a shortcut for fs.Stat().IsDir(). Only returns true if
// the path is a directory. If not a directory (or an error prevented
// confirming it is a directory) then returns false.
func IsDirFS(fsys fs.FS, path string) bool {
	info, err := fs.Stat(fsys, path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

var (
	ExtractFilePerms = fs.FileMode(0600)
	ExtractDirPerms  = fs.FileMode(0700)
)

// ExtractEmbed walks the embedded file system and duplicates its
// structure into the target directory beginning with root. Since
// embed.FS drops all permissions information from the original files
// all files and directories are written as read/write (0600) for the
// current effective user (0600 for file, 0700 for directories). This
// default can be changed by setting the package variables
// ExtractFilePerms and ExtractDirPerms. Note that each embedded file is
// full buffered into memory before writing.
func ExtractEmbed(it embed.FS, root, target string) error {
	return fs.WalkDir(it, root,

		func(path string, i fs.DirEntry, err error) error {

			to := filepath.Join(target, strings.TrimPrefix(path, root))

			if i.IsDir() {
				return os.MkdirAll(to, ExtractDirPerms)
			}

			buf, err := fs.ReadFile(it, path)
			if err != nil {
				return err
			}

			return os.WriteFile(to, buf, ExtractFilePerms)
		})

}

// Paths returns a list of full paths to each of the directories or
// files from the root but with the path to the root stripped resulting
// in relative paths.
func RelPaths(it fs.FS, root string) []string {
	var paths []string
	fs.WalkDir(it, root,
		func(path string, i fs.DirEntry, err error) error {
			to := strings.TrimPrefix(path, root)
			if to == "" {
				return nil
			}
			paths = append(paths, to[1:])
			return nil
		})
	return paths
}

// LatestChange walks the directory rooted at root looking at each file or
// directory within it recursively (including itself) and returns the
// time of the most recent change along with its full path. LastChange
// returns nil FileInfo and empty string if root is not a directory.
func LatestChange(root string) (string, fs.FileInfo) {

	if !IsDir(root) {
		return "", nil
	}

	latest := struct {
		Info fs.FileInfo
		Path string
	}{}

	err := filepath.WalkDir(root, func(p string, f fs.DirEntry, _ error) error {

		if latest.Path == "" {
			i, err := f.Info()
			if err != nil {
				return err
			}
			latest.Info = i
			latest.Path = p
			return nil
		}
		i, err := f.Info()
		if err != nil {
			return nil
		}
		if i.ModTime().After(latest.Info.ModTime()) {
			latest.Info = i
			latest.Path = p
		}
		return nil
	})
	if err != nil {
		return "", nil
	}
	return latest.Path, latest.Info
}

// PathEntry contains the fully qualified path to a DirEntry and the
// returned FileInfo for that specific path. PathEntry saves the work of
// fetching this information a second time when functions in this
// package have already retrieved it.
type PathEntry struct {
	Path string
	Info fs.FileInfo
}

// NameIsInt returns true if the filepath.Base name of the passed path
// is a valid positive integer (including 0).
func NameIsInt(path string) bool {
	name := filepath.Base(path)
	if i, err := strconv.Atoi(name); i < 0 || err != nil {
		return false
	}
	return true
}

// IntDirs returns all the directory entries within the target directory
// that have integer names. The lowest integer and highest integer
// values are also returned. Only positive integers are checked. This is
// useful when using directory names as database-friendly unique primary
// keys for other file system content.
//
// IntDirs returns an empty slice and -1 values if no matches are
// found.
//
// Errors looking up the FileInfo cause Info to be nil.
func IntDirs(target string) (paths []PathEntry, low, high int) {
	low, high = -1, -1
	entries, err := os.ReadDir(target)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		var val int
		name := entry.Name()
		if v, err := strconv.Atoi(name); err != nil || v < 0 {
			continue
		} else {
			val = v
		}
		if low < 0 {
			low = val
		}
		if high < 0 {
			high = val
		}
		if val > high {
			high = val
		}
		if val < low {
			low = val
		}
		var pe PathEntry
		if abs, err := filepath.Abs(filepath.Join(target, name)); err == nil {
			pe.Path = abs
		}
		if i, err := entry.Info(); err == nil {
			pe.Info = i
		}
		paths = append(paths, pe)
	}
	return
}

// Preserve moves the target to a new name with an "~" isosec suffix
// usually in anticipation of eventual deletion or restoration for
// transactionally safe dealings with directories and files. Returns the
// name of the new file or directory. Does *not* return an error if the
// target doesn't exist (since there is nothing to preserve) and returns
// an empty string. Pair this with defer RevertIfMissing to provide
// transactional backups of files and directories.
func Preserve(target string) (string, error) {
	if NotExists(target) {
		return "", nil
	}
	nname := target + `~` + Isosec()
	if err := os.Rename(target, nname); err != nil {
		return "", err
	}
	return nname, nil
}

// Isosec returns the GMT current time in ISO8601 (RFC3339) without
// any punctuation or the T.  This is frequently a very good unique
// suffix that has the added advantage of being chronologically sortable
// and more readable than the epoch. (Also see Second())
func Isosec() string {
	return fmt.Sprintf("%v", time.Now().In(time.UTC).Format("20060102150405"))
}

// RevertIfMissing is designed to be called from defer after calling
// a Preserve and storing the backup path to a variable. If the target
// is missing, backup is renamed back to the target. If the target is
// *not* missing, backup is removed. If the backup is an empty string no
// backup restoration is assumed and a nil error is returned with no
// action. This allows the output of Preserve to be passed directly to
// RevertIfMissing without complicating the defer.
func RevertIfMissing(target, backup string) error {
	if target == "" {
		return nil
	}
	if NotExists(target) {
		if err := os.Rename(backup, target); err != nil {
			return err
		}
	}
	if Exists(target) {
		if err := os.RemoveAll(backup); err != nil {
			return err
		}
	}
	return nil
}

// IsSymLink returns true if the file at the base of path (not any
// intermediate directories) is itself a symbolic link.
func IsSymLink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

var UserHomeDir = os.UserHomeDir

// UserConfigDir returns the path to the user's config directory,
// checking [$XDG_CONFIG_HOME] environment variable first. It defaults to
// joining [UserHomeDir] with `.config` if the variable is unset on Unix
// systems. Note that this is unlike [pkg/os.UserConfigDir] by design
// for better intuitive consistency despite the defined "standards" from
// operating system vendors. This difference is particularly important
// given the prevalence of WSL2 usage requiring a clear distinction.
func UserConfigDir() (string, error) {
	path, has := os.LookupEnv(`XDG_CONFIG_HOME`)
	if !has {
		dir, err := UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, `.config`)
	}
	return path, nil
}

// UserCacheDir returns the path to the user's cache directory,
// checking [$XDG_CACHE_HOME] environment variable first. It defaults to
// joining [UserHomeDir] with `.cache` if the variable is unset on Unix
// systems (including darwin unlike [pkg/os.UserCacheDir]). On Windows it
// returns %LocalAppData%. It returns an error if [UserHomeDir] fails.
func UserCacheDir() (string, error) {
	if runtime.GOOS == "windows" {
		return os.Getenv(`LocalAppData`), nil
	}
	path, has := os.LookupEnv(`XDG_CACHE_HOME`)
	if !has {
		dir, err := UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, `.cache`)
	}
	return path, nil
}

// UserStateDir returns the path to the user's state directory,
// checking [XDG_STATE_HOME] environment variable first. It defaults
// to joining [UserHomeDir] with `.local/state` if the variable is
// unset. It returns an error if [UserHomeDir] fails.
func UserStateDir() (string, error) {
	path, has := os.LookupEnv(`XDG_STATE_HOME`)
	if !has {
		dir, err := UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(dir, `.local`, `state`)
	}
	return path, nil
}
