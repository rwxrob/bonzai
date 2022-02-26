/*
Package loop shamelessly attempts to bring the better parts of Lisp loops to Go specifically in order to enable rapid, and clean applications development --- particularly when replacing shell scripts with Go.
*/
package loop

import "fmt"

// Do executes the given function for each item in the slice. If any
// error is encountered processing stops and error returned.
func Do[T any](set []T, p func(i T) error) error {
	for _, i := range set {
		if err := p(i); err != nil {
			return err
		}
	}
	return nil
}

// Println prints ever element of the set.
func Println[T any](set []T) {
	Do(set, func(i T) error { fmt.Println(i); return nil })
}

// Print prints ever element of the set.
func Print[T any](set []T) {
	Do(set, func(i T) error { fmt.Print(i); return nil })
}

// Printf prints ever element of the set using format string.
func Printf[T any](set []T, form string) {
	Do(set, func(i T) error { fmt.Printf(form, i); return nil })
}
