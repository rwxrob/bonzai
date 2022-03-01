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
	p := new(Scanner)
	if err := p.Init(i); err != nil {
		return nil, err
	}
	return p, nil
}

// Init reads all of passed parsable data (io.Reader, string, []byte)
// into buffered memory, scans the first rune, and sets the internals of
// scanner appropriately returning an error if anything happens while
// attempting to read and buffer the data (OOM, etc.).
func (p *Scanner) Init(i any) error {
	if err := p.buffer(i); err != nil {
		return err
	}
	r, ln := utf8.DecodeRune(p.Buf) // scan first
	if ln == 0 {
		r = tk.EOD
		return fmt.Errorf("scanner: failed to scan first rune")
	}
	p.Cur = new(Cur)
	p.Cur.Rune = r
	p.Cur.Len = ln
	p.Cur.Next = ln
	p.Cur.Pos.Line = 1
	p.Cur.Pos.LineRune = 1
	p.Cur.Pos.LineByte = 1
	p.Cur.Pos.Rune = 1
	return nil
}

// reads and buffers io.Reader, string, or []byte types
func (p *Scanner) buffer(i any) error {
	var err error
	switch in := i.(type) {
	case io.Reader:
		p.Buf, err = io.ReadAll(in)
		if err != nil {
			return err
		}
	case string:
		p.Buf = []byte(in)
	case []byte:
		p.Buf = in
	default:
		return fmt.Errorf("scanner: unsupported input type: %t", i)
	}
	if len(p.Buf) == 0 {
		return fmt.Errorf("scanner: no input")
	}
	return err
}

// Scan decodes the next rune and advances the scanner cursor by one.
func (p *Scanner) Scan() {
	if p.Done() {
		return
	}
	r, ln := utf8.DecodeRune(p.Buf[p.Cur.Next:])
	if ln != 0 {
		p.Cur.Byte = p.Cur.Next
		p.Cur.Pos.LineByte += p.Cur.Len
	} else {
		r = tk.EOD
	}
	p.Cur.Rune = r
	p.Cur.Pos.Rune += 1
	p.Cur.Next += ln
	p.Cur.Pos.LineRune += 1
	p.Cur.Len = ln
}

// ScanN scans the next n runes advancing n runes forward or returns
// p.Done() if attempted after already at end of data.
func (p *Scanner) ScanN(n int) {
	for i := 0; i < n; i++ {
		p.Scan()
	}
}

// Done returns true if current cursor rune is tk.EOD and the cursor length
// is also zero.
func (p *Scanner) Done() bool {
	return p.Cur.Rune == tk.EOD && p.Cur.Len == 0
}

// String delegates to internal cursor String.
func (p *Scanner) String() string { return p.Cur.String() }

// Print delegates to internal cursor Print.
func (p *Scanner) Print() { p.Cur.Print() }

// Mark returns a copy of the current scanner cursor to preserve like
// a bookmark into the buffer data. See Cur, Look, LookSlice.
func (p *Scanner) Mark() *Cur {
	if p.Cur == nil {
		return nil
	}
	// force a copy
	cp := *p.Cur
	return &cp
}

// Snap sets an extra internal cursor to the current cursor. See Mark.
func (p *Scanner) Snap() { p.Snapped = p.Mark() }

// Back jumps the current cursor to the last Snap (Snapped).
func (p *Scanner) Back() { p.Jump(p.Snapped) }

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data. Beware, however, that the new cursor must originate
// from the same (or identical) data buffer or the values will be out of
// sync.
func (p *Scanner) Jump(c *Cur) { nc := *c; p.Cur = &nc }

// Peek returns a string containing all the runes from the current
// scanner cursor position forward to the number of runes passed.
// If end of data is encountered it will return everything up until that
// point.  Also see Look and LookSlice.
func (p *Scanner) Peek(n uint) string {
	buf := ""
	pos := p.Cur.Byte
	for c := uint(0); c < n; c++ {
		r, ln := utf8.DecodeRune(p.Buf[pos:])
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
func (p *Scanner) Look(to *Cur) string {
	if to.Byte < p.Cur.Byte {
		return string(p.Buf[to.Byte:p.Cur.Next])
	}
	return string(p.Buf[p.Cur.Byte:to.Next])
}

// LookSlice returns a string containing all the bytes from the first
// cursor up to the second cursor. Neither cursor position is changed.
func (p *Scanner) LookSlice(beg *Cur, end *Cur) string {
	return string(p.Buf[beg.Byte:end.Next])
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
func (p *Scanner) Expect(scannables ...any) (*Cur, error) {
	var beg, end *Cur
	beg = p.Cur

	for _, m := range scannables {

		// please keep the most common at the top
		switch v := m.(type) {

		case string: // ----------------------------------------------------
			if v == "" {
				return nil, fmt.Errorf("expect: cannot parse empty string")
			}
			for _, r := range []rune(v) {
				if p.Cur.Rune != r {
					err := p.ErrorExpected(r)
					p.Jump(beg)
					return nil, err
				}
				end = p.Mark()
				p.Scan()
			}

		case rune: // ------------------------------------------------------
			if p.Cur.Rune != v {
				err := p.ErrorExpected(m)
				p.Jump(beg)
				return nil, err
			}
			end = p.Mark()
			p.Scan()

		case is.Not: // ----------------------------------------------------
			for _, i := range v {
				if _, e := p.Check(i); e == nil {
					err := p.ErrorExpected(v, i)
					p.Jump(beg)
					return nil, err
				}
			}
			end = p.Mark()

		case is.In: // -----------------------------------------------------
			var m *Cur
			for _, i := range v {
				var err error
				last := p.Mark()
				m, err = p.Expect(i)
				if err == nil {
					break
				}
				p.Jump(last)
			}
			if m == nil {
				return nil, p.ErrorExpected(v)
			}
			end = m

		case is.Seq: // ----------------------------------------------------
			m, err := p.Expect(v...)
			if err != nil {
				p.Jump(beg)
				return nil, err
			}
			end = m

		case is.Opt: // ----------------------------------------------------
			var m *Cur
			for _, i := range v {
				var err error
				m, err = p.Expect(is.MMx{0, 1, i})
				if err != nil {
					p.Jump(beg)
					return nil, p.ErrorExpected(v)
				}
			}
			end = m

		case is.MMx: // ----------------------------------------------------
			c := 0
			last := p.Mark()
			var err error
			var m *Cur
			for {
				m, err = p.Expect(v.This)
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
				p.Jump(last)
				return nil, p.ErrorExpected(v)
			}
			end = last

		case is.Min: // ----------------------------------------------------
			c := 0
			last := p.Mark()
			var err error
			var m *Cur
			for {
				m, err = p.Expect(v.This)
				if err != nil {
					break
				}
				last = m
				c++
			}
			if c < v.Min {
				p.Jump(beg)
				return nil, p.ErrorExpected(v)
			}
			end = last

		case is.X: // ------------------------------------------------------
			m, err := p.Expect(is.MMx{v.X, v.X, v.This})
			if err != nil {
				p.Jump(beg)
				return nil, p.ErrorExpected(v)
			}
			end = m

		case is.Rng: // ----------------------------------------------------
			if !(v.First <= p.Cur.Rune && p.Cur.Rune <= v.Last) {
				err := p.ErrorExpected(v)
				p.Jump(beg)
				return nil, err
			}
			end = p.Mark()
			p.Scan()

		default: // --------------------------------------------------------
			if p.ExtendExpect != nil {
				return p.ExtendExpect(p, scannables...)
			}
			return nil, fmt.Errorf("expect: unscannable type (%T)", m)
		}
	}
	return end, nil
}

// ErrorExpected returns a verbose, one-line error describing what was
// expected when it encountered whatever the scanner last scanned. All
// scannable types are supported. See Expect.
func (p *Scanner) ErrorExpected(this any, args ...any) error {
	var msg string
	but := fmt.Sprintf(` at %v`, p)
	if p.Done() {
		runes := `runes`
		if p.Cur.Pos.Rune == 1 {
			runes = `rune`
		}
		but = fmt.Sprintf(`, exceeded data length (%v %v)`,
			p.Cur.Pos.Rune, runes)
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
func (p *Scanner) NewLine() { p.Cur.NewLine() }

// Check behaves exactly like Expect but jumps back to the original
// cursor position after scanning for expected scannable values.
func (p *Scanner) Check(scannables ...any) (*Cur, error) {
	defer p.Jump(p.Mark())
	return p.Expect(scannables...)
}
