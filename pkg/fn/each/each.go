// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*
Package each shamelessly attempts to bring the better parts of Lisp loops to Go specifically in order to enable rapid, and clean applications development --- particularly when replacing shell scripts with Go.
*/
package each

import (
	"fmt"
	"log"
)

// Do executes the given function for every item in the slice.
func Do[T any](set []T, p func(i T)) {
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
	Do(set, func(i T) { fmt.Println(i) })
}

// Print prints ever element of the set.
func Print[T any](set []T) {
	Do(set, func(i T) { fmt.Print(i) })
}

// Printf calls fmt.Printf on itself with the given form. For more
// substantial printing consider calling each.Do instead.
func Printf[T any](set []T, form string) {
	Do(set, func(i T) { fmt.Printf(form, i) })
}

// Log calls log.Print on ever element of the set.
func Log[T any](set []T) {
	Do(set, func(i T) { log.Print(i) })
}

// Logf calls log.Printf on ever element of the set.
func Logf[T any](set []T, form string) {
	Do(set, func(i T) { log.Printf(form, i) })
}
