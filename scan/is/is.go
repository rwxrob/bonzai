// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package is contains the scannable structs recognized as a pseudo-grammar
accepted by the scan.Expect and scan.Check methods. These structs are
guaranteed never to have their structure order change in any way
allowing them to be used in key-less, in-line composable syntax notation
(which, despite the many editor warnings, is entirely supported by Go
and always will be).
*/
package is

// Not scannable slice represents set of negative lookahead expressions.
// If any are seen at the current cursor position the scan will fail.
type Not []any

// In scannable slice represents a set group of scannables. If any
// scannable in the slice is found the Is scannable itself is true. If
// nothing is found the result is false. The search through the In
// scannable group is always linear allowing parser developers to
// establish the priority for common finds themselves.
type In []any

// Seq scannable slice represents a sequence group of scannables. It
// ensures that the slice of scannables always appears in that specific
// order. Avoid over use of Seq since Expect/Check already expect
// a sequence of scannables.
type Seq []any

// Opt scannable slice represents a set of optional positive look-ahead
// scannables much like In, except that if nothing is found no error is
// generated. This is the equivalent of the question mark on a set
// ([]?) from regular expressions.
type Opt []any

// MMx scannable struct represents the inclusive minimum and maximum
// count of the give scannable item (This).
type MMx struct {
	Min  int
	Max  int
	This any
}

// Min scannable struct represents the inclusive minimum (Min,) count
// of the given scannable item (This).
type Min struct {
	Min  int
	This any
}

// N scannable struct represents exactly X number of the given scannable
// items (This).
type N struct {
	N    int
	This any
}

// Rng scannable struct represents any single rune from an
// inclusive consecutive set from the First to Last (First,Last).
type Rng struct {
	First rune
	Last  rune
}

// ---------------------------- composites ----------------------------
//                    (keep most common to the left)

var WS = In{' ', '\n', '\t', '\r'}
var Digit = Rng{0, 9}
