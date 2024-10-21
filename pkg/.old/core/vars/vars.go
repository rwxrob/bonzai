// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package vars provides high-level functions that are called from the
// Go Bonzai branch of the same name providing universal access to the
// core functionality.
package vars

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	_fs "io/fs"

	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/rwxrob/bonzai/pkg/futil"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

var BlankLine = regexp.MustCompile(`^[ \t\r\n]*$`)

type Map struct {
	sync.Mutex
	M    map[string]string
	Id   string // usually application name
	Dir  string // usually os.UserCacheDir
	File string // usually vars
}

// Map returns a Map with the M initialized. No other initialization is
// performed. See Init.
func New() Map {
	m := Map{}
	m.M = map[string]string{}
	return m
}

// DirPath is the Dir and Id joined.
func (c Map) DirPath() string { return filepath.Join(c.Dir, c.Id) }

// Path returns the combined Dir and File.
func (c Map) Path() string { return filepath.Join(c.Dir, c.Id, c.File) }

// Init initializes the cache directory (Dir) for the current user and
// given application name (Id) using the standard os.UserCacheDir
// location.  The directory is completely removed and new configuration
// file(s) are created.
//
// Consider placing a confirmation prompt before calling this function
// when term.IsInteractive is true. Since Init uses fs/{dir,file}.Create
// you can set the file.DefaultPerms and dir.DefaultPerms if you prefer
// a different default for your permissions.
//
// Permissions in the fs package are restrictive (0700/0600) by default
// to  allow tokens to be stored within configuration files (as other
// applications are known to do). Still, saving of critical secrets is
// not encouraged within any flat file. But anything that a web browser
// would need to cache in order to operate is appropriate (cookies,
// session tokens, etc.).
//
// Fulfills the bonzai.Vars interface.
func (c Map) Init() error {
	d := c.DirPath()

	// safety checks before blowing things away
	if d == "" {
		return fmt.Errorf("could not resolve cache path for %q", c.Id)
	}
	if len(c.Id) == 0 && len(c.Dir) == 0 {
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

	return futil.Touch(c.Path())
}

// Exists returns true if a configuration file exists at Path.
func (c Map) Exists() bool {
	return futil.Exists(c.Path())
}

// SoftInit calls Init if not Exists.
func (c Map) SoftInit() error {
	if !c.Exists() {
		return c.Init()
	}
	return nil
}

// Data returns the cache data in text marshaled format: k=v, no equal
// sign in key, carriage return and line returns escaped, terminated by
// line return on each line. Logs an error if source of data is
// unavailable. Fulfills the bonzai.Vars interface.
func (m Map) Data() string {
	byt, err := os.ReadFile(m.Path())
	if err != nil {
		log.Print(err)
	}
	return string(byt)
}

// Print prints the text version of the cache. See Data for format.
// Fulfills the bonzai.Vars interface.
func (m Map) Print() { fmt.Print(m.Data()) }

// Get returns a value from the persisted cache (if it has one). No
// locking is done.  Fulfills the bonzai.Vars interface.
func (m Map) Get(key string) string {
	m.Load()
	val, _ := m.M[key]
	return val
}

// Set sets a persistent variable in the cache or returns an error if
// not. Fulfills the bonzai.Vars interface.
func (m Map) Set(key, val string) error {
	path := m.Path()
	mod := futil.ModTime(path)
	if mod.IsZero() {
		return fmt.Errorf("failed to read mod time on file: %q", path)
	}
	if err := m.Load(); err != nil {
		return err
	}
	nmod := futil.ModTime(path)
	if mod.IsZero() {
		return fmt.Errorf("failed to read mod time on file: %q", path)
	}
	if nmod.After(mod) {
		return fmt.Errorf("file has changed since read: %q", path)
	}
	prev, has := m.M[key]
	if has {
		defer func() { m.M[key] = prev }()
	}
	m.M[key] = val
	text, err := m.MarshalText()
	if err != nil {
		return err
	}
	return m.OverWrite(string(text))
}

// Del deletes an entry from the persistent cache. Fulfills the
// bonzai.Vars interface.
func (m Map) Del(key string) error {
	if err := m.Load(); err != nil {
		return err
	}
	delete(m.M, key)
	return m.Save()
}

// Loads the latest from File and Unmarshals into M.
func (m *Map) Load() error {
	return m.UnmarshalText([]byte(m.Data()))
}

// Save persists the current map to file. See OverWrite.
func (m *Map) Save() error {
	byt, err := m.MarshalText()
	if err != nil {
		return err
	}
	return m.OverWrite(string(byt))
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

// OverWrite overwrites the cache File in a way that is safe for all
// callers of OverWrite in this current system for any operating system
// using go-internal/lockedfile (taken from the to internal project
// itself, https://github.com/golang/go/issues/33974) but applying the
// file.DefaultPerms instead of the 0666 Go default.
// The format of the cache string must be key:value with carriage
// returns and line returns escaped. No colons are allowed in the
// key. Each line must be terminated with a single line return.
func (c Map) OverWrite(with string) error {
	if err := c.mkdir(); err != nil {
		return err
	}
	return lockedfile.Write(c.Path(),
		strings.NewReader(with), _fs.FileMode(futil.DefaultFilePerms))
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

// UnmarshalText fulfills encoding.TextUnmarshaler interface
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
// fs/file.Edit for more.
func (c Map) Edit() error {
	if err := c.mkdir(); err != nil {
		return err
	}
	path := c.Path()
	if path == "" {
		return fmt.Errorf("unable to locate cache vars for %q", c.Id)
	}
	return futil.Edit(path)
}
