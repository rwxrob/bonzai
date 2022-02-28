package scanner

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

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
// and []byte types. Also see the Init method.
func New(i interface{}) *Scanner {
	p := new(Scanner)
	if err := p.Init(i); err != nil {
		return p
	}
	return nil
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

// CopyCur returns a copy of the current scanner cursor. See Cur.
func (p *Scanner) CopyCur() *Cur {
	if p.Cur == nil {
		return nil
	}
	// force a copy
	cp := *p.Cur
	return &cp
}

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the scanner's current position in the
// buffered data.
func (p *Scanner) Jump(c *Cur) { nc := *c; p.Cur = &nc }

// Look returns a string containing all the bytes from the current
// scanner cursor position up to the passed cursor position, forward or
// backward. Neither the internal nor the passed  cursor position is
// changed.
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
// string, Class, Check, Opt, Not, Seq, One, Min, MinMax, Count. This
// allows grammars to be represented simply and parsed easily without
// exceptional overhead from additional function calls and indirection.
func (p *Scanner) Expect(scannable ...interface{}) (*Cur, error) {
	var beg, end *Cur
	beg = p.Cur
	for _, m := range scannable {
		switch v := m.(type) {

		case rune:
			if p.Cur.Rune != v {
				err := p.ErrorExpected(m)
				p.Jump(beg)
				return nil, err
			}
			end = p.CopyCur()
			p.Scan()

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
				end = p.CopyCur()
				p.Scan()
			}
			/*
				case Class:
					if !v.Check(p.Cur.Rune) {
						err := p.ErrorExpected(v)
						p.Jump(beg)
						return nil, err
					}
					end = p.CopyCur()
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

				case is.Not:
					if _, e := p.Check(v.This); e == nil {
						err := p.ErrorExpected(v)
						p.Jump(beg)
						return nil, err
					}
					end = p.CopyCur()

				case is.Min:
					c := 0
					last := p.CopyCur()
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
					last := p.CopyCur()
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
						last := p.CopyCur()
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
			return nil, fmt.Errorf("expect: unsupported argument type (%T)", m)
		}
	}
	return end, nil
}

func (p *Scanner) ErrorExpected(this interface{}) error {
	var msg string
	but := fmt.Sprintf(` but got %v`, p)
	if p.Done() {
		runes := `runes`
		if p.Cur.Pos.Rune == 1 {
			runes = `rune`
		}
		but = fmt.Sprintf(` but exceeded data length (%v %v)`, p.Cur.Pos.Rune, runes)
	}
	// TODO add verbose errors for *all* types in Grammar
	switch v := this.(type) {
	case string:
		msg = fmt.Sprintf(`expected string %q`, v)
	case rune:
		msg = fmt.Sprintf(`expected rune %q`, v)
		/*
			case Class:
				msg = fmt.Sprintf(`expected class %v (%v)`, v.Ident(), v.Desc())
			default:
				msg = fmt.Sprintf(`expected %T %q`, v, v)
		*/
	}
	return errors.New(msg + but)
}

// NewLine delegates to interval Curs.NewLine.
func (p *Scanner) NewLine() { p.Cur.NewLine() }

func (p *Scanner) Check(ms ...interface{}) (*Cur, error) {
	defer p.Jump(p.CopyCur())
	return p.Expect(ms...)
}
