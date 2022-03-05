package scan

import (
	"encoding/json"
	"fmt"
	"log"
)

// Node structs are for use with bonzai/scan.R and compliment those from
// bonzai/tree. They are simply lighter and optimized for quick parsing
// during the scan process with scan.R.Parse and the z.P Expect
// expression.
type Node struct {
	T string  `json:",omitempty"` // type
	V string  `json:",omitempty"` // type
	U []*Node `json:",omitempty"` // type
}

// ---------------------------- PrintAsJSON ---------------------------

// JSON implements PrintAsJSON multi-line, 2-space indent JSON output.
func (s *Node) JSON() string { b, _ := json.Marshal(s); return string(b) }

// String implements PrintAsJSON and fmt.Stringer interface as JSON.
func (s Node) String() string { return s.JSON() }

// Print implements PrintAsJSON.
func (s *Node) Print() { fmt.Println(s.JSON()) }

// Log implements PrintAsJSON.
func (s Node) Log() { log.Print(s.JSON()) }
