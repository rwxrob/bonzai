// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package scanner

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/rwxrob/bonzai/scanner/is"
	"github.com/rwxrob/bonzai/scanner/tk"
)

// Scanner implements a non-linear, rune-centric, buffered data scanner.
// See New for creating a usable struct that implements Scanner. The
// buffer and cursor are directly exposed to facilitate
// higher-performance, direct access when needed.
type Scanner struct {
	Buf []byte
	Cur *Cur
}

// New returns a newly initialized non-linear, rune-centric, buffered
// data scanner with support for parsing data from io.Reader, string,
// and []byte types. Returns nil and the error if any encountered during
// initialization. Also see the Init method.
func New(i interface{}) (*Scanner, error) {
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
func (p *Scanner) Init(i interface{}) error {
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
func (p *Scanner) buffer(i interface{}) error {
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

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data. Beware, however, that the new cursor must originate
// from the same (or identical) data buffer or the values will be out of
// sync.
func (p *Scanner) Jump(c *Cur) { nc := *c; p.Cur = &nc }

// Peek returns a string containing all the runes from the current
// scanner cursor position forward to the number of runes passed.
// If end of data is countered will everything up until that point.
// Also so Look and LookSlice.
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

// Expect takes a variable list of parsable types including rune,
// string, is.In, is.Opt, is.Not, is.Seq, is.Min, is.MinMax, is.Count,
// and all strings from the tk subpackage. This allows for very readable
// functional grammar parsers to be created quickly without exceptional
// overhead from additional function calls and indirection. As some have
// said, "it's regex without the regex."
func (p *Scanner) Expect(scannable ...interface{}) (*Cur, error) {
	var beg, end *Cur
	beg = p.Cur
	for _, m := range scannable {

		// please keep the most common at the top
		switch v := m.(type) {

		case string:
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

		case rune:
			if p.Cur.Rune != v {
				err := p.ErrorExpected(m)
				p.Jump(beg)
				return nil, err
			}
			end = p.Mark()
			p.Scan()

		case is.Not:
			if _, e := p.Check(v.This); e == nil {
				err := p.ErrorExpected(v)
				p.Jump(beg)
				return nil, err
			}
			end = p.Mark()

			/*
				case Class:
					if !v.Check(p.Cur.Rune) {
						err := p.ErrorExpected(v)
						p.Jump(beg)
						return nil, err
					}
					end = p.Mark()
					p.Scan()

				case Check:
					rv, err := v.Check(p)
					if err != nil {
						p.Jump(beg)
						return nil, err
					}
					end = rv
					p.Jump(rv)
					p.Scan()

				case is.Opt:
					m, err := p.Expect(is.MinMax{v.This, 0, 1})
					if err != nil {
						p.Jump(beg)
						return nil, err
					}
					end = m

				case is.Min:
					c := 0
					last := p.Mark()
					var err error
					var m *Cur
					for {
						m, err = p.Expect(v.Match)
						if err != nil {
							break
						}
						last = m
						c++
					}
					if c < v.Min {
						p.Jump(beg)
						return nil, err
					}
					end = last

				case is.Count:
					m, err := p.Expect(is.MinMax{v.Match, v.Count, v.Count})
					if err != nil {
						p.Jump(beg)
						return nil, err
					}
					end = m

				case is.MinMax:
					c := 0
					last := p.Mark()
					var err error
					var m *Cur
					for {
						m, err = p.Expect(v.Match)
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
						return nil, err
					}
					end = last

				case is.Seq:
					m, err := p.Expect(v...)
					if err != nil {
						p.Jump(beg)
						return nil, err
					}
					end = m

				case is.OneOf:
					var m *Cur
					var err error
					for _, i := range v {
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
			*/

		default:
			return nil, fmt.Errorf("expect: unscannable type (%T)", m)
		}
	}
	return end, nil
}

// ErrorExpected returns a verbose, one-line error describing what was
// expected when it encountered whatever the scanner last scanned. All
// scannable types are supported. See Expect.
func (p *Scanner) ErrorExpected(this interface{}) error {
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
		msg = fmt.Sprintf(`not expecting %q`, v.This)
	default:
		msg = fmt.Sprintf(`expected %T %q`, v, v)
	}
	return errors.New(msg + but)
}

// NewLine delegates to interval Curs.NewLine.
func (p *Scanner) NewLine() { p.Cur.NewLine() }

// Check behaves exactly like Expect but jumps back to the original
// cursor position after scanning for expected scannable values.
func (p *Scanner) Check(scannable ...interface{}) (*Cur, error) {
	defer p.Jump(p.Mark())
	return p.Expect(scannable...)
}
