package scan

/*

// Parse creates a new Node reference from the string returned by Look.
func (s *R) Parse(to *Cur) *Node {
	n := new(Node)
	n.V = s.Look(to)
	return n
}

// ParseSlice creates a new Node reference from the string returned by
// LookSlice.
func (s *R) ParseSlice(b *Cur, e *Cur) *Node {
	n := new(Node)
	n.V = s.LookSlice(b, e)
	return n
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
		if s.Cur.Rune != v {
			err := s.ErrorExpected(v)
			return nil, err
		}
		s.Scan()
		return s.Last, nil

	case tk.Token: // --------------------------------------------------
		switch v {
		case tk.ANY: // A, A1
			s.Scan()
		case tk.A2:
			s.ScanN(2)
		case tk.A3:
			s.ScanN(3)
		case tk.A4:
			s.ScanN(4)
		case tk.A5:
			s.ScanN(5)
		case tk.A6:
			s.ScanN(6)
		case tk.A7:
			s.ScanN(7)
		case tk.A8:
			s.ScanN(8)
		case tk.A9:
			s.ScanN(9)
		}
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
	default:
		return fmt.Errorf("scanner: unsupported input type: %t", i)
	}
	s.BufLen = len(s.Buf)
	if s.BufLen == 0 {
		return fmt.Errorf("scanner: no input")
	}
	return err
}

/*

// Parse creates a new Node reference from the string returned by Look.
func (s *R) Parse(to *Cur) *Node {
	n := new(Node)
	n.V = s.Look(to)
	return n
}

// ParseSlice creates a new Node reference from the string returned by
// LookSlice.
func (s *R) ParseSlice(b *Cur, e *Cur) *Node {
	n := new(Node)
	n.V = s.LookSlice(b, e)
	return n
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
		if s.Cur.Rune != v {
			err := s.ErrorExpected(v)
			return nil, err
		}
		s.Scan()
		return s.Last, nil

	case tk.Token: // --------------------------------------------------
		switch v {
		case tk.ANY: // A, A1
			s.Scan()
		case tk.A2:
			s.ScanN(2)
		case tk.A3:
			s.ScanN(3)
		case tk.A4:
			s.ScanN(4)
		case tk.A5:
			s.ScanN(5)
		case tk.A6:
			s.ScanN(6)
		case tk.A7:
			s.ScanN(7)
		case tk.A8:
			s.ScanN(8)
		case tk.A9:
			s.ScanN(9)
		}
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

	case z.Ti: // -----------------------------------------------------
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

	case z.T: // -----------------------------------------------------
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

	case z.Y: // ----------------------------------------------------
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

	case z.N: // ----------------------------------------------------
		m := s.Mark()
		for _, i := range v {
			if c, _ := s.Expect(i); c != nil {
				s.Jump(m)
				err := s.ErrorExpected(v, i)
				return nil, err
			}
		}
		return m, nil

	case z.I: // -----------------------------------------------------
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

	case z.MM: // ----------------------------------------------------
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

	case z.M1: // ----------------------------------------------------
		m := s.Mark()
		c, err := s.Expect(z.M{1, v.This})
		if err != nil {
			s.Jump(m)
			return nil, s.ErrorExpected(v)
		}
		return c, nil

	case z.M: // ----------------------------------------------------
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

	case z.A: // ----------------------------------------------------
		for n := 0; n < v.N; n++ {
			s.Scan()
		}
		s.Scan()
		return s.Last, nil
		// see rune for A2-9

	case z.R: // ----------------------------------------------------
		if !(v.First <= s.Cur.Rune && s.Cur.Rune <= v.Last) {
			err := s.ErrorExpected(v)
			return nil, err
		}
		s.Scan()
		return s.Last, nil

	case z.C: // ------------------------------------------------------
		return s.expcount(v.N, v.This)
	case z.C2:
		return s.expcount(2, v.This)
	case z.C3:
		return s.expcount(3, v.This)
	case z.C4:
		return s.expcount(4, v.This)
	case z.C5:
		return s.expcount(5, v.This)
	case z.C6:
		return s.expcount(6, v.This)
	case z.C7:
		return s.expcount(7, v.This)
	case z.C8:
		return s.expcount(8, v.This)
	case z.C9:
		return s.expcount(9, v.This)

	case Hook: // ------------------------------------------------------
		if err := v(s); err != nil {
			return nil, s.ErrorExpected(v, err)
		}
		return s.Cur, nil

	case func(r *R) error:
		if err := v(s); err != nil {
			return nil, s.ErrorExpected(v, err)
		}
		return s.Cur, nil

	}
	return nil, fmt.Errorf("unknown expression (%T)", expr)
}

// handles all C* expressions
func (s *R) expcount(n int, expr any) (*Cur, error) {
	b := s.Mark()
	m := s.Mark()
	for i := 0; i < n; i++ {
		m, _ := s.Expect(expr)
		if m == nil {
			s.Jump(b)
			return nil, s.ErrorExpected(expr)
		}
	}
	return m, nil
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
	case z.Y:
		if len(v) > 1 {
			msg = fmt.Sprintf(`expected one of %q`, v)
		} else {
			msg = fmt.Sprintf(`expected %q`, v[0])
		}
	case z.N:
		msg = fmt.Sprintf(`unexpected %q`, args[0])
	case z.I:
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
	case z.M1:
		str := `expected one or more %q`
		msg = fmt.Sprintf(str, v.This)
	case z.M:
		str := `expected min %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.This)
	case z.MM:
		str := `expected min %v, max %v of %q`
		msg = fmt.Sprintf(str, v.Min, v.Max, v.This)
	case z.R:
		str := `expected range [%v-%v]`
		msg = fmt.Sprintf(str, string(v.First), string(v.Last))
	case z.Ti:
		str := `%q not found`
		if len(v) > 1 {
			str = `none of %q found`
		}
		msg = fmt.Sprintf(str, v)
	case z.T:
		msg = fmt.Sprintf(`none of %q found`, v)
	case Hook:
		msg = fmt.Sprintf("%v: %v", strings.ToLower(util.FuncName(v)), args[0])
	case func(s *R) error:
		msg = fmt.Sprintf("%v: %v", strings.ToLower(util.FuncName(v)), args[0])
	default:
		msg = fmt.Sprintf(`expected %T %q`, v, v)
	}
	return errors.New(msg + but)
}

// --------------------------- new functions --------------------------

func (s *R) Any(n int) bool {
	for i := 0; i < n; i++ {
		s.Scan()
	}
	return true
}

func (s *R) Str(strs ...string) {
	for _, it := range strs {
		for _, r := range []rune(it) {
			if r != s.Cur.Rune {
				panic(fmt.Sprintf("expecting %q", r))
			}
			s.Scan()
		}
	}
}

func (s *R) Opt(strs ...string) {
	s.Snap()
OUT:
	for _, it := range strs {
		for _, r := range []rune(it) {
			if r != s.Cur.Rune {
				s.Back()
				continue OUT
			}
			s.Scan()
		}
	}
}
