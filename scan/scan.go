// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package scan implements a non-linear, rune-centric, buffered data,
extendable scanner that includes its own high-level parsing expression
grammar (BPEGN) comprised of Go constants, slices, and structs used as
expressions from the z ("is") subpackage making parser generation (by
hand or code) trivial from any structured meta languages such as PEGN,
PEG, EBNF, ABNF, etc. Most will use the scanner to create parsers
quickly where a regular expression will not suffice. See the z ("is")
and tk packages for a growing number of common, centrally maintain
expressions for your parsing pleasure. Also see the mark (BonzaiMark)
subpackage for a working examples of the scanner in action, which is
used by the included Bonzai help.Cmd command and others.
*/
package scan

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	z "github.com/rwxrob/bonzai/scan/is"
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

	// Last contains the previous Cur value when Scan was called.
	Last *Cur

	// Snapped contains the latest Cur when Snap was called.
	Snapped *Cur

	// State contains scanner state bitfield using a combination of the
	// different constants. Currently, only EOD is used but others may be
	// added as more state-modifying single-token expressions are
	// considered (like tk.IS and tk.NOT now).
	State int
}

const (
	EOD = 1 << iota // reached EOD, current rune is tk.EOD
)

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
		s.State |= EOD
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

// Scan decodes the next rune and advances the scanner cursor by one
// saving the last cursor into s.Last.
func (s *R) Scan() {
	s.Last = s.Mark()

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

// Rewind Jumps to Last.
func (s *R) Rewind() { s.Jump(s.Last) }

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

// The Expect method is primarily designed for ease of use quickly by
// allowing BPEGN expressions to be coded much faster than traditional
// parser functions. In fact, entire applications can be generated
// automatically from BPEGN, PEGN, ABNF, EBNF, or any other meta syntax
// language by adding additional scan.Hooks to do work besides just
// parsing what is scanned. Generate the BPEGN Go expressions to pass to
// Expect and code the hook callbacks. That's it. Coding in BPEGN comes
// naturally when putting together a parser quickly and without the
// cognitive overhead of switching between grammars. The BPEGN
// expressions passed to Expect are 100% Go. For full documentation of
// BPEGN expressions see the z ("is") package and the source of this
// Expect method.
//
// Warning: While it is nothing for most developers to be concerned
// with, the Expect method does do a fair amount of functional recursion
// for the sake of simplicity and to support BPEGN syntax. Those wishing
// to gain the maximum performance should consider using other scan.R
// methods instead in such cases. Developers are encouraged to do their
// own benchmarking and perhaps start with BPEGN until they can create
// more optimized parsers when and if necessary. Most will discover
// other more substantial bottlenecks. The Bonzai project places
// priority on speed and quality of developer delivery over run-time
// performance. Delivery time is far more costly than the minimal gains
// in run-time performance. "Premature optimization is the root of all
// evil," as they say.
func (s *R) Expect(expr any) (*Cur, error) {
	s.Last = s.Cur

	// please keep the most common expressions types at the top

	switch v := expr.(type) {

	case rune: // ------------------------------------------------------
		if v != tk.ANY && s.Cur.Rune != v {
			err := s.ErrorExpected(v)
			return nil, err
		}
		s.Scan()
		return s.Last, nil

	case string: // ----------------------------------------------------
		if v == "" {
			return s.Mark(), nil
		}
		// avoid the temptation to look directly at bytes since it would
		// allow none runes to be passed within "strings"
		for _, v := range []rune(v) {
			if v != s.Cur.Rune {
				return nil, s.ErrorExpected(v)
			}
			s.Scan()
		}
		return s.Last, nil

	case z.X: // -----------------------------------------------------
		var err error
		b := s.Mark()
		m := s.Mark()
		for _, i := range v {
			m, err = s.Expect(i)
			if err != nil {
				s.Jump(b)
				return nil, err
			}
		}
		return m, nil

	case z.O: // -------------------------------------------------------
		for _, i := range v {
			m, _ := s.Expect(i)
			if m != nil {
				return m, nil
			}
		}

	case z.Toi: // -----------------------------------------------------
		back := s.Mark()
		for s.Cur.Rune != tk.EOD {
			for _, i := range v {
				m, _ := s.Expect(i)
				if m != nil {
					return m, nil
				}
			}
			s.Scan()
		}
		s.Jump(back)
		return nil, s.ErrorExpected(v)

	case z.To: // -----------------------------------------------------
		m := s.Mark()
		b4 := s.Mark()
		for s.Cur.Rune != tk.EOD {
			for _, i := range v {
				b := s.Mark()
				c, _ := s.Expect(i)
				if c != nil {
					s.Jump(b)
					return b4, nil
				}
			}
			b4 = s.Mark()
			s.Scan()
		}
		s.Jump(m)
		return nil, s.ErrorExpected(v)

	case z.It: // ----------------------------------------------------
		var m *Cur
		b := s.Mark()
		for _, i := range v {
			m, _ = s.Expect(i)
			if m != nil {
				break
			}
		}
		if m == nil {
			return nil, s.ErrorExpected(v)
		}
		s.Jump(b)
		return b, nil

	case z.Not: // ----------------------------------------------------
		m := s.Mark()
		for _, i := range v {
			if c, _ := s.Expect(i); c != nil {
				s.Jump(m)
				err := s.ErrorExpected(v, i)
				return nil, err
			}
		}
		return m, nil

	case z.In: // -----------------------------------------------------
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
		return m, nil

	case z.MMx: // ----------------------------------------------------
		c := 0
		last := s.Mark()
		var err error
		var m, end *Cur
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
		}
		if !(v.Min <= c && c <= v.Max) {
			s.Jump(last)
			return nil, s.ErrorExpected(v)
		}
		return end, nil

	case z.Mn1: // ----------------------------------------------------
		m := s.Mark()
		c, err := s.Expect(z.Min{1, v.This})
		if err != nil {
			s.Jump(m)
			return nil, s.ErrorExpected(v)
		}
		return c, nil

	case z.Min: // ----------------------------------------------------
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
			s.Jump(last)
			return nil, s.ErrorExpected(v)
		}
		return last, nil

	case z.C: // ------------------------------------------------------
		b := s.Mark()
		m, err := s.Expect(z.MMx{v.N, v.N, v.This})
		if err != nil {
			s.Jump(b)
			return nil, s.ErrorExpected(v)
		}
		return m, nil

	case z.Any: // ----------------------------------------------------
		for n := 0; n < v.N; n++ {
			s.Scan()
		}
		m := s.Mark()
		s.Scan()
		return m, nil

	case z.Rng: // ----------------------------------------------------
		if !(v.First <= s.Cur.Rune && s.Cur.Rune <= v.Last) {
			err := s.ErrorExpected(v)
			return nil, err
		}
		m := s.Mark()
		s.Scan()
		return m, nil

	case Hook: // ------------------------------------------------------
		if !v(s) {
			return nil, fmt.Errorf(
				"expect: hook function failed (%v)",
				util.FuncName(v),
			)
		}
		return s.Cur, nil

	case func(r *R) bool:
		if !v(s) {
			return nil, fmt.Errorf(
				"expect: hook function failed (%v)",
				util.FuncName(v),
			)
		}
		return s.Cur, nil

	}
	return nil, fmt.Errorf("unknown expression (%T)", expr)
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
	case z.It:
		if len(v) > 1 {
			msg = fmt.Sprintf(`expected one of %q`, v)
		} else {
			msg = fmt.Sprintf(`expected %q`, v[0])
		}
	case z.Not:
		msg = fmt.Sprintf(`unexpected %q`, args[0])
	case z.In:
		str := `expected one of %q`
		msg = fmt.Sprintf(str, v)
	case z.X:
		//str := `expected %q in sequence %q at %v beginning`
		//msg = fmt.Sprintf(str, args[0], v, args[1])
		str := `expected %q in sequence`
		msg = fmt.Sprintf(str, v)
	case z.O:
		str := `expected an optional %v`
		msg = fmt.Sprintf(str, v)
	case z.Mn1:
		str := `expected one or more %q`
		msg = fmt.Sprintf(str, v.This)
	case z.Min:
		str := `expected min %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.This)
	case z.MMx:
		str := `expected min %v, max %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.Max, v.This)
	case z.C:
		str := `expected exactly %v of %q`
		msg = fmt.Sprintf(str, v.N, v.This)
	case z.Rng:
		str := `expected range [%v-%v]`
		msg = fmt.Sprintf(str, string(v.First), string(v.Last))
	case z.Toi:
		str := `%q not found`
		if len(v) > 1 {
			str = `none of %q found`
		}
		msg = fmt.Sprintf(str, v)
	case z.To:
		msg = fmt.Sprintf(`none of %q found`, v)
	default:
		msg = fmt.Sprintf(`expected %T %q`, v, v)
	}
	return errors.New(msg + but)
}

// NewLine delegates to interval Curs.NewLine.
func (s *R) NewLine() { s.Cur.NewLine() }
