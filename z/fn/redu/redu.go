package redu

import "github.com/rwxrob/bonzai/z/to"

// Longest will convert everything to.String and return the length of
// the longest string in the set
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
