// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scan

import (
	"fmt"
	"log"

	"github.com/rwxrob/bonzai/scan/tk"
)

// Pos contains the user-land position for reporting back when there is
// a problem with parsing. This is different from internal position
// information used to parse from the internal data buffer.
type Pos struct {
	Line     int // lines (rows) starting at 1
	LineRune int // offset rune in line starting at 1
	LineByte int // offset byte in line starting at 1
	Rune     int // offset rune pos starting with 1
}

// NewLine forces the internal cursor position to increment it's
// user-land line counter and reset LineRune and LineByte to 1. This
// enables specific grammar parser implementations to explicitly create
// a new line as defined by that grammar.
func (p *Pos) NewLine() {
	p.Line++
	p.LineRune = 1
	p.LineByte = 1
}

// Cur is a cursor structure that points to specific position within
// buffered data. An internal Cur is used within the Parser. Cursors
// are returned by methods of Parser implementations.  Cursors must be
// set to EOD and have Len set to 0 if the Parser is asked to read
// beyond the end of the data.  Manipulating the values of a Curs
// directly is strongly discouraged (but not worth the performance hit
// of stopping it with interface encapsulation).
type Cur struct {
	Pos
	Rune rune // last rune decoded
	Byte int  // beginning of last rune decoded
	Len  int  // length of last rune decoded (0-4)
	Next int  // beginning of next rune to decode
}

// String fulfills the fmt.Stringer interface by printing
// the Cursor's location in a human-friendly way:
//
//   U+1F47F '👿' 1,3-5 (3-5)
//                | | |  | |
//             line | |  | overall byte offset
//   line rune offset |  overall rune offset
//     line byte offset
//
func (c *Cur) String() string {
	if c == nil {
		return "<nil>"
	}
	if c.Rune == tk.EOD {
		return "<EOD>"
	}
	s := fmt.Sprintf(`%U %q %v,%v-%v (%v-%v)`,
		c.Rune, c.Rune,
		c.Pos.Line, c.Pos.LineRune, c.Pos.LineByte,
		c.Pos.Rune, c.Byte+1,
	)
	return s
}

// Print prints the cursor itself in String form. See String.
func (c *Cur) Print() { fmt.Println(c.String()) }

// Log calls log.Println on the cursor itself in String form. See String.
func (c *Cur) Log() { log.Println(c.String()) }
