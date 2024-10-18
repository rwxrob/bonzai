// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*

Package set provides structs, functions, and methods for dealing with
sets composed of strings (or things that can be converted to strings).

*/
package set

import "fmt"

// String converts whatever is passed to its fmt.Sprintf("%v") string
// version (but avoids calling it if possible). Be sure you use things
// with consistent string representations.
func String(in any) string {
	switch v := in.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case []rune:
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Minus performs a set "minus" operation by returning a new set with
// the elements of the second set removed from it.
func Minus[T any, M any](set []T, min []T) []string {
	m := []string{}
	for _, i := range set {
		x := String(i)
		var seen bool
		for _, n := range min {
			if x == String(n) {
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
