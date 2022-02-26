package filter

import (
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
