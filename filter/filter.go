// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"fmt"
	"strings"
)

// Number combines the primitives generally considered numbers by JSON
// and other high-level structure data representations.
type Number interface {
	int16 | int32 | int64 | float32 | float64
}

// Text combines byte slice and string.
type Text interface {
	[]byte | string
}

// P is for "principle" in this case. These are the types that have
// representations in JSON and other high-level structured data
// representations.
type P interface {
	int16 | int32 | int64 | float32 | float64 |
		[]byte | string | bool
}

// Println prints ever element of the set.
func Println[T P](set []T) {
	for _, i := range set {
		fmt.Println(i)
	}
}

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
