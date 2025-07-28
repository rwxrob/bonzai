package inprops

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/lockedfile"
	"github.com/rwxrob/bonzai/futil"
)

// Persister represents a simple key-value storage system using the
// Java Persister file format. Deviations from the standard Properties
// specification:
//
// - Comments re not supported and will result in an error if present.
// - Blank lines are not allowed and will also result in an error.
// - Every line must be a valid `key=value` pair; otherwise, ignored.
//
// # Supported features
//
// - Keys and values are trimmed of leading and trailing whitespace.
// - Special characters are escaped and unescaped automatically.
// - The file format is compatible with most Persister parsers.
//
// # Usage
//
//		  storage := &inprops.Persister{File: "data.properties"}
//		  storage.Setup()
//		  storage.Set("key", "value")
//	    value := storage.Get("key")
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
// It opens the file in read-only mode, creating it if it doesn't already exist,
// and applies secure file permissions (0600) to restrict access.
// The file is immediately closed after being verified or created.
// If the file cannot be opened or created, an error is returned.
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
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := unescapeValue(strings.TrimSpace(parts[1]))
			data[key] = value
		}
	}
	return data
}

func (p *Persister) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return
	}

	writer := bufio.NewWriter(f)
	for key, value := range data {
		escapedValue := escapeValue(value)
		writer.WriteString(key + "=" + escapedValue + "\n")
	}
	writer.Flush()
}

func escapeValue(value string) string {
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\\", "\\\\",
		"=", "\\=",
		":", "\\:",
	)
	return replacer.Replace(value)
}

func unescapeValue(value string) string {
	replacer := strings.NewReplacer(
		"\\n", "\n",
		"\\r", "\r",
		"\\t", "\t",
		"\\\\", "\\",
		"\\=", "=",
		"\\:", ":",
	)
	return replacer.Replace(value)
}
