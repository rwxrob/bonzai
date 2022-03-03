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
simplify grammar development at a higher level. Pull requests are welcome for missing, commonly used composite candidates.
*/
package is

// ------------------------------- sets -------------------------------

// Seq groups expressions into a sequence. It ensures that all
// expressions appears in that specific order. If any are not the scan
// fails.
type Seq []any

// Lk is a set of positive lookahead expressions. If any are seen at
// the current cursor position the scan will proceed without consuming
// them (unlike is.Opt and is.Any). If none are found the scan fails.
type Lk []any

// Not is a set of negative lookahead expressions. If any are seen at
// the current cursor position the scan will fail. Otherwise, the scan
// proceeds from that same position.
type Not []any

// Any is a set of advancing expressions. If any scannable in the slice
// is found the scan advances to the end of that expression and
// continues. If none of the expressions is found the scan fails.
// Evaluation of expressions is always left to right allowing allowing
// parser developers to prioritize common expressions at the beginning
// of the slice.
type Any []any

// Opt is a set of optional advancing expressions. If any expression is
// found the scan is advanced (unlike Lk, which does not advance).
type Opt []any

// --------------------------- parameterized --------------------------

// MMx parameterized advancing expression scans for the inclusive
// minimum and maximum count of the given expression (This). Use within
// is.Lk to disable advancement.
type MMx struct {
	Min  int
	Max  int
	This any
}

// Min parameterized advancing expression scans for the inclusive minimum
// number of the given expression item (This). Use within is.Lk to
// disable advancement.
type Min struct {
	Min  int
	This any
}

// Mn1 parameterized advancing expression is shorthand for is.Min{1,This}.
type Mn1 struct{ This any }

// N parameterized advancing expression scans for exactly N number of
// the given expression (This). Use within is.Lk to disable advancement.
type N struct {
	N    int
	This any
}

// Rng parameterized advancing expression scans for a single Unicode
// code point (rune, uint32) from an inclusive consecutive set from
// First to Last (First,Last). Use within is.Lk to disable advancement.
type Rng struct {
	First rune
	Last  rune
}

// ---------------------------- composites ----------------------------
//                    (keep most common to the left)

var WS = In{' ', '\n', '\t', '\r'}
var Digit = Rng{0, 9}
