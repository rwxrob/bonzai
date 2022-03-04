// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package is defines the Bonzai Parsing Expression Grammar Notation
(BPEGN) implemented entirely using Go types (mostly slices and structs).
BPEGN can be 100% transpiled to and from the Parsing Expression Grammer
Notation (PEGN). Code in one and use "bonzai pegn" command to easily
generate the other. BPEGN is sometimes referred to informally as "Bonzai
Scanner Expect" language as well since it is passed directly to
scan.Expect or scan.Check.

Parameterized Struct and Set Slice Expressions

BPEGN uses Go structs and slices to represent scannable expressions with
obvious similarities to regular expressions. These expressions have
one-to-parity with PEGN expressions. Slices represent sets of
possibilities. Structs provide parameters for more complex expressions
and are are guaranteed never to change allowing them to be dependably
used in assignment without struct field names using Go's inline
composable syntax. Some editors may need configuring to allow this since
in general practice this can create subtle (but substantial) foot-guns
for maintainers.

Look-Ahead and Advancing Expressions

"Advancing" expressions will advance the scan to the end of the
expression match. "Look-ahead" expressions simply check for a match but
do not advance the scan. Developers should take careful note of the
difference in the documentation.

Composites

Composites are compound expressions composed of others. They represent
the tokens and classes from PEGN and other grammars and are designed to
simplify grammar development at a higher level. Pull requests are
welcome for missing, commonly used composite candidates.

Hooks

Hooks are not strictly an expression type and are declared in the scan
package itself (to avoid a cyclical import dependency since it is passed
a scan.R). A Hook is passed only the scanner struct and must return a bool
indicating if the scan should proceed.  See scan.Hook for more
information.
*/
package is

// ------------------------------- sets -------------------------------

// Seq groups expressions into a sequence. It ensures that all
// expressions appears in that specific order. If any are not the scan
// fails.
type Seq []any

// Lk is a set of positive lookahead expressions. If any are seen at
// the current cursor position the scan will proceed without consuming
// them (unlike is.Opt and is.In). If none are found the scan fails.
type Lk []any

// Not is a set of negative lookahead expressions. If any are seen at
// the current cursor position the scan will fail and the scan is never
// advanced.
type Not []any

// In is a set of advancing expressions. If any expression in the slice
// is found the scan advances to the end of that expression and
// continues. If none of the expressions is found the scan fails.
// Evaluation of expressions is always left to right allowing
// parser developers to prioritize common expressions at the beginning
// of the slice.
type In []any

// Opt is a set of optional advancing expressions. If any expression is
// found the scan is advanced (unlike Lk, which does not advance).
type Opt []any

// To is a set of advancing expressions that mark an exclusive boundary
// at which the scan should stop. The matching expression itself will
// not be advanced.
//
// In order to work with escaped boundaries use a negative
// look-ahead sequence combined with the boundary:
//
//     is.To{s.Seq{is.Not{`\\`,`"`}}}
//
type To []any

// Toi is a set of advancing expressions that mark an inclusive boundary
// after which the scan should stop. The matching expression will be
// included in the advance (unlike is.To).
type Toi []any

// --------------------------- parameterized --------------------------

// MMx is a parameterized advancing expression that matches an inclusive
// minimum and maximum count of the given expression (This). Use within
// is.Lk to disable advancement.
type MMx struct {
	Min  int
	Max  int
	This any
}

// Min is a parameterized advancing expression that matches an inclusive
// minimum number of the given expression item (This). Use within is.Lk
// to disable advancement.
type Min struct {
	Min  int
	This any
}

// Mn1 is shorthand for is.Min{1,This}.
type Mn1 struct{ This any }

// N is a parameterized advancing expression that matches exactly
// N number of the given expression (This). Use within is.Lk to disable
// advancement.
type N struct {
	N    int
	This any
}

// Any is short for is.N{tk.ANY}.
type Any struct {
	N int
}

// Rng is a parameterized advancing expression that matches a single
// Unicode code point (rune, uint32) from an inclusive consecutive set
// from First to Last (First,Last). Use within is.Lk to disable
// advancement.
type Rng struct {
	First rune
	Last  rune
}
