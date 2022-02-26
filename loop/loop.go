/*
Package loop shamelessly attempts to bring the better parts of Lisp loops to Go specifically in order to enable rapid, and clean applications development --- particularly when replacing shell scripts with Go.
*/
package loop

import "fmt"

// All executes the given function for every item in the slice.
func All[T any](set []T, p func(i T)) {
	for _, i := range set {
		p(i)
	}
}

// UntilError executes the give function for every item in the slice
// until it encounters and error and returns the error, if any.
func UntilError[T any](set []T, p func(i T) error) error {
	for _, i := range set {
		if err := p(i); err != nil {
			return err
		}
	}
	return nil
}

// Println prints ever element of the set.
func Println[T any](set []T) {
	All(set, func(i T) { fmt.Println(i) })
}

// Print prints ever element of the set.
func Print[T any](set []T) {
	All(set, func(i T) { fmt.Print(i) })
}

// Printf prints ever element of the set using format string.
func Printf[T any](set []T, form string) {
	All(set, func(i T) { fmt.Printf(form, i) })
}
