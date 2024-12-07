package persisters

import (
	"bufio"
	"os"
	"strings"

	"github.com/rogpeppe/go-internal/lockedfile"
)

// Properties represents a simple key-value storage system using the
// Java Properties file format. Deviations from the standard Properties
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
// - The file format is compatible with most Properties parsers.
//
// # Usage
//
//		  storage := &persisters.Properties{File: "data.properties"}
//		  storage.Setup()
//		  storage.Set("key", "value")
//	    value := storage.Get("key")
type Properties struct {
	File string // File is the path to the properties file used for storage.
}

func (p *Properties) Setup() error {
	_, err := lockedfile.OpenFile(p.File, os.O_RDONLY|os.O_CREATE, 0600)
	return err // Ensure the file exists
}

func (p *Properties) Get(key string) string {
	data := p.loadFile()
	return data[key] // Returns "" if the key is not present
}

func (p *Properties) Set(key, value string) {
	data := p.loadFile()
	data[key] = value // Update the in-memory state
	p.saveFile(data)  // Save changes back to the file
}

func (p *Properties) loadFile() map[string]string {
	data := make(map[string]string)
	f, err := lockedfile.OpenFile(p.File, os.O_RDONLY, 0600)
	if err != nil {
		return data
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Parse key=value pairs without skipping blank or commented lines
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := unescapeValue(strings.TrimSpace(parts[1]))
			data[key] = value
		}
	}
	return data
}

func (p *Properties) saveFile(data map[string]string) {
	f, err := lockedfile.OpenFile(p.File, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for key, value := range data {
		escapedValue := escapeValue(value)
		_, _ = writer.WriteString(key + "=" + escapedValue + "\n")
	}
	writer.Flush()
}

// Escapes special characters in properties values
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

// Unescapes special characters in properties values
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
