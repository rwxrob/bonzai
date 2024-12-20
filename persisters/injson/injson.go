package injson

import (
	"encoding/json"
	"os"

	"github.com/rogpeppe/go-internal/lockedfile"
)

type Persister struct {
	File string
}

// Setup ensures that the persistence file exists and is ready for use.
// It opens the file in read-only mode, creating it if it doesn't
// already exist, and applies secure file permissions (0600) to restrict
// access. The file is immediately closed after being verified or
// created. If the file cannot be opened or created, an error is
// returned.
func (p *Persister) Setup() error {
	_, err := lockedfile.OpenFile(p.File, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return data
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&data)
	return data
}

func (p *Persister) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(data)
}
