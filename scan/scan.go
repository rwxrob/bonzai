// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan implements a non-linear, rune-centric, buffered data,
extendable scanner that includes its own high-level parsing expression
grammar (BPEGN) comprised of Go slices and structs used as expressions
from the "is" subpackage making parser generation (by hand or code)
trivial from any structured meta languages such as PEGN, PEG, EBNF,
ABNF, etc. Most will use the scanner to create parsers quickly where
a regular expression will not suffice. See the "is" and "tk" packages
for a growing number of common, centrally maintain expressions for your
parsing pleasure. Also see the "mark" (BonzaiMark) subpackage for
a working examples of the scanner in action, which is used by the
included Bonzai help.Cmd command and others.
*/
package scan

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/rwxrob/bonzai/scan/is"
	"github.com/rwxrob/bonzai/scan/tk"
	"github.com/rwxrob/bonzai/util"
)

// Hook is a function expression that accepts a reference to the current
// scanner and simply returns true or false. Hook functions are allowed
// to do whatever they need and must advance the scan.R themselves (if
// necessary) and should not be abused and are given the lowest priority
// when searching for expressions. Static scanning expressions will
// usually be faster than any Hook. Hook allows PEGN (and others) to
// indicate Hook names for executable code that must be run during the
// scanning of a specific grammar (indicated as "rhetorical" in some
// grammars). In fact, scan.Rs can be converted into parsers relatively
// easily simply by implementing a set of Hook functions to capture or
// render scanned data at specific points during the scan process. Since
// only the name of the Hook function is required BPEGN remains
// compatible with PEGN one-for-one transpiling.
type Hook func(s *R) bool

// R (as in scan.R or "scanner") implements a non-linear, rune-centric,
// buffered data scanner and provides full support for BPEGN. See New
// for creating a usable struct that implements scan.R. The buffer and
// cursor are directly exposed to facilitate higher-performance, direct
// access when needed.
type R struct {

	// Buf is the data buffer providing infinite look-ahead and behind.
	Buf    []byte
	BufLen int

	// Cur is the active current cursor pointing to the Buf data.
	Cur *Cur

	// Snapped contains the latest Cur when Snap was called.
	Snapped *Cur
}

// New returns a newly initialized non-linear, rune-centric, buffered
// data scanner with support for parsing data from io.Reader, string,
// and []byte types. Returns nil and the error if any encountered during
// initialization. Also see the Init method.
func New(i any) (*R, error) {
	s := new(R)
	if err := s.Init(i); err != nil {
		return nil, err
	}
	return s, nil
}

// Init reads all of passed parsable data (io.Reader, string, []byte)
// into buffered memory, scans the first rune, and sets the internals of
// scanner appropriately returning an error if anything happens while
// attempting to read and buffer the data (OOM, etc.).
func (s *R) Init(i any) error {
	s.Cur = new(Cur)
	s.Cur.Pos = Pos{}
	s.Cur.Pos.Line = 1
	s.Cur.Pos.LineRune = 1
	s.Cur.Pos.LineByte = 1
	s.Cur.Pos.Rune = 1

	if err := s.buffer(i); err != nil {
		return err
	}

	r, ln := utf8.DecodeRune(s.Buf) // scan first
	if ln == 0 {
		r = tk.EOD
		return fmt.Errorf("scanner: failed to scan first rune")
	}

	s.Cur.Rune = r
	s.Cur.Len = ln
	s.Cur.Next = ln

	return nil
}

// reads and buffers io.Reader, string, or []byte types
func (s *R) buffer(i any) error {
	var err error
	switch v := i.(type) {
	case io.Reader:
		s.Buf, err = io.ReadAll(v)
		if err != nil {
			return err
		}
	case string:
		s.Buf = []byte(v)
	case []byte:
		s.Buf = v
	default:
		return fmt.Errorf("scanner: unsupported input type: %t", i)
	}
	s.BufLen = len(s.Buf)
	if s.BufLen == 0 {
		return fmt.Errorf("scanner: no input")
	}
	return err
}

// Scan decodes the next rune and advances the scanner cursor by one.
// The method of scanning isn't as optimized as other scanner (for
// example, the scanner from the bonzai/json package), but it is
// sufficient for most high level needs.
func (s *R) Scan() {

	if s.Cur.Next == s.BufLen {
		s.Cur.Rune = tk.EOD
		return
	}

	ln := 1
	r := rune(s.Buf[s.Cur.Next])
	if r > utf8.RuneSelf {
		r, ln = utf8.DecodeRune(s.Buf[s.Cur.Next:])
	}
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
// EOD if attempted after already at end of data.
func (s *R) ScanN(n int) {
	for i := 0; i < n; i++ {
		s.Scan()
	}
}

// String delegates to internal cursor String.
func (s *R) String() string { return s.Cur.String() }

// Print delegates to internal cursor Print.
func (s *R) Print() { s.Cur.Print() }

// Log delegates to internal cursor Log.
func (s *R) Log() { s.Cur.Log() }

// Mark returns a copy of the current scanner cursor to preserve like
// a bookmark into the buffer data. See Cur, Look, LookSlice.
func (s *R) Mark() *Cur {
	if s.Cur == nil {
		return nil
	}
	// force a copy
	cp := *s.Cur
	return &cp
}

// Snap sets an extra internal cursor to the current cursor. See Mark.
func (s *R) Snap() { s.Snapped = s.Mark() }

// Back jumps the current cursor to the last Snap (Snapped).
func (s *R) Back() { s.Jump(s.Snapped) }

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data. Beware, however, that the new cursor must originate
// from the same (or identical) data buffer or the values will be out of
// sync.
func (s *R) Jump(c *Cur) { nc := *c; s.Cur = &nc }

// Peek returns a string containing all the runes from the current
// scanner cursor position forward to the number of runes passed.
// If end of data is encountered it will return everything up until that
// point.  Also see Look and LookSlice.
func (s *R) Peek(n uint) string {
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
func (s *R) Look(to *Cur) string {
	if to.Byte < s.Cur.Byte {
		return string(s.Buf[to.Byte:s.Cur.Next])
	}
	return string(s.Buf[s.Cur.Byte:to.Next])
}

// LookSlice returns a string containing all the bytes from the first
// cursor up to the second cursor. Neither cursor position is changed.
func (s *R) LookSlice(beg *Cur, end *Cur) string {
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
func (s *R) Expect(expr ...any) (*Cur, error) {
	var beg, end *Cur
	beg = s.Cur

	for _, m := range expr {

		// please keep the most common at the top
		switch v := m.(type) {

		case rune: // ------------------------------------------------------
			if v != tk.ANY && s.Cur.Rune != v {
				err := s.ErrorExpected(m)
				s.Jump(beg)
				return nil, err
			}
			end = s.Mark()
			s.Scan()

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

		case is.Toi: // -----------------------------------------------------
			var m *Cur
			for m == nil && s.Cur.Rune != tk.EOD {
				for _, i := range v {
					m, _ = s.check(i)
				}
				s.Scan()
			}
			if m == nil {
				err := s.ErrorExpected(v)
				s.Jump(beg)
				return nil, err
			}
			end = m

		case is.To: // -----------------------------------------------------
			var m, b4 *Cur
		OUT:
			for s.Cur.Rune != tk.EOD {
				for _, i := range v {
					m, _ = s.check(i)
					if m != nil {
						break OUT
					}
				}
				b4 = s.Mark()
				s.Scan()
			}
			if m == nil {
				err := s.ErrorExpected(v)
				return nil, err
			}
			end = b4

		case is.It: // ----------------------------------------------------
			var m *Cur
			for _, i := range v {
				m, _ = s.check(i)
				if m != nil {
					break
				}
			}
			if m == nil {
				return nil, s.ErrorExpected(v)
			}
			end = s.Mark()

		case is.Not: // ----------------------------------------------------
			m := s.Mark()
			for _, i := range v {
				if c, _ := s.check(i); c != nil {
					err := s.ErrorExpected(v, i)
					return nil, err
				}
			}
			end = m

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
			m := s.Mark()
			c, err := s.Expect(v...)
			if err != nil {
				s.Jump(m)
				return nil, err
			}
			end = c

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

		case is.Mn1: // ----------------------------------------------------
			m, err := s.Expect(is.Min{1, v.This})
			if err != nil {
				s.Jump(beg)
				return nil, s.ErrorExpected(v)
			}
			end = m

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

		case is.N: // ------------------------------------------------------
			m, err := s.Expect(is.MMx{v.N, v.N, v.This})
			if err != nil {
				s.Jump(beg)
				return nil, s.ErrorExpected(v)
			}
			end = m

		case is.Any: // ----------------------------------------------------
			for n := 0; n < v.N; n++ {
				s.Scan()
			}
			end = s.Mark()
			s.Scan()

		case is.Rng: // ----------------------------------------------------
			if !(v.First <= s.Cur.Rune && s.Cur.Rune <= v.Last) {
				err := s.ErrorExpected(v)
				s.Jump(beg)
				return nil, err
			}
			end = s.Mark()
			s.Scan()

		case Hook: // ------------------------------------------------------
			if !v(s) {
				return nil, fmt.Errorf(
					"expect: hook function failed (%v)",
					util.FuncName(v),
				)
			}
			end = s.Mark()

		case func(r *R) bool:
			if !v(s) {
				return nil, fmt.Errorf(
					"expect: hook function failed (%v)",
					util.FuncName(v),
				)
			}
			end = s.Mark()

		default: // --------------------------------------------------------
			return nil, fmt.Errorf("expect: unexpr expression (%T)", m)
		}
	}
	return end, nil
}

// ErrorExpected returns a verbose, one-line error describing what was
// expected when it encountered whatever the scanner last scanned. All
// expression types are supported. See Expect.
func (s *R) ErrorExpected(this any, args ...any) error {
	var msg string
	but := fmt.Sprintf(` at %v`, s)
	if s.Cur != nil && s.Cur.Rune == tk.EOD && s.Cur.Len == 0 {
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
	case is.It:
		msg = fmt.Sprintf(`expected %q`, v)
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
	case is.Mn1:
		str := `expected one or more %q`
		msg = fmt.Sprintf(str, v.This)
	case is.Min:
		str := `expected min %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.This)
	case is.MMx:
		str := `expected min %v, max %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.Max, v.This)
	case is.N:
		str := `expected exactly %v of %q`
		msg = fmt.Sprintf(str, v.N, v.This)
	case is.Rng:
		str := `expected range [%v-%v]`
		msg = fmt.Sprintf(str, string(v.First), string(v.Last))
	case is.To, is.Toi:
		str := `%q not found`
		msg = fmt.Sprintf(str, v)
	default:
		msg = fmt.Sprintf(`expected %T %q`, v, v)
	}
	return errors.New(msg + but)
}

// NewLine delegates to interval Curs.NewLine.
func (s *R) NewLine() { s.Cur.NewLine() }

// check behaves exactly like Expect but jumps back to the original
// cursor position after scanning for expected expression values.
func (s *R) check(expr ...any) (*Cur, error) {
	defer s.Jump(s.Mark())
	return s.Expect(expr...)
}
