// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package set provides structs, functions, and methods for dealing with
sets.
*/
package set

import "github.com/rwxrob/bonzai/pkg/to"

// MinusAsString performs a set "minus" operation by returning a new set with
// the elements of the second set removed from it.
func MinusAsString[T any, M any](set []T, min []T) []string {
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
