/*
Package json contains interface specifications for representing any Go
type as JSON where possible. Using the goprintasjson tool allows for quick code generation of scaffolding to make any Go type easily used as JSON.

For querying JSON/YAML/TOML/XML structured data the github.com/rwxrob/yq functions are recommended but not included due to the large number of dependencies that would forced into this project.
*/
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// specification (unlike the encoding/json standard which defaults to
// escaping many other characters as well unnecessarily).
func Escape(in string) string {
	out := ``
	for _, r := range in {
		switch r {
		case '\t':
			out += `\t`
		case '\b':
			out += `\b`
		case '\f':
			out += `\f`
		case '\n':
			out += `\n`
		case '\r':
			out += `\r`
		case '\\':
			out += `\\`
		case '"':
			out += `\"`
		default:
			out += string(r)
		}
	}
	return out
}

// Marshal mimics json.Marshal from the encoding/json package without
// the broken, unnecessary HTML escapes and extraneous newline that the
// json.Encoder adds. Call this from your own MarshalJSON methods to get
// JSON rendering that is more readable and compliant with the JSON
// specification (unless you are using the extremely rare case of
// dumping that into HTML, for some reason). Note that this cannot be
// called from any structs MarshalJSON method on itself because it will
// cause infinite functional recursion. Write a proper MarshalJSON
// method or create a dummy struct and call json.Marshal on that
// instead.
func Marshal(v any) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	return []byte(strings.TrimSpace(buf.String())), err
}

// MarshalIndent mimics json.Marshal from the encoding/json package but
// without the escapes, etc. See Marshal.
func MarshalIndent(v any, a, b string) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent(a, b)
	err := enc.Encode(v)
	return []byte(strings.TrimSpace(buf.String())), err
}

// Unmarshal mimics json.Unmarshal from the encoding/json package.
func Unmarshal(buf []byte, v any) error {
	return json.Unmarshal(buf, v)
}

// This encapsulates anything with the AsJSON interface from this package
// by simply assigning a new variable with that item as the only value
// in the structure:
//
//	something := []string{"some","thing"}
//	jsonified := json.This{something}
//	jsonified.Print()
type This struct{ This any }

// UnmarshalJSON implements AsJSON
func (s *This) UnmarshalJSON(buf []byte) error {
	return json.Unmarshal(buf, &s.This)
}

// JSON implements AsJSON.
func (s This) JSON() ([]byte, error) { return json.Marshal(s.This) }

// String implements AsJSON and logs any error.
func (s This) String() string {
	byt, err := s.JSON()
	if err != nil {
		log.Print(err)
	}
	return string(byt)
}

// Print implements AsJSON printing with fmt.Println (adding a line
// return).
func (s This) Print() { fmt.Println(s.String()) }

// Log implements AsJSON.
func (s This) Log() { log.Print(s.String()) }
