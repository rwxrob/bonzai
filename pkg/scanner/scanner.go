// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scanner

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"unicode/utf8"
)

// Trace sets the trace for everything that uses this package. Use
// TraceOn/Off for specific scanner tracing.
var Trace int

// S (to avoid stuttering) implements a buffered data, non-linear,
// rune-centric, scanner with regular expression support
type S struct {
	Buf     []byte   // full buffer for lookahead or behind
	R       rune     // last decoded/scanned rune, maybe >1byte
	B       int      // index pointing beginning of R
	E       int      // index pointing to end (after) R
	NewLine []string // []string{"\r\n","\n"} by default
	Trace   int      // non-zero activates tracing

	viewlen int // length of bytes to show in preview
	errors  []error
	maxerr  int
}

var ViewLenDefault = 10 // default length of preview window

// New is a high-level scanner constructor and initializer that takes
// a single optional argument containing any valid Buffer() argument.
// Invalid arguments will fail (not fatal) with log output.
func New(args ...any) *S {
	s := new(S)
	switch len(args) {
	case 2:
		if c, ok := args[1].(Cursor); ok {
			s.Goto(c)
		}
		fallthrough
	case 1:
		s.Buffer(args[0])
	}
	return s
}

func (s *S) SetViewLen(a int) { s.viewlen = a }
func (s *S) SetMaxErr(i int)  { s.maxerr = i }
func (s *S) Bytes() *[]byte   { return &s.Buf }
func (s *S) Rune() rune       { return s.R }
func (s *S) RuneB() int       { return s.B }
func (s *S) RuneE() int       { return s.E }
func (s *S) Mark() Cursor     { return Cursor{&s.Buf, s.R, s.B, s.E} }
func (s *S) Goto(c Cursor)    { s.R, s.B, s.E = c.R, c.B, c.E }
func (s *S) ViewLen() int     { return s.viewlen }
func (s *S) TraceOff()        { s.Trace = 0 }
func (s *S) TraceOn()         { s.Trace++ }

// CopyEE returns copy (n,m] fulfilling Scanner interface.
func (s *S) CopyEE(m Cursor) string {
	if m.B <= s.B {
		return string(s.Buf[m.E:s.E])
	}
	return string(s.Buf[s.E:m.E])
}

// CopyBB returns copy [n,m] fulfilling Scanner interface.
func (s *S) CopyBE(m Cursor) string {
	if m.B <= s.B {
		return string(s.Buf[m.B:s.E])
	}
	return string(s.Buf[s.B:m.E])
}

// CopyBB returns copy [n,m) fulfilling Scanner interface.
func (s *S) CopyBB(m Cursor) string {
	if m.B <= s.B {
		return string(s.Buf[m.B:s.B])
	}
	return string(s.Buf[s.B:m.B])
}

// CopyEB returns copy (n,m) fulfilling Scanner interface.
func (s *S) CopyEB(m Cursor) string {
	if m.B <= s.B {
		return string(s.Buf[m.E:s.B])
	}
	return string(s.Buf[s.E:m.B])
}

// Open opens the file at path and passes it to Buffer. Fulfills
// Scanner.
func (s *S) Open(path string) error {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return err
	}
	return s.Buffer(f)
}

// Buffer sets the internal bytes buffer (Buf) and resets the existing
// cursor values to their initial state (null, 0,0). This is useful when
// testing in order to buffer strings as well as content from any
// io.Reader, []byte, []rune, or string. Fulfills Scanner.
func (s *S) Buffer(b any) error {
	switch v := b.(type) {
	case string:
		s.Buf = []byte(v)
	case []byte:
		s.Buf = v
	case []rune:
		s.Buf = []byte(string(v))
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			return err
		}
		s.Buf = b
	}
	s.R = '\x00'
	s.B = 0
	s.E = 0
	return nil
}

// Revert is a shortcut for s.Goto always returning false.
func (s *S) Revert(m Cursor) bool {
	s.Goto(m)
	return false
}

// Position contains the human-friendly information about the position
// within a give text file. Note that all values begin with 1 and not
// 0.
type Position struct {
	Rune    rune // rune at this location
	BufByte int  // byte offset in file
	BufRune int  // rune offset in file
	Line    int  // line offset
	LByte   int  // line column byte offset
	LRune   int  // line column rune offset
}

// String fulfills the fmt.Stringer interface by printing
// the Position in a human-friendly way:
//
//	U+1F47F '👿' 1,3-5 (3-5)
//	             | | |  | |
//	          line | |  | overall byte offset
//	line rune offset |  overall rune offset
//	  line byte offset
func (p Position) String() string {
	s := fmt.Sprintf(`%U %q %v,%v-%v (%v-%v)`,
		p.Rune, p.Rune,
		p.Line, p.LRune, p.LByte,
		p.BufRune, p.BufByte,
	)
	return s
}

// Print prints the cursor itself in String form. See String.
func (p Position) Print() { fmt.Println(p.String()) }

// Log calls log.Println on the cursor itself in String form. See String.
func (p Position) Log() { log.Println(p.String()) }

// Pos returns a human-friendly Position for the current location.
// When multiple positions are needed use Positions instead.

func (s S) Pos() Position { return s.Positions(s.E)[0] }

// Positions returns human-friendly Position information (which can easily
// be used to populate a text/template) for each raw byte offset (s.E).
// Only one pass through the buffer (s.Buf) is required to count lines and
// runes since the raw byte position (s.E) is frequently changed
// directly.  Therefore, when multiple positions are wanted, consider
// caching the raw byte positions (s.E) and calling Positions() once for
// all of them.
func (s S) Positions(p ...int) []Position {
	pos := make([]Position, len(p))

	if len(p) == 0 {
		return pos
	}

	if s.NewLine == nil {
		s.NewLine = []string{"\r\n", "\n"}
	}

	_rune, line, lbyte, lrune := 1, 1, 1, 1
	_s := S{Buf: s.Buf}
	//_s.Trace++

	for _s.Scan() {

		for _, nl := range s.NewLine {
			if _s.Is(nl) {
				line++
				_s.E += len(nl) - 1
				_rune += len(nl) - 1
				lbyte = 0
				lrune = 0
				continue
			}
		}

		for i, v := range p {
			if _s.E == v {
				pos[i] = Position{
					Rune:    _s.R,
					BufByte: _s.E,
					BufRune: _rune,
					Line:    line,
					LByte:   lbyte,
					LRune:   lrune,
				}
			}
		}

		rlen := len([]byte(string(s.R)))
		lbyte += rlen
		lrune++
		_rune++

	}

	return pos
}

// String implements fmt.Stringer with simply the position (E) and
// quoted rune (R) along with its Unicode. For printing more human
// friendly information about the current scanner state use Report.
func (s S) String() string {
	if s.viewlen == 0 {
		s.viewlen = ViewLenDefault
	}
	end := s.E + s.viewlen
	if end > len(s.Buf) {
		end = len(s.Buf)
	}
	return fmt.Sprintf("%v %q",
		Cursor{&s.Buf, s.R, s.B, s.E}, s.Buf[s.E:end])
}

// Print is shorthand for fmt.Println(s).
func (s S) Print() { fmt.Println(s) }

// Log is shorthand for log.Print(s).
func (s S) Log() { log.Println(s) }

// Scan decodes the next rune, setting it to R, and advances position
// (P) by the size of the rune (R) in bytes returning false then there
// is nothing left to scan. Only runes bigger than utf8.RuneSelf are
// decoded since most runes (ASCII) will usually be under this number.
func (s *S) Scan() bool {

	if s.E >= len(s.Buf) {
		return false
	}

	ln := 1
	r := rune(s.Buf[s.E])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.Buf[s.E:])
		if ln == 0 {
			return false
		}
	}

	s.B = s.E
	s.E += ln
	s.R = r

	if s.Trace > 0 || Trace > 0 {
		s.Log()
	}

	return true
}

// Peek returns true if the passed string matches from current position
// in the buffer (s.B) forward. Returns false if the string
// would go beyond the length of buffer (len(s.Buf)). Peek does not
// advance the Scanner.
func (s *S) Peek(a string) bool {
	if len(a)+s.E > len(s.Buf) {
		return false
	}
	if string(s.Buf[s.E:s.E+len(a)]) == a {
		return true
	}
	return false
}

// Finished returns true if scanner has nothing more to scan.
func (s *S) Finished() bool { return s.E == len(s.Buf) }

// Beginning returns true if and only if the scanner is currently
// pointing to the beginning of the buffer without anything scanned at
// all.
func (s *S) Beginning() bool { return s.E == 0 }

// Is returns true if the passed string matches the last scanned rune
// and the runes ahead matching the length of the string.  Returns false
// if the string would go beyond the length of buffer (len(s.Buf)).
func (s *S) Is(a string) bool {

	if len(a)+s.B > len(s.Buf) {
		return false
	}

	if string(s.Buf[s.B:s.B+len(a)]) == a {
		return true
	}

	return false
}

// PeekMatch checks for a regular expression match at the current
// position in the buffer providing a mechanism for positive and
// negative lookahead expressions. It returns the length of the match.
// Successful matches might be zero (see regexp.Regexp.FindIndex).
// A negative value is returned if no match is found. Note that Go
// regular expressions now include the Unicode character classes (ex:
// \p{L|d}) that should be used over dated alternatives (ex: \w).
func (s *S) PeekMatch(re *regexp.Regexp) int {

	loc := re.FindIndex(s.Buf[s.E:])
	if loc == nil {
		return -1
	}

	if loc[0] == 0 {
		return loc[1]
	}

	return -1

}

// Match checks for a regular expression match at the last position in
// the buffer (s.B) providing a mechanism for positive and negative
// lookahead expressions. It returns the length of the match.
// Successful matches might be zero (see regexp.Regexp.FindIndex).
// A negative value is returned if no match is found.  Note that Go
// regular expressions now include the Unicode character classes (ex:
// \p{L|d}) that should be used over dated alternatives (ex: \w).
func (s *S) Match(re *regexp.Regexp) int {
	loc := re.FindIndex(s.Buf[s.B:])
	if loc == nil {
		return -1
	}
	if loc[0] == 0 {
		return loc[1]
	}
	return -1
}
