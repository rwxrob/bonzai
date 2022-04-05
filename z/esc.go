package Z

import "github.com/rwxrob/fn"

// EscThese is set to the default UNIX shell characters which require
// escaping to be used safely on the terminal. It can be changed to suit
// the needs of different host shell environments.
var EscThese = " \r\t\n|&;()<>![]"

// Esc returns a shell-escaped version of the string s. The returned value
// is a string that can safely be used as one token in a shell command line.
func Esc(s string) string {
	var buf []rune
	for _, r := range s {
		for _, esc := range EscThese {
			if r == esc {
				buf = append(buf, '\\')
			}
		}
		buf = append(buf, r)
	}
	return string(buf)
}

// EscAll calls Esc on all passed strings.
func EscAll(args []string) []string { return fn.Map(args, Esc) }
