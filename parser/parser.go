package parser

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// Parser must implement a parser with all the trimmings. See New for
// a creating a specific instance of Parser.
type Parser interface {
	Check(ms ...any) (*Cur, error)
	Cur() *Cur
	CopyCur() *Cur
	Done() bool
	ErrorExpected(this any) error
	Expect(ms ...any) (*Cur, error)
	Init(i any)
	Jump(m *Cur)
	Move(n int)
	NewLine()
	Next()
	Look(to *Cur) string
	LookSlice(beg *Cur, end *Cur) string
	Print()
	String() string
}

// EOD is a special value that is returned when the end of data is
// reached enabling functional parser functions to look for it reliably
// no matter what is being parsed.
const EOD = 1<<31 - 1 // max int32

// see Parser interface
type parser struct {
	in  io.Reader
	buf []byte
	cur *Cur
}

// New returns a newly initialized rune-centric Parser with support for
// parsing data from io.Reader, string, and []byte types. See the Parser
// interface and Init method.
func New(i interface{}) *parser {
	p := new(parser)
	if err := p.Init(i); err != nil {
		return p
	}
	return nil
}

// Init reads all of passed parsable (io.Reader, string, []byte) into
// memory, parses the first rune, and sets the internals of parser
// appropriately returning an error if anything happens while attempting
// to read and buffer the data (OOM, etc.).
func (p *parser) Init(i interface{}) error {
	if err := p.buffer(i); err != nil {
		return err
	}
	r, ln := utf8.DecodeRune(p.buf) // scan first
	if ln == 0 {
		r = EOD
		return fmt.Errorf("parser: failed to scan first rune")
	}
	p.cur = new(Cur)
	p.cur.Rune = r
	p.cur.Len = ln
	p.cur.Next = ln
	p.cur.Pos.Line = 1
	p.cur.Pos.LineRune = 1
	p.cur.Pos.LineByte = 1
	p.cur.Pos.Rune = 1
	return nil
}

// reads and buffers io.Reader, string, or []byte types
func (p *parser) buffer(i interface{}) error {
	var err error
	switch in := i.(type) {
	case io.Reader:
		p.buf, err = io.ReadAll(in)
		if err != nil {
			return err
		}
	case string:
		p.buf = []byte(in)
	case []byte:
		p.buf = in
	default:
		return fmt.Errorf("parser: unsupported input type: %t", i)
	}
	if len(p.buf) == 0 {
		return fmt.Errorf("parser: no input")
	}
	return err
}

// Next parses the next rune advancing a single rune forward or sets
// current cursor rune to EOD if at end of data. Returns p.Done() if
// attempted after already at end of data.
func (p *parser) Next() {
	if p.Done() {
		return
	}
	r, ln := utf8.DecodeRune(p.buf[p.cur.Next:])
	if ln != 0 {
		p.cur.Byte = p.cur.Next
		p.cur.Pos.LineByte += p.cur.Len
	} else {
		r = EOD
	}
	p.cur.Rune = r
	p.cur.Pos.Rune += 1
	p.cur.Next += ln
	p.cur.Pos.LineRune += 1
	p.cur.Len = ln
}
func (p *parser) Move(n int) {
	for i := 0; i < n; i++ {
		p.Next()
	}
}

// Done returns true if current cursor rune is EOD and the cursor length
// is also zero.
func (p *parser) Done() bool {
	return p.cur.Rune == EOD && p.cur.Len == 0
}

// String delegates to internal cursor String.
func (p *parser) String() string { return p.cur.String() }

// Print delegates to internal cursor Print.
func (p *parser) Print() { p.cur.Print() }

// Cur returns exact cursor used by parser. See CopyCur and Cur struct.
func (p *parser) Cur() *Cur { return p.cur }

// CopyCur returns a copy of the current parser cursor. See Cur.
func (p *parser) CopyCur() *Cur {
	if p.cur == nil {
		return nil
	}
	// force a copy
	cp := *p.cur
	return &cp
}

// Jump replaces the internal cursor with a copy of the one passed
// effectively repositioning the parser's current position in the
// buffered data.
func (p *parser) Jump(c *Cur) { nc := *c; p.cur = &nc }

// Look returns a string containing all the bytes from the current
// parser cursor position up to the passed cursor position, forward or
// backward. Neither the internal nor the passed  cursor position is
// changed.
func (p *parser) Look(to *Cur) string {
	if to.Byte < p.cur.Byte {
		return string(p.buf[to.Byte:p.cur.Next])
	}
	return string(p.buf[p.cur.Byte:to.Next])
}

// LookSlice returns a string containing all the bytes from the first
// cursor up to the second cursor. Neither cursor position is changed.
func (p *parser) LookSlice(beg *Cur, end *Cur) string {
	return string(p.buf[beg.Byte:end.Next])
}

// Expect takes a variable list of parsable types including rune,
// string, Class, Check, Opt, Not, Seq, One, Min, MinMax, Count. This
// allows grammars to be represented simply and parsed easily without
// exceptional overhead from additional function calls and indirection.
func (p *parser) Expect(ms ...interface{}) (*Cur, error) {
	var beg, end *Cur
	beg = p.Cur()
	for _, m := range ms {
		switch v := m.(type) {

		case rune:
			if p.cur.Rune != v {
				err := p.ErrorExpected(m)
				p.Jump(beg)
				return nil, err
			}
			end = p.CopyCur()
			p.Next()

		case string:
			if v == "" {
				return nil, fmt.Errorf("expect: cannot parse empty string")
			}
			for _, r := range []rune(v) {
				if p.cur.Rune != r {
					err := p.ErrorExpected(r)
					p.Jump(beg)
					return nil, err
				}
				end = p.CopyCur()
				p.Next()
			}
			/*
				case Class:
					if !v.Check(p.cur.Rune) {
						err := p.ErrorExpected(v)
						p.Jump(beg)
						return nil, err
					}
					end = p.CopyCur()
					p.Next()

				case Check:
					rv, err := v.Check(p)
					if err != nil {
						p.Jump(beg)
						return nil, err
					}
					end = rv
					p.Jump(rv)
					p.Next()

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

func (p *parser) ErrorExpected(this interface{}) error {
	var msg string
	but := fmt.Sprintf(` but got %v`, p)
	if p.Done() {
		runes := `runes`
		if p.cur.Pos.Rune == 1 {
			runes = `rune`
		}
		but = fmt.Sprintf(` but exceeded data length (%v %v)`, p.cur.Pos.Rune, runes)
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
func (p *parser) NewLine() { p.cur.NewLine() }

func (p *parser) Check(ms ...interface{}) (*Cur, error) {
	defer p.Jump(p.CopyCur())
	return p.Expect(ms...)
}
