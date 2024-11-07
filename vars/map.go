// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package vars provides common ways to work with persistent and as is
// commonly needed when creating command line applications.
package vars

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/fn/maps"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

const DefaultFileName = `vars.properties`

type Map struct {
	sync.Mutex
	M    map[string]string
	File string

	lastload time.Time
}

// Map implements [Driver]
var _ Driver = new(Map)

// All returns the state data file as text marshaled in the format: k=v,
// no equal sign in key, carriage return and line returns escaped,
// terminated by line return on each line. Logs an error if source of
// data is unavailable. Fulfills the [Driver] interface.
func (m *Map) All() (string, error) {
	byt, err := os.ReadFile(m.File)
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

// Print retrieves and displays the data from the file [m.File]
// by calling the [All] method. Returns any error encountered during
// the data retrieval.
func (m *Map) Print() error {
	out, err := m.All()
	fmt.Print(out)
	return err
}

// NewMap returns a pointer to a [Map] without any additional
// initialization.
func NewMap() *Map {
	m := new(Map)
	m.M = map[string]string{}
	return m
}

// NewMapFrom calls [NewMap], sets [Map.File], and [ values from the
// parsed file. If any error is encountered it is returned along with
// the new map. The map is always returned (never nil).
func NewMapFrom(file string) (*Map, error) {
	m := NewMap()
	if len(file) == 0 {
		return nil, ErrMissingArg{`file`}
	}
	m.File = file
	err := m.loadFile(m.File)
	return m, err
}

// NewMapFromInit is the same as [NewMapFrom] except it creates a file
// at the location if it does not already exist by calling [Init] on
// a new [Map] created from [NewMap].
func NewMapFromInit(file string) (*Map, error) {
	var m *Map
	if len(file) == 0 {
		return nil, ErrMissingArg{`file`}
	}
	if futil.Exists(file) {
		var err error
		m, err = NewMapFrom(file)
		if err != nil {
			return nil, err
		}
	} else {
		m = NewMap()
		m.File = file
		if err := m.Init(); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// loadFile calls [Load] after buffering the file.
func (c *Map) loadFile(file string) (err error) {
	var f *os.File
	f, err = os.Open(file)
	if err != nil {
		return
	}

	defer func() {
		closeerr := f.Close()
		err = errors.Join(closeerr, err)
	}()

	var buf []byte
	buf, err = io.ReadAll(f)
	if err != nil {
		return
	}

	var info fs.FileInfo
	info, err = f.Stat()
	if err != nil {
		return
	}

	c.lastload = info.ModTime()
	return c.Load(string(buf))
}

// Load accepts a string of key-value pairs and adds them to the map,
// where pairs are separated by newlines and each pair is in the format
// "k=v". Implements [Driver].
func (m *Map) Load(keyvals string) error {
	m.refresh()
	err := m.UnmarshalText([]byte(keyvals))
	if err != nil {
		return err
	}
	return m.save()
}

// UnmarshalText fulfills [encoding.TextUnmarshaler] interface and locks
// while doing so. The internal map reference is not replaced and
// existing values remain unless overwritten by the incoming data.
func (c *Map) UnmarshalText(in []byte) error {
	c.Lock()
	defer c.Unlock()
	lines := to.Lines(string(in))
	for _, line := range lines {
		parts := strings.SplitN(line, `=`, 2)
		if len(parts) == 2 {
			c.M[parts[0]] = to.UnEscReturns(parts[1])
		}
	}
	return nil
}

// fileHasChanged returns true if the cache file has a newer last
// modified time [fs.FileInfo] ModeTime() that is after that of the last
// operation. Note that true is returned if there is any error (which is
// logged).
func (m *Map) fileHasChanged() bool {
	info, err := os.Stat(m.File)
	if err != nil {
		return true // safe to trigger a new load
	}
	return info.ModTime().After(m.lastload)
}

// Refresh calls [Map.loadFile] on itself pointing to its internal
// [Map.File] if [Map.fileHasChanged]. No attempt to create the cache
// file is made if missing. For that use [Map.Init] instead.
func (m *Map) refresh() error {
	if !m.fileHasChanged() {
		return nil
	}
	return m.loadFile(m.File)
}

// defaultFile returns the default file path for the Map by combining the
// executable cache directory with the properties file name. It calls the
// RealExeStateDir method from the run package to get the directory. If
// there is an error retrieving the directory, it returns an empty string
// and the error.
func (m *Map) defaultFile() (string, error) {
	dir, err := run.RealExeStateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, DefaultFileName), nil
}

// Init checks if the [Map]'s [File] is set and exists. If not, uses
// a default Ensures file existence with [futil.Touch]. Returns an error
// if operations fail. Fulfills the [Driver] interface.
func (m *Map) Init() error {
	file := m.File
	var err error
	if len(file) == 0 {
		file, err = m.defaultFile()
		if err != nil {
			return err
		}
		m.File = file
	}
	if futil.Exists(file) {
		return nil
	}
	return futil.Touch(file)
}

// Edit opens the default editor to edit the file specified in m.File.
func (m *Map) Edit() error { return edit.Files(m.File) }

// Clear removes all elements from the map [m.M] while ensuring thread safety
// by locking before the operation and unlocking afterward.
func (m *Map) Clear() error {
	m.Lock()
	maps.Clear(m.M)
	m.Unlock()
	if len(m.File) > 0 {
		return m.save()
	}
	return nil
}

// Get retrieves the value associated with a key, returning an error if
// the key does not exist. Get calls [Map.FileHasChanged] and if true
// calls [Map.Load].
func (m *Map) Get(key string) (string, error) {
	m.refresh()
	if val, exists := m.M[key]; exists {
		return val, nil
	}
	return "", ErrNotFound{key}
}

// Has checks if the map [m.M] contains the specified [key]
// after refreshing the map's data. Returns true if the
// key is present; otherwise, it returns false.
func (m *Map) Has(key string) bool {
	m.refresh()
	_, has := m.M[key]
	return has
}

// GrepK returns all key-value pairs associated with a key that matches
// the given regular expression. An empty string is a valid (non-error)
// result indicating nothing matched. Fulfills the [Driver] interface.
func (m *Map) GrepK(regx string) (string, error) {
	var buf strings.Builder
	if len(m.File) > 0 {
		if err := m.refresh(); err != nil {
			return "", err
		}
	}
	x := regexp.MustCompile(regx)
	for k, v := range m.M {
		if x.MatchString(k) {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(v)
			buf.WriteString("\n")
		}
	}
	return buf.String(), nil
}

// GrepV returns all key-value pairs associated with a value that matches
// the given regular expression. Fulfills the [Driver] interface.
func (m *Map) GrepV(regx string) (string, error) {
	var buf strings.Builder
	if len(m.File) > 0 {
		if err := m.refresh(); err != nil {
			return "", err
		}
	}
	x := regexp.MustCompile(regx)
	for k, v := range m.M {
		if x.MatchString(v) {
			buf.WriteString(k)
			buf.WriteString(`=`)
			buf.WriteString(v)
			buf.WriteString("\n")
		}
	}
	return buf.String(), nil
}

// Set adds or updates the value associated with a key.
func (m *Map) Set(key, val string) error {
	m.refresh()
	if cur, has := m.M[key]; has && val == cur {
		return nil
	}
	m.M[key] = val
	return m.save()
}

// Save persists the current map to file. See OverWrite.
func (m *Map) save() error {
	byt, err := m.MarshalText()
	if err != nil {
		return err
	}
	return futil.Overwrite(m.File, string(byt))
}

// MarshalText fulfills [encoding.TextMarshaler] interface.
func (c *Map) MarshalText() ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	lines := make([]string, 0, len(c.M))
	for k, v := range c.M {
		lines = append(lines, k+"="+to.EscReturns(v))
	}
	slices.Sort(lines)
	buf := new(bytes.Buffer)
	for _, line := range lines {
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf.Bytes(), nil
}

// Delete deletes an entry from the persistent cache. Fulfills the
// [Driver] interface.
func (m *Map) Delete(key string) error {
	if err := m.refresh(); err != nil {
		return err
	}
	delete(m.M, key)
	return m.save()
}

// KeysWithPrefix returns a slice of keys from the map [m.M] that start with the
// specified prefix (pre), refreshing the map first. If an error occurs during
// refresh, it returns an empty slice and the error. Note that this
// linearly passes through every map rather than attempt to sort the map
// keys first.
func (m *Map) KeysWithPrefix(pre string) ([]string, error) {
	if err := m.refresh(); err != nil {
		return []string{}, err
	}
	list := maps.KeysWithPrefix(m.M, pre)
	return list, nil
}
