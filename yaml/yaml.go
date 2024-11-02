package yaml

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

// AsYAML specifies a type that must support marshaling using the
// gopkg.in/yaml.v3 package with its defaults for marshaling and
// unmarshaling. The yaml.v3 package supports `yaml:",inline"` and is
// fully compatible for encoding and decoding from encoding/json as well
// (yaml.v2 is not).
//
// String is from fmt.Stringer, but fulfilling this interface in this
// package promises to render the string specifically using yaml.v3
// default output marshaling --- especially when it comes to consistent
// indentation and wrapping. While YAML is a flexible format,
// consistency ensures the most efficient and sustainable creation of
// tests and other systems that require such consistency, whether or not
// dependency on such consistency is a "good idea" (as demonstrated by
// the Kubernetes project).
//
// Printer specifies methods for printing self as YAML and will log any
// error if encountered. Printer provides a consistent representation of
// any structure such that it an easily be read and compared as YAML
// whenever printed and test. Sadly, the default string representations
// for most types in Go are virtually unusable for consistent
// representations of any structure. And while it is true that YAML data
// should be supported in any way that is it presented, some consistent
// output makes for more consistent debugging, documentation, and
// testing.
//
// AsYAML implementations must Print and Log the output of String from
// the same interface.
//
// MarshalYAML and UnmarshalYAML must be explicitly defined and use the
// gopkg.in/yaml.v3 to avoid confusion. Use of the helper yaml.This struct
// may facilitate this for existing types that do not wish to implement
// the full interface.
type AsYAML interface {
	YAML() ([]byte, error)
	String() string
	Print()
	Log() string
	MarshalYAML() ([]byte, error)
	UnmarshalYAML(buf []byte) error
}

// This encapsulates anything with the AsYAML interface from this package
// by simply assigning a new variable with that item as the only value
// in the structure:
//
//	something := []string{"some","thing"}
//	yamlified := yaml.This{something}
//	yamlified.Print()
type This struct{ This any }

// UnmarshalYAML implements AsYAML
func (s *This) UnmarshalYAML(buf []byte) error {
	return yaml.Unmarshal(buf, &s.This)
}

// YAML implements AsYAML.
func (s This) YAML() ([]byte, error) { return yaml.Marshal(s.This) }

// String implements AsYAML and logs any error.
func (s This) String() string {
	byt, err := s.YAML()
	if err != nil {
		log.Print(err)
	}
	return string(byt)
}

// Print implements AsYAML printing with fmt.Print (no additional line).
func (s This) Print() { fmt.Print(s.String()) }

// Log implements AsYAML.
func (s This) Log() { log.Print(s.String()) }
