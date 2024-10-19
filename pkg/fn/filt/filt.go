// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filt

import (
	"path/filepath"
	"strings"
)

type Text interface {
	string | []byte
}

// HasPrefix filters the Text input set and returns only those elements
// that have the give prefix.
func HasPrefix[T Text](set []T, pre string) []T {
	var s []T
	for _, i := range set {
		if strings.HasPrefix(string(i), pre) {
			s = append(s, i)
		}
	}
	return s
}

// BaseHasPrefix filters the input of file paths and returns only those
// elements where the base name has the given prefix.
func BaseHasPrefix[T Text](paths []T, pre string) []T {
	var s []T
	for _, i := range paths {
		if strings.HasPrefix(filepath.Base(string(i)), pre) {
			s = append(s, i)
		}
	}
	return s
}

// NotEmpty filters only strings that are not empty.
func NotEmpty[T Text](set []T) []T {
	var s []T
	for _, i := range set {
		if string(i) != "" {
			s = append(s, i)
		}
	}
	return s
}

// RemoveIndex removes the item at the given index returning a new slice
// while preserving the references to each item in the original slice.
func RemoveIndex[T any](set []T, pos int) []T {
	nset := make([]T, (len(set) - 1))
	k := 0
	for i, it := range set {
		if i == pos {
			continue
		}
		nset[k] = it
		k++
	}
	return nset
}
