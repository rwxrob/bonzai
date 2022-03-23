package scan

import (
	"fmt"
	"log"
)

// LogPanic is used with defer to trap any panic and log it.
func (s *R) LogPanic() {
	r := recover()
	if r != nil {
		log.Printf("%v at %v", r, s)
	}
}

// LogPanic is used with defer to trap any panic and print it.
func (s *R) PrintPanic() {
	r := recover()
	if r != nil {
		fmt.Printf("%v at %v", r, s)
	}
}

// --------------------------- PEGN support ---------------------------
// The following are directly correlated to supported PEGN expressions
// and are intended to be generated from PEGN grammars specifically
// (altough other languages generators could easily be adapted).

// Rune matches the exact rune specified or panics.
func (s *R) Rune(r rune) {
	if s.Cur.Rune != r {
		panic(fmt.Sprintf("expected %q", r))
	}
	s.Scan()
}

// Str iterates over all of the runes in the string as if Rune were
// called on each.
func (s *R) Str(v string) {
	for _, v := range []rune(v) {
		if v != s.Cur.Rune {
			panic(fmt.Sprintf("expected %q", v))
		}
		s.Scan()
	}
}
