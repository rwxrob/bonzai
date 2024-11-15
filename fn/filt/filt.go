// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filt

import (
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/fn"
)

type Strings fn.Filterer[string]
type Ints fn.Filterer[int]
type Anys fn.Filterer[any]

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

// HasPrefixSorted filters the Text input set and returns only those elements
// that have the give prefix assuming that the input is already sorted
// so that early return on the first miss can be trusted.
func HasPrefixSorted[T Text](set []T, pre string) []T {
	var s []T
	for _, i := range set {
		if !strings.HasPrefix(string(i), pre) {
			return s
		}
		s = append(s, i)
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

// HasSuffix filters the Text input set and returns only those elements
// that have the give suffix.
func HasSuffix[T Text](set []T, suf string) []T {
	var s []T
	for _, i := range set {
		if strings.HasSuffix(string(i), suf) {
			s = append(s, i)
		}
	}
	return s
}

// HasSuffixSorted filters the Text input set and returns only those elements
// that have the give suffix assuming that the input is already sorted
// so that early return on the first miss can be trusted.
func HasSuffixSorted[T Text](set []T, pre string) []T {
	var s []T
	for _, i := range set {
		if !strings.HasPrefix(string(i), pre) {
			return s
		}
		s = append(s, i)
	}
	return s
}

// BaseHasSuffix filters the input of file paths and returns only those
// elements where the base name has the given suffix.
func BaseHasSuffix[T Text](paths []T, pre string) []T {
	var s []T
	for _, i := range paths {
		if strings.HasSuffix(filepath.Base(string(i)), pre) {
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
