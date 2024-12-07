package persisters

import (
	"encoding/json"
	"os"

	"github.com/rogpeppe/go-internal/lockedfile"
)

type JSON struct {
	File string // Exported field for direct access
}

func (p *JSON) Setup() error {
	// Ensure the file exists but do not pre-load it into memory
	_, err := lockedfile.OpenFile(p.File, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (p *JSON) Get(key string) string {
	// Load the latest state of the file
	data := p.loadFile()
	return data[key] // Returns "" if key is not present
}

func (p *JSON) Set(key, value string) {
	// Load the latest state of the file
	data := p.loadFile()
	data[key] = value // Modify the in-memory state
	p.saveFile(data)  // Save the updated state
}

func (p *JSON) loadFile() map[string]string {
	data := make(map[string]string)
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY, 0600)
	if err != nil {
		// Handle file read error
		return data
	}
	defer f.Close()

	_ = json.NewDecoder(f).Decode(&data) // Ignore errors for empty files
	return data
}

func (p *JSON) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		// Handle file write error
		return
	}
	defer f.Close()

	_ = json.NewEncoder(f).Encode(data) // Ignore encoding errors
}
