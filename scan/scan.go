// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan implements a non-linear, rune-centric, buffered data
scanner that includes its own high-level syntax comprised of scannable
structures from the "is" subpackage making parser generation (by hand or
code generation) trivial from any structured meta languages such as
PEGN, PEG, EBNF, ABNF, etc. Most will use the scanner to create parsers
quickly by hand where a regular expression will not suffice. See the
"is" and "tk" packages for a growing number of common, centrally
maintain scannables for your parsing pleasure. Also see the "mark"
(BonzaiMark) subpackage for a working example of the scanner in action,
which is used by the included help.Cmd command.
*/
package scan

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/rwxrob/bonzai/scan/is"
	"github.com/rwxrob/bonzai/scan/tk"
)

// Scanner implements a non-linear, rune-centric, buffered data scanner.
// See New for creating a usable struct that implements Scanner. The
// buffer and cursor are directly exposed to facilitate
// higher-performance, direct access when needed.
type Scanner struct {

	// Buf is the data buffer providing infinite look-ahead and behind.
	Buf []byte

	// Cur is the active current cursor pointing to the Buf data.
	Cur *Cur

	// Snapped contains the latest Cur when Snap was called.
	Snapped *Cur

	// ExtendExpect provides a hook to support additional custom
	// scannables for both Expect and Check Scanner methods. Take note of
	// the ErrorExpected errors in order to construct similar errors where
	// returning ErrorExpected itself would not provide clear error
	// messages.
	ExtendExpect func(s *Scanner, scannable ...any) (*Cur, error)
}

// New returns a newly initialized non-linear, rune-centric, buffered
// data scanner with support for parsing data from io.Reader, string,
// and []byte types. Returns nil and the error if any encountered during
// initialization. Also see the Init method.
func New(i any) (*Scanner, error) {
	s := new(Scanner)
	if err := s.Init(i); err != nil {
		return nil, err
	}
	return s, nil
}

// Init reads all of passed parsable data (io.Reader, string, []byte)
// into buffered memory, scans the first rune, and sets the internals of
// scanner appropriately returning an error if anything happens while
// attempting to read and buffer the data (OOM, etc.).
func (s *Scanner) Init(i any) error {
	if err := s.buffer(i); err != nil {
		return err
	}
	r, ln := utf8.DecodeRune(s.Buf) // scan first
	if ln == 0 {
		r = tk.EOD
		return fmt.Errorf("scanner: failed to scan first rune")
	}
	s.Cur = new(Cur)
	s.Cur.Rune = r
	s.Cur.Len = ln
	s.Cur.Next = ln
	s.Cur.Pos.Line = 1
	s.Cur.Pos.LineRune = 1
	s.Cur.Pos.LineByte = 1
	s.Cur.Pos.Rune = 1
	return nil
}

// reads and buffers io.Reader, string, or []byte types
func (s *Scanner) buffer(i any) error {
	var err error
	switch in := i.(type) {
	case io.Reader:
		s.Buf, err = io.ReadAll(in)
		if err != nil {
			return err
		}
	case string:
		s.Buf = []byte(in)
	case []byte:
		s.Buf = in
	default:
		return fmt.Errorf("scanner: unsupported input type: %t", i)
	}
	if len(s.Buf) == 0 {
		return fmt.Errorf("scanner: no input")
	}
	return err
}

// Scan decodes the next rune and advances the scanner cursor by one.
func (s *Scanner) Scan() {
	if s.Done() {
		return
	}
	r, ln := utf8.DecodeRune(s.Buf[s.Cur.Next:])
	if ln != 0 {
		s.Cur.Byte = s.Cur.Next
		s.Cur.Pos.LineByte += s.Cur.Len
	} else {
		r = tk.EOD
	}
	s.Cur.Rune = r
	s.Cur.Pos.Rune += 1
	s.Cur.Next += ln
	s.Cur.Pos.LineRune += 1
	s.Cur.Len = ln
}

// ScanN scans the next n runes advancing n runes forward or returns
// s.Done() if attempted after already at end of data.
func (s *Scanner) ScanN(n int) {
	for i := 0; i < n; i++ {
		s.Scan()
	}
}

// Done returns true if current cursor rune is tk.EOD and the cursor length
// is also zero.
func (s *Scanner) Done() bool {
	return s.Cur.Rune == tk.EOD && s.Cur.Len == 0
}

// String delegates to internal cursor String.
func (s *Scanner) String() string { return s.Cur.String() }

// Print delegates to internal cursor Print.
func (s *Scanner) Print() { s.Cur.Print() }

// Mark returns a copy of the current scanner cursor to preserve like
// a bookmark into the buffer data. See Cur, Look, LookSlice.
func (s *Scanner) Mark() *Cur {
	if s.Cur == nil {
		return nil
	}
	// force a copy
	cp := *s.Cur
	return &cp
}

// Snap sets an extra internal cursor to the current cursor. See Mark.
func (s *Scanner) Snap() { s.Snapped = s.Mark() }

// Back jumps the current cursor to the last Snap (Snapped).
func (s *Scanner) Back() { s.Jump(s.Snapped) }

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data. Beware, however, that the new cursor must originate
// from the same (or identical) data buffer or the values will be out of
// sync.
func (s *Scanner) Jump(c *Cur) { nc := *c; s.Cur = &nc }

// Peek returns a string containing all the runes from the current
// scanner cursor position forward to the number of runes passed.
// If end of data is encountered it will return everything up until that
// point.  Also see Look and LookSlice.
func (s *Scanner) Peek(n uint) string {
	buf := ""
	pos := s.Cur.Byte
	for c := uint(0); c < n; c++ {
		r, ln := utf8.DecodeRune(s.Buf[pos:])
		if ln == 0 {
			break
		}
		buf += string(r)
		pos += ln
	}
	return buf
}

// Look returns a string containing all the bytes from the current
// scanner cursor position ahead or behind to the passed cursor
// position. Neither the internal nor the passed cursor position is
// changed. Also see Peek and LookSlice.
func (s *Scanner) Look(to *Cur) string {
	if to.Byte < s.Cur.Byte {
		return string(s.Buf[to.Byte:s.Cur.Next])
	}
	return string(s.Buf[s.Cur.Byte:to.Next])
}

// LookSlice returns a string containing all the bytes from the first
// cursor up to the second cursor. Neither cursor position is changed.
func (s *Scanner) LookSlice(beg *Cur, end *Cur) string {
	return string(s.Buf[beg.Byte:end.Next])
}

// Expect takes a variable list of parsable types including the
// following (in order of priority):
//
//     string          - "foo" simple string
//     rune            - 'f' uint32, but "rune" in errors
//     is.Not{any...}  - negative look-ahead set (slice)
//     is.In{any...}   - one positive look-ahead from set (slice)
//     is.Seq{any...}  - required positive look-ahead sequence (slice)
//     is.Opt{any...}  - optional positive look-ahead set (slice)
//     is.Min{n,any}   - minimum positive look-aheads
//     is.MMx{n,m,any} - minimum and maximum positive look-aheads
//     is.X{n,any}     - exactly n positive look-aheads
//     is.Rng{n,m}     - inclusive range from rune n to rune m (n,m)
//
// Any token string from the tk subpackage can be used as well (and will
// be noted as special in error out as opposed to simple strings. This
// allows for very readable functional grammar parsers to be created
// quickly without exceptional overhead from additional function calls
// and indirection. As some have said, "it's regex without the regex."
func (s *Scanner) Expect(scannables ...any) (*Cur, error) {
	var beg, end *Cur
	beg = s.Cur

	for _, m := range scannables {

		// please keep the most common at the top
		switch v := m.(type) {

		case string: // ----------------------------------------------------
			if v == "" {
				return nil, fmt.Errorf("expect: cannot parse empty string")
			}
			for _, r := range []rune(v) {
				if s.Cur.Rune != r {
					err := s.ErrorExpected(r)
					s.Jump(beg)
					return nil, err
				}
				end = s.Mark()
				s.Scan()
			}

		case rune: // ------------------------------------------------------
			if s.Cur.Rune != v {
				err := s.ErrorExpected(m)
				s.Jump(beg)
				return nil, err
			}
			end = s.Mark()
			s.Scan()

		case is.Not: // ----------------------------------------------------
			for _, i := range v {
				if _, e := s.Check(i); e == nil {
					err := s.ErrorExpected(v, i)
					s.Jump(beg)
					return nil, err
				}
			}
			end = s.Mark()

		case is.In: // -----------------------------------------------------
			var m *Cur
			for _, i := range v {
				var err error
				last := s.Mark()
				m, err = s.Expect(i)
				if err == nil {
					break
				}
				s.Jump(last)
			}
			if m == nil {
				return nil, s.ErrorExpected(v)
			}
			end = m

		case is.Seq: // ----------------------------------------------------
			m, err := s.Expect(v...)
			if err != nil {
				s.Jump(beg)
				return nil, err
			}
			end = m

		case is.Opt: // ----------------------------------------------------
			var m *Cur
			for _, i := range v {
				var err error
				m, err = s.Expect(is.MMx{0, 1, i})
				if err != nil {
					s.Jump(beg)
					return nil, s.ErrorExpected(v)
				}
			}
			end = m

		case is.MMx: // ----------------------------------------------------
			c := 0
			last := s.Mark()
			var err error
			var m *Cur
			for {
				m, err = s.Expect(v.This)
				if err != nil {
					break
				}
				last = m
				c++
			}
			if c == 0 && v.Min == 0 {
				if end == nil {
					end = last
				}
				continue
			}
			if !(v.Min <= c && c <= v.Max) {
				s.Jump(last)
				return nil, s.ErrorExpected(v)
			}
			end = last

		case is.Min: // ----------------------------------------------------
			c := 0
			last := s.Mark()
			var err error
			var m *Cur
			for {
				m, err = s.Expect(v.This)
				if err != nil {
					break
				}
				last = m
				c++
			}
			if c < v.Min {
				s.Jump(beg)
				return nil, s.ErrorExpected(v)
			}
			end = last

		case is.X: // ------------------------------------------------------
			m, err := s.Expect(is.MMx{v.X, v.X, v.This})
			if err != nil {
				s.Jump(beg)
				return nil, s.ErrorExpected(v)
			}
			end = m

		case is.Rng: // ----------------------------------------------------
			if !(v.First <= s.Cur.Rune && s.Cur.Rune <= v.Last) {
				err := s.ErrorExpected(v)
				s.Jump(beg)
				return nil, err
			}
			end = s.Mark()
			s.Scan()

		default: // --------------------------------------------------------
			if s.ExtendExpect != nil {
				return s.ExtendExpect(s, scannables...)
			}
			return nil, fmt.Errorf("expect: unscannable type (%T)", m)
		}
	}
	return end, nil
}

// ErrorExpected returns a verbose, one-line error describing what was
// expected when it encountered whatever the scanner last scanned. All
// scannable types are supported. See Expect.
func (s *Scanner) ErrorExpected(this any, args ...any) error {
	var msg string
	but := fmt.Sprintf(` at %v`, s)
	if s.Done() {
		runes := `runes`
		if s.Cur.Pos.Rune == 1 {
			runes = `rune`
		}
		but = fmt.Sprintf(`, exceeded data length (%v %v)`,
			s.Cur.Pos.Rune, runes)
	}
	switch v := this.(type) {
	case rune: // otherwise will use uint32
		msg = fmt.Sprintf(`expected rune %q`, v)
	case is.Not:
		msg = fmt.Sprintf(`unexpected %q`, args[0])
	case is.In:
		str := `expected one of %q`
		msg = fmt.Sprintf(str, v)
	case is.Seq:
		str := `expected %q in sequence %q`
		msg = fmt.Sprintf(str, args[0], v)
	case is.Opt:
		str := `expected an optional %v`
		msg = fmt.Sprintf(str, v)
	case is.Min:
		str := `expected min %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.This)
	case is.MMx:
		str := `expected min %v, max %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.Max, v.This)
	case is.X:
		str := `expected exactly %v of %q`
		msg = fmt.Sprintf(str, v.X, v.This)
	case is.Rng:
		str := `expected range [%v-%v]`
		msg = fmt.Sprintf(str, string(v.First), string(v.Last))
	default:
		msg = fmt.Sprintf(`expected %T %q`, v, v)
	}
	return errors.New(msg + but)
}

// NewLine delegates to interval Curs.NewLine.
func (s *Scanner) NewLine() { s.Cur.NewLine() }

// Check behaves exactly like Expect but jumps back to the original
// cursor position after scanning for expected scannable values.
func (s *Scanner) Check(scannables ...any) (*Cur, error) {
	defer s.Jump(s.Mark())
	return s.Expect(scannables...)
}
