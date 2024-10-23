// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package vars provides an implementation of the [bonzai.VarsDriver]
// interface employing high-performance file persistence using the
// [os.UserCacheDir] directory containing a single standard properties
// text file per [Id].
package vars

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

var BlankLine = regexp.MustCompile(`^[ \t\r\n]*$`)

type Map struct {
	sync.Mutex
	M    map[string]string
	Id   string // usually application name
	Path string // usually os.UserCacheDir
}

// Map returns a Map with the M initialized. No other initialization is
// performed. See Init.
func New() Map {
	m := Map{}
	m.M = map[string]string{}
	return m
}

// DirPath is the [Path] and [Id] joined.
func (c Map) DirPath() string { return filepath.Join(c.Path, c.Id) }

// FullPath returns the combined [Path], [Id], and `vars.properties` string.
func (c Map) FullPath() string { return filepath.Join(c.Path, c.Id, `vars.properties`) }

// Init initializes the cache directory ([Path]) for the current user and
// given application name ([Id]) using the standard [os.UserCacheDir]
// location. The directory is completely removed and new configuration
// file(s) are created.
//
// Consider placing a confirmation prompt before calling this function
// when [term.IsInteractive] is true. Since Init uses fs/{dir,file}.Create
// you can set the file.DefaultPerms and dir.DefaultPerms if you prefer
// a different default for your permissions.
//
// Permissions in the fs package are restrictive (0700/0600) by default
// to  allow tokens to be stored within configuration files (as other
// applications are known to do). Still, saving of critical secrets is
// not encouraged within any flat file. But anything that a web browser
// would need to cache in order to operate is appropriate (cookies,
// session tokens, etc.).
func (c Map) Init() error {
	d := c.DirPath()

	// safety checks before blowing things away
	if d == "" {
		return fmt.Errorf("could not resolve cache path for %q", c.Id)
	}
	if len(c.Id) == 0 && len(c.Path) == 0 {
		return fmt.Errorf("empty directory id")
	}

	if futil.Exists(d) {
		if err := os.RemoveAll(d); err != nil {
			return err
		}
	}

	if err := futil.CreateDir(d); err != nil {
		return err
	}

	return futil.Touch(c.FullPath())
}

// Exists returns true if a configuration file exists at [FullPath].
func (c Map) Exists() bool {
	return futil.Exists(c.FullPath())
}

// SoftInit calls Init if not Exists.
func (c Map) SoftInit() error {
	if !c.Exists() {
		return c.Init()
	}
	return nil
}

// Fetch returns all the persistent cache data in text marshaled format:
// k=v, no equal sign in key, carriage return and line returns escaped,
// terminated by line return on each line. Logs an error and returns -1
// (FAILED) if source of data is unavailable.
func (m Map) Fetch() (string, int) {
	byt, err := os.ReadFile(m.FullPath())
	if err != nil {
		log.Print(err)
		return "", -1
	}
	return string(byt), 1
}

// Print prints the text version of the cache. See Data for format.
// Fulfills the bonzai.Vars interface.
func (m Map) Print() { fmt.Print(m.Fetch()) }

// Get returns a value from the persisted cache (if it has one). No
// locking is done.  Fulfills the [bonzai.Vars] interface.
func (m Map) Get(key string) (string, int) {
	m.Load()
	val, has := m.M[key]
	if !has {
		return val, 0
	}
	return val, 1
}

// Set sets a persistent variable in the cache or returns -1 (FAILED) if
// not. Fulfills the [bonzai.Vars] interface. Never returns 0 (NOTFOUND)
// since creates if not. Errors are logged.
func (m Map) Set(key, val string) int {
	path := m.FullPath()
	mod := futil.ModTime(path)
	if mod.IsZero() {
		log.Printf("failed to read mod time on file: %q", path)
		return -1
	}
	if err := m.Load(); err != nil {
		log.Print(err)
		return -1
	}
	nmod := futil.ModTime(path)
	if mod.IsZero() {
		log.Printf("failed to read mod time on file: %q", path)
		return -1
	}
	if nmod.After(mod) {
		log.Printf("file has changed since read: %q", path)
		return -1
	}
	prev, has := m.M[key]
	if has {
		defer func() { m.M[key] = prev }()
	}
	m.M[key] = val
	text, err := m.MarshalText()
	if err != nil {
		log.Print(err)
		return -1
	}
	err = m.Overwrite(string(text))
	if err != nil {
		log.Print(err)
		return -1
	}
	return 1
}

// Del deletes an entry from the persistent cache. Fulfills the
// [bonzai.Vars] interface. Logs errors and return -1 (FAILED) on
// failure. If key was not found does nothing. Never returns 0 (NOTFOUND)
// nor 1 (SUCCESS) since go's delete operator does not return any
// status. This operation fully loads the full cache and overwrites it
// every time. See [Load] and [Save].
func (m Map) Del(key string) int {
	if err := m.Load(); err != nil {
		log.Print(err)
		return -1
	}
	delete(m.M, key)
	if err := m.Save(); err != nil {
		return -1
	}
	return 1
}

// Loads the latest from File and Unmarshals into M.
func (m *Map) Load() error {
	buf, code := m.Fetch()
	if code == 1 {
		return m.UnmarshalText([]byte(buf))
	}
	return fmt.Errorf(`failed to retrieve cache while loading`)
}

// Save persists the current map to file. See Overwrite. See [Load] and
// [Save].
func (m *Map) Save() error {
	byt, err := m.MarshalText()
	if err != nil {
		return err
	}
	return m.Overwrite(string(byt))
}

func (c Map) mkdir() error {
	d := c.DirPath()
	if d == "" {
		return fmt.Errorf("failed to find config for %q", c.Id)
	}
	if futil.NotExists(d) {
		if err := futil.CreateDir(d); err != nil {
			return err
		}
	}
	return nil
}

// Overwrite overwrites the cache File in a way that is safe for all
// callers from any executable on this machine (see [futil.Overwrite]).
// The format of the cache string must be key=value with carriage
// returns and line returns escaped. No colons are allowed in the
// key. Each line must be terminated with a single line return.
func (c Map) Overwrite(with string) error {
	if err := c.mkdir(); err != nil {
		return err
	}
	return futil.Overwrite(c.FullPath(), with)
}

// MarshalText fulfills encoding.TextMarshaler interface
func (c Map) MarshalText() ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	var out string
	for k, v := range c.M {
		out += k + "=" + to.EscReturns(v) + "\n"
	}
	return []byte(out), nil
}

// UnmarshalText fulfills encoding.TextUnmarshaler interface.
func (c *Map) UnmarshalText(in []byte) error {
	c.Lock()
	defer c.Unlock()
	lines := to.Lines(string(in))
	for _, line := range lines {
		if BlankLine.MatchString(line) {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			c.M[parts[0]] = to.UnEscReturns(parts[1])
		}
	}
	return nil
}

// Edit opens the given cached variables files in the local editor. See
// [futil.Edit] for more.
func (c Map) Edit() error {
	if err := c.mkdir(); err != nil {
		return err
	}
	path := c.FullPath()
	if path == "" {
		return fmt.Errorf("unable to locate cache vars for %q", c.Id)
	}
	return futil.Edit(path)
}
