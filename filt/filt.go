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
	m := []T{}
	for _, i := range set {
		if strings.HasPrefix(string(i), pre) {
			m = append(m, i)
		}
	}
	return m
}

// BaseHasPrefix filters the input of file paths and returns only those
// elements where the base name has the given prefix.
func BaseHasPrefix[T Text](paths []T, pre string) []T {
	m := []T{}
	for _, i := range paths {
		if strings.HasPrefix(filepath.Base(string(i)), pre) {
			m = append(m, i)
		}
	}
	return m
}
