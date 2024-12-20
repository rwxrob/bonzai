package injson

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/rwxrob/bonzai/futil"
)

type Persister struct {
	File string
}

func NewUserConfig(name, file string) *Persister {
	this := new(Persister)
	dir, err := futil.UserConfigDir()
	if err != nil {
		panic(err)
	}
	f := filepath.Join(dir, name, file)
	err = futil.Touch(f)
	if err != nil {
		panic(err)
	}
	this.File = f
	return this
}

func NewUserCache(name, file string) *Persister {
	this := new(Persister)
	dir, err := futil.UserCacheDir()
	if err != nil {
		panic(err)
	}
	f := filepath.Join(dir, name, file)
	err = futil.Touch(f)
	if err != nil {
		panic(err)
	}
	this.File = f
	return this
}

func NewUserState(name, file string) *Persister {
	this := new(Persister)
	dir, err := futil.UserStateDir()
	if err != nil {
		panic(err)
	}
	f := filepath.Join(dir, name, file)
	err = futil.Touch(f)
	if err != nil {
		panic(err)
	}
	this.File = f
	return this
}

// Setup ensures that the persistence file exists and is ready for use.
// It opens the file in read-only mode, creating it if it doesn't
// already exist, and applies secure file permissions (0600) to restrict
// access. The file is immediately closed after being verified or
// created. If the file cannot be opened or created, an error is
// returned.
func (p *Persister) Setup() error {
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY|os.O_CREATE, 0600)
	defer f.Close()
	return err
}

// Get retrieves the value associated with the given key from the persisted
// data. The method loads the data from the persistence file in a way
// that is both safe for concurrency and locked against use by other
// programs using the lockedfile package (like go binary itself does).
// If the key is not present in the data, a nil value will be
// returned.
func (p *Persister) Get(key string) string {
	data := p.loadFile()
	return data[key]
}

// Set stores a key-value pair in the persisted data. The method loads
// the existing data from the persistence file, updates the data with
// the new key-value pair, and saves it back to the file.
func (p *Persister) Set(key, value string) {
	data := p.loadFile()
	data[key] = value
	p.saveFile(data)
}

func (p *Persister) loadFile() map[string]string {
	data := make(map[string]string)
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return data
	}
	json.NewDecoder(f).Decode(&data)
	return data
}

func (p *Persister) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return
	}
	json.NewEncoder(f).Encode(data)
}
