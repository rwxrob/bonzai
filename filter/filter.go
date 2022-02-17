package filter

import "strings"

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
	Number
	Text
	bool
}

// HasPrefix filters the Text input set and returns only those elements
// that have the give prefix.
func HasPrefix[T Text](set []T, pre string) []T {
	m := []T{}
	for _, s := range set {
		if strings.HasPrefix(string(s), pre) {
			m = append(m, s)
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
