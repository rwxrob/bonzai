package tinout

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Spec represents a test specification of inputs and outputs.
type Spec struct {
	Name     string    `yaml:"Name,omitempty"`
	Version  string    `yaml:"Version,omitempty"`
	Source   string    `yaml:"Source,omitempty"`
	Issues   string    `yaml:"Issues,omitempty"`
	Discuss  string    `yaml:"Discuss,omitempty"`
	Notes    string    `yaml:"Notes,omitempty"`
	Date     string    `yaml:"Date,omitempty"`
	License  string    `yaml:"License,omitempty"`
	Tests    []Test    `yaml:"Tests,omitempty"`
	Sections []Section `yaml:"Sections,omitempty"`
}

// Section of tests.
type Section struct {
	Name  string `yaml:"Name,omitempty"`
	Notes string `yaml:"Notes,omitempty"`
	Tests []Test `yaml:"Tests,omitempty"`
}

// Test has the input, output, and notes for a given test.
type Test struct {
	I   string `yaml:"I,omitempty"`   // input
	O   string `yaml:"O,omitempty"`   // output
	N   string `yaml:"N,omitempty"`   // notes
	Got string `yaml:"Got,omitempty"` // result of last check
}

// Passing returns if t.Got is equal to t.O.
func (t *Test) Passing() bool {
	return t.Got == t.O
}

// CheckMethod is any function that takes a test and returns true if the test
// passed. The value it got when testing can be stored in t.Got and t.Passed
// should be set to true if passed.
type CheckMethod func(t *Test) bool

// State returns a string describing the current state of the test.
func (t *Test) State() string {
	passing := "failing"
	if t.Passing() {
		passing = "passing"
	}
	return fmt.Sprintf("\nState:    %q\nInput:    %q\nWanted:   %q\nGot:      %q\n", passing, t.I, t.O, t.Got)
}

// Load loads the [Spec] from a YAML file at path.
func Load(path string) (Spec, error) {
	s := Spec{}
	byt, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	err = yaml.Unmarshal(byt, &s)
	return s, err
}

// Read reads the Spec from a YAML stream reader.
func Read(r io.Reader) (Spec, error) {
	s := Spec{}
	byt, err := io.ReadAll(r)
	if err != nil {
		return s, err
	}
	err = yaml.Unmarshal(byt, &s)
	return s, err
}

// Check takes a function as an argument that takes in [CheckMethod]
// function. It then calls the CheckMethod on all the s.Tests and
// s.Sections[].Tests until it finds the first one that does not pass
// and returns a pointer to it. Returns nil if all pass.
func (s *Spec) Check(ok CheckMethod) *Test {
	for _, t := range s.Tests {
		if !ok(&t) {
			return &t
		}
	}
	for _, sc := range s.Sections {
		for _, t := range sc.Tests {
			if !ok(&t) {
				return &t
			}
		}
	}
	return nil
}
