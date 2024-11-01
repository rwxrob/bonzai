package futil

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/uniq"
)

// Multipart is meant to contain the delimited sections of output and can
// be marshalled into a single delimited string safely and automatically
// simply by using it in a string context.
type Multipart struct {
	Delimiter string
	Map       map[string]string
}

// MarshalText fulfills the encoding.TextMarshaler interface by
// delimiting each section of output with a unique delimiter line that
// contains a space and the key for each section. Order of sections is
// indeterminate officially (but consistent for testing, per Go). The
// special "break" delimiter is always the last line. The Delimiter is
// used if defined, otherwise one is automatically assigned.
func (o Multipart) MarshalText() ([]byte, error) {
	var out string
	if o.Delimiter == "" {
		o.Delimiter = uniq.Base32()
	}
	if o.Delimiter == "" {
		return nil, fmt.Errorf(`unable to get random data`)
	}
	for k, v := range o.Map {
		out += o.Delimiter + " " + k + "\n" + v + "\n"
	}
	out += o.Delimiter + " break"
	return []byte(out), nil
}

// UnmarshalText fulfills the encoding.TextUnmarshaler interface by
// using its internal Delimiter or sensing the delimiter as the first
// text field (up to the first space) if not set and using that
// delimiter to parse the remaining data into the key/value pairs ending
// when either the end of text is encountered or the special "break"
// delimiter is read.
func (o *Multipart) UnmarshalText(text []byte) error {

	var cur string
	s := bufio.NewScanner(bytes.NewReader(text))
	o.Map = map[string]string{}

	// detect delimiter as first line if Delimiter is unset
	if o.Delimiter == "" {

		if !s.Scan() {
			return fmt.Errorf(`failed to scan first line`)
		}

		f := strings.Fields(s.Text())
		if len(f) < 2 {
			return fmt.Errorf(`first line is not delimiter`)
		}

		if f[1] == "break" {
			return nil
		}

		o.Delimiter = f[0]
		cur = f[1]
	}

	for s.Scan() {
		line := s.Text()

		// delimiter?
		if strings.HasPrefix(line, o.Delimiter) {
			f := strings.Fields(line)
			if len(f) < 2 {
				return fmt.Errorf(`delimiter missing key`)
			}
			if cur != "" {
				o.Map[cur] = o.Map[cur][:len(o.Map[cur])-1]
			}
			if f[1] == `break` {
				return nil
			}
			cur = f[1]
			continue
		}

		if cur == "" {
			continue
		}

		o.Map[cur] += line + "\n"

	}

	return nil
}

// String fulfills the fmt.Stringer interface by calling MarshalText.
func (o Multipart) String() string {
	buf, err := o.MarshalText()
	if err != nil {
		return ""
	}
	return string(buf)
}
