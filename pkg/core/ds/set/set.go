// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package set provides structs, functions, and methods for dealing with
sets of anything convert to its string notation with to.String.
*/
package set

import "github.com/rwxrob/bonzai/pkg/core/to"

// Minus performs a set "minus" operation by returning a new set with
// the elements of the second set removed from it. Any type type may be
// passed but the string value returned from to.String will be used.
// This operation is not highly performant, but fast enough for
// convenience.
func Minus[S any, T any](set []S, min []T) []string {
	m := []string{}
	for _, i := range set {
		x := to.String(i)
		var seen bool
		for _, n := range min {
			if x == to.String(n) {
				seen = true
				break
			}
		}
		if !seen {
			m = append(m, x)
		}
	}
	return m
}
