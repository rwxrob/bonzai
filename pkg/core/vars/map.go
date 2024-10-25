// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package vars provides common ways to work with persistent and
// environment variables as is commonly needed when creating command
// line applications. It includes an implementation of the
// [bonzai.VarsDriver] interface using secure, performant cache file
// persistence.
package vars

import (
	"sync"
)

const PropsFileName = `vars.properties`

// implements [Driver]
type Map struct {
	sync.Mutex
	M    map[string]string
	File string
}

var _ Driver = new(Map)

func NewMap() (*Map, error) {
	// TODO
	return nil, nil
}

func NewMapFrom(file string) (*Map, error) {
	// TODO
	return nil, nil
}

func (m *Map) Init() error {
	// TODO
	// TODO if File is empty, infer it from the current user and binary name
	return nil
}

func (m *Map) Clear() error {
	// TODO
	// TODO if File is empty, infer it from the current user and binary name
	return nil
}

// Get retrieves the value associated with a key, returning an error if
// the key does not exist.
func (m *Map) Get(key string) (string, error) {
	// TODO fetch a fresh copy of cache file if it has changed
	if val, exists := m.M[key]; exists {
		return val, nil
	}
	return "", NotFound{key}
}

// Match retrieves the value associated with a key that matches the
// given regular expression, returning [NotFound] if the key does not
// exist.
func (m *Map) Match(regx string) (string, error) {
	/*
		// TODO fetch a fresh copy of cache file if it has changed
		return "", NotFound{key}
	*/
	return "", nil
}

// Set adds or updates the value associated with a key.
func (m *Map) Set(key, val string) error {
	m.M[key] = val
	// TODO persist if the value has changed
	// TODO throw an error if something else has changed the cache file
	// since last save
	return nil
}

// Delete deletes the key-value pair associated with the given key
// synchronizing the cache file (see [Sync]).
func (m *Map) Delete(key string) error {
	/*
		if _, exists := m.M[key]; !exists {
			return errors.New("key does not exist")
		}
		delete(m.M, key)
		// TODO persist
		return nil
	*/
	return nil
}

// Fetch returns all key-value pairs in the format "k=v" separated by newlines, escaping "\r" and "\n" in values.
func (m *Map) Fetch() (string, error) {
	// TODO sync and return m.String()
	return "", nil
}

// Load accepts a string of key-value pairs and adds them to the map, where pairs are separated by newlines
// and each pair is in the format "k=v".
func (m *Map) Load(keyvals string) error {
	// TODO use existing code
	/*
		lines := strings.Split(keyvals, "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			key := parts[0]
			value := strings.ReplaceAll(parts[1], "\\n", "\n")
			value = strings.ReplaceAll(value, "\\r", "\r")
			m.M[key] = value
		}
	*/
	return nil
}

/*
// From returns a Map parsed and initialized from the properties
// file at path.
func From(filepath string) (Map, error) {
	m := Map{}
	m.File = filepath
	m.Load()
	return m, nil
}

// Loads the latest from [File] using [Fetch] and calls [UnmarshalText]
// to get it into [M].
func (m *Map) Load() error {
	buf, err := m.Fetch()
	if err != nil {
		return err
	}
	return m.UnmarshalText([]byte(buf))
}

// Clear removes all data from the map without resetting any internal
// references. This allows for safe dependencies on the internal but
// exported [M] map directly as applications require.
func (m *Map) Clear() {
	for k := range m.M {
		delete(m.M, k)
	}
}

// UnmarshalText fulfills [encoding.TextUnmarshaler] interface and locks
// while doing so. The internal map reference is not replaced and
// existing values remain unless overwritten by the incoming data.
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

// Fetch returns the persistent cache data in marshaled text format:
//
// * key=value pairs
// * no equal sign in key
// * equal signs valid in value
// * carriage returns (\r) and line feeds (\n) escaped
// * lines ending with unescaped line feed (\n)
//
// Fetch keeps track of the last time it fetched and checks the [os.Stat]
//
// Note that an error does not always indicate no data was retrieved.
func (m *Map) Fetch() (string, error) {
	info, err := os.Stat(m.File)
	if err != nil {
		return "", err
	}
	byt, err := os.ReadFile(m.File)
	return string(byt), err
}

// Print prints the text version of the cache. See Data for format.
// Fulfills the [bonzai.VarsDriver] interface.
func (m Map) Print() { fmt.Print(m.Fetch()) }

// Get returns a value from the persisted cache (if it has one). No
// locking is done.  Fulfills the [bonzai.VarsDriver] interface.
func (m Map) Get(key string) (string, error) {
	if err := m.Load(); err != nil {
		return "", err
	}
	if val, has := m.M[key]; has {
		return val, nil
	}
	return "", NotFound{key}
}

/*

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

// Keys calls [maps.Keys] on the internal map [M].
func (c Map) Keys() []string { return maps.Keys(c.M) }

// Keys calls [maps.KeysWithPrefix] on the internal map [M].
func (c Map) KeysWithPrefix(pre string) []string {
	return maps.KeysWithPrefix(c.M, pre)
}
*/
