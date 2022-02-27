// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

/*

Package fn implements the traditional map/filter/reduce/each functions
and an array type (A) for those who prefer a more object-oriented
approach. Unlike other implementations, the array (slice) is always
first preventing the first-class in-line anonymous function from
obfuscating the parameter list of the functional function.

"Why have a functional package in a commander?"

These functions are included because commander frameworks are about
creating less friction during command applications development. Since
many such commands are creating from ported shell scripts where the UNIX
philosophy and filters reign supreme in makes sense to enable similar
functional approaches (filters, pipelines) to approximate the simplicity
of and speed shell scripting development (at the very slight cost of
runtime performance). In short, this package speeds and simplifies
command application development by making it more compatible with shell
scripting in general.

*/
package fn

import "github.com/rwxrob/bonzai/each"

// A is the equivalent of an array primitive in other functional
// languages. It is a generic slice of anything.
type A[T any] []T

func (a A[any]) Each(f func(i any))                { each.Do(a, f) }
func (a A[any]) E(f func(i any))                   { each.Do(a, f) }
func (a A[any]) Map(f func(i any) any) A[any]      { return Map(a, f) }
func (a A[any]) M(f func(i any) any) A[any]        { return Map(a, f) }
func (a A[any]) Filter(f func(i any) bool) A[any]  { return Filter(a, f) }
func (a A[any]) F(f func(i any) bool) A[any]       { return Filter(a, f) }
func (a A[any]) Reduce(f func(i any, r *any)) *any { return Reduce(a, f) }
func (a A[any]) R(f func(i any, r *any)) *any      { return Reduce(a, f) }
func (a A[any]) Print()                            { each.Print(a) }
func (a A[any]) Println()                          { each.Println(a) }
func (a A[any]) P()                                { each.Println(a) }
func (a A[any]) Printf(t string)                   { each.Printf(a, t) }
func (a A[any]) Log()                              { each.Log(a) }
func (a A[any]) Logf(f string)                     { each.Logf(a, f) }

// Map executes an operator function provided on each item in the slice
// returning a new slice with items of a potentially different type
// completely (which is different from using the Array.Map method which
// requires returning the same type). If error handling is needed it
// should be handled within an enclosure within the function. This keeps
// signatures simple and functional.
func Map[I any, O any](slice []I, f func(in I) O) []O {
	list := []O{}
	for _, i := range slice {
		list = append(list, f(i))
	}
	return list
}

// Filter applies the boolean function on each item only returning those
// items that evaluate to true.
func Filter[T any](slice []T, f func(in T) bool) []T {
	list := []T{}
	for _, i := range slice {
		if f(i) {
			list = append(list, i)
		}
	}
	return list
}

// Reduce calls the given reducer function for every item in the slice
// passing a required reference to an item to hold the results If error
// handling is needed it should be handled within an enclosure within
// the function.  This keeps signatures simple and functional.
func Reduce[T any, R any](slice []T, f func(in T, ref *R)) *R {
	r := new(R)
	for _, i := range slice {
		f(i, r)
	}
	return r
}
