// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package z (often imported as "is") defines the Bonzai Parsing Expression
Grammar Notation (BPEGN) (aka "Bonzai Scanner Expect" language)
implemented entirely using Go types (mostly slices and structs).  BPEGN
can be 100% transpiled to and from the Parsing Expression Grammer
Notation (PEGN). Code in one grammar and use the bonzai command to
easily generate the other.

Nodes and Expressions

Nodes (z.N) indicate something to be captured as a part of the resulting
abstract syntax tree. They are functionally equivalent to parenthesis in
regular expressions but with the obvious advantage of capturing a rooted
node tree instead of an array. Expressions (z.X) indicate a sequence to be scanned but not captured unless the expression itself is within a node (z.N).

Tokens

The BPEGN token strings are contained in the "tk" package can be used as
is. See the "tk" package for details.

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
package z

// ------------------------------- core -------------------------------

// N ("node") is a named sequence of expressions that will be captured
// into a node. The first string must always be the name which can be
// any valid Go string. If any expression fails to match the scan fails.
// Otherwise, a new tree.Node is added under the current node and the
// scan proceeds. Nodes must either contain other nodes or no nodes at
// all. If the first item in the sequence after the name is not also
// a node (z.N) then the node is marked as "edge" (or "leaf") and any
// nodes detected further in the sequence will cause the scan to fail
// with a syntax error.
type N []any

// X ("expression") is a sequence of expressions.  If any are not the
// scan fails. (Equal to (?foo) in regular expressions.)
type X []any

// ------------------------------- sets -------------------------------

// It (the slice) is a set of positive lookahead expressions. If any are
// seen at the current cursor position the scan will proceed without
// consuming them (unlike is.O and is.In). If none are found the scan
// fails. This is useful when everything from one expression is wanted
// except for a few positive exceptions. (Equal to ampersand (&) in
// PEGN.) Also see the tk.IS token.
type It []any

// Not (the slice) is a set of negative lookahead expressions. If any
// are seen at the current cursor position the scan will fail and the
// scan is never advanced. This is useful when everything from one
// expression is wanted except for a few negative exceptions. (Equal to
// exclamation point (!) in PEGN.) Also see the tk.NOT token.
type Not []any

// In is a set of advancing expressions. If any expression in the slice
// is found the scan advances to the end of that expression and
// continues. If none of the expressions is found the scan fails.
// Evaluation of expressions is always left to right allowing
// parser developers to prioritize common expressions at the beginning
// of the slice.
type In []any

// O is a set of optional advancing expressions. If any expression is
// found the scan is advanced (unlike is.It, which does not advance).
type O []any

// To is a set of advancing expressions that mark an exclusive boundary
// at which the scan should stop returning a cursor just before the
// boundary.
type To []any

// Ti ("to inclusive") is an inclusive version of is.To returning
// a cursor that points to the last rune of the boundary itself.
type Ti []any

// --------------------------- parameterized --------------------------

// MM ("minmax") is a parameterized advancing expression that matches an
// inclusive minimum and maximum count of the given expression (This).
type MM struct {
	Min  int
	Max  int
	This any
}

// M ("min") is a parameterized advancing expression that matches an
// inclusive minimum number of the given expression item (This). Use
// within is.It to disable advancement.
type M struct {
	Min  int
	This any
}

// M1 is shorthand for z.M{1,This}.
type M1 struct{ This any }

// C is a parameterized advancing expression that matches an exact count
// of the given expression (This). Use within is.It to disable
// advancement.
type C struct {
	N    int
	This any
}

// Any is short for is.C{tk.ANY}.
type Any struct {
	N int
}

// Rng is a parameterized advancing expression that matches a single
// Unicode code point (rune, uint32) from an inclusive consecutive set
// from First to Last (First,Last). Use within is.It to disable
// advancement.
type Rng struct {
	First rune
	Last  rune
}
