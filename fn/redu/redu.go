package redu

import (
	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/to"
)

type Strings fn.Reducer[string]
type Ints fn.Reducer[int]
type Anys fn.Reducer[any]

// Longest will convert everything [to.String] and return the length of
// the longest string in the set.
func Longest[T any](set []T) int {
	var longest int
	for _, v := range set {
		s := to.String(v)
		if len(s) > longest {
			longest = len(s)
		}
	}
	return longest
}

// Unique removed duplicates.
func Unique[T comparable](set []T) []T {
	var list []T
	seen := map[T]bool{}
	for _, v := range set {
		if _, has := seen[v]; !has {
			list = append(list, v)
			seen[v] = true
		}
	}
	return list
}
