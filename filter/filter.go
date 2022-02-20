// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter

import "strings"

// HasPrefix filters the Text input set and returns only those elements
// that have the give prefix.
func HasPrefix[T Text](set []T, pre string) []T {
	m := []T{}
	for _, i := range set {
		if strings.HasPrefix(string(i), pre) {
			m = append(m, i)
		}
	}
	return m
}

// Minus performs a set "minus" operation by returning a new set with
// the elements of the second set removed from it.
func Minus[T Text, M Text](set []T, min []M) []T {
	m := []T{}
	for _, i := range set {
		var seen bool
		for _, n := range min {
			if string(n) == string(i) {
				seen = true
				break
			}
		}
		if !seen {
			m = append(m, i)
		}
	}
	return m
}
