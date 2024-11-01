package scanner

import (
	"fmt"
)

// Cursor contains a cursor pointer to a bytes slice buffer and
// information pointing to a specific location in the bytes buffer. The
// order of fields is guaranteed to never change.
type Cursor struct {
	Buf *[]byte // pointer to actual bytes buffer
	R   rune    // last rune scanned
	B   int     // beginning of last rune scanned
	E   int     // effective end of last rune scanned (beginning of next)
}

// String implements fmt.Stringer with the last rune scanned (R/Rune),
// and the beginning and ending byte positions joined with
// a dash.
func (c Cursor) String() string {
	return fmt.Sprintf("%q %v-%v", c.R, c.B, c.E)
}
