package qstack

import (
	"unicode"
)

// Fields is exactly the same as strings.Fields (returning strings
// separated by unicode.IsSpace) except it returns a qstack.QS instead.
func Fields(in string) *QS[string] {
	var field string
	fields := New[string]()
	for _, r := range []rune(in) {
		if unicode.IsSpace(r) {
			if len(field) > 0 {
				fields.Push(field)
				field = ""
			}
			continue
		}
		field += string(r)
	}
	if len(field) > 0 {
		fields.Push(field)
	}
	return fields
}
