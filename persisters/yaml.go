package persisters

import (
	"io"
	"os"

	"github.com/rogpeppe/go-internal/lockedfile"
	"gopkg.in/yaml.v3"
)

// YAML represents a key-value storage system using the YAML format.
//
// Features:
//
// - Values are stored as strings and loaded directly into a flat key-value map.
// - The file format is compatible with standard YAML parsers.
//
// # Usage
//
//		     storage := &persisters.YAML{File: "data.yaml"}
//	    	 storage.Setup()
//		     storage.Set("key", "value")
//		     value := storage.Get("key")
type YAML struct {
	File string // File is the path to the YAML file used for storage.
}

func (p *YAML) Setup() error {
	// Ensure the file exists with secure permissions
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func (p *YAML) Get(key string) string {
	data := p.loadFile()
	return data[key] // Returns "" if the key is not present
}

func (p *YAML) Set(key, value string) {
	data := p.loadFile()
	data[key] = value // Update the in-memory state
	p.saveFile(data)  // Save changes back to the file
}

func (p *YAML) loadFile() map[string]string {
	data := make(map[string]string)
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY, 0600)
	if err != nil {
		return data
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil || len(content) == 0 {
		return data
	}

	_ = yaml.Unmarshal(content, &data) // Ignore errors for simplicity
	return data
}

func (p *YAML) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	content, _ := yaml.Marshal(data) // Ignore errors for simplicity
	f.Write(content)
}
