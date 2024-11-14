// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package fn

import (
	"fmt"
	"log"

	"github.com/rwxrob/bonzai/fn/each"
)

type Transformer[I any, O any] interface {
	Transform(in []I) []O
}

type Filterer[T any] interface {
	Filter(in []T) []T
}

type Reducer[T any] interface {
	Reduce(in []T) []T
}

// Number combines the primitives generally considered numbers by JSON
// and other high-level structure data representations.
type Number interface {
	int | int64 | int32 | int16 | int8 |
		uint64 | uint32 | uint16 | uint8 |
		float64 | float32
}

// Text combines byte slice and string.
type Text interface {
	[]byte | string | []rune
}

// Sharable are the types that have representations in JSON, YAML, TOML
// and other high-level structured data representations.
type Sharable interface {
	int | int64 | int32 | int16 | int8 |
		uint64 | uint32 | uint16 | uint8 |
		float64 | float32 |
		[]byte | string |
		bool
}

// A is the equivalent of an array primitive in other functional
// languages. It is a generic slice of anything.
type A[T any] []T

// Each calls each.Do on self
func (a A[any]) Each(f func(i any)) { each.Do(a, f) }

// E calls Do function on self.
func (a A[any]) E(f func(i any)) { each.Do(a, f) }

// Map calls Map function on self.
func (a A[any]) Map(f func(i any) any) A[any] { return Map(a, f) }

// M calls Map function on self.
func (a A[any]) M(f func(i any) any) A[any] { return Map(a, f) }

// Filter calls Filter function on self.
func (a A[any]) Filter(f func(i any) bool) A[any] { return Filter(a, f) }

// Filter calls Filter function on self.
func (a A[any]) F(f func(i any) bool) A[any] { return Filter(a, f) }

// Reduce calls Reduce function on self.
func (a A[any]) Reduce(f func(i any, r *any)) *any { return Reduce(a, f) }

// Reduce calls Reduce function on self.
func (a A[any]) R(f func(i any, r *any)) *any { return Reduce(a, f) }

// Print calls each.Print on self.
func (a A[any]) Print() { each.Print(a) }

// Println calls each.Println on self.
func (a A[any]) Println() { each.Println(a) }

// Printf calls each.Printf on self.
func (a A[any]) Printf(t string) { each.Printf(a, t) }

// Log calls each.Log on self.
func (a A[any]) Log() { each.Log(a) }

// Logf calls each.Logf on self.
func (a A[any]) Logf(f string) { each.Logf(a, f) }

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

// Pipe implements the closest thing to UNIX pipelines possible in Go by
// passing each argument to the next assuming a func(in any) any format
// where the input (in) is converted to a string (if not already
// a string). If any return an error the pipeline returns an empty
// string and logs an error.
func Pipe(filter ...any) string {
	if len(filter) == 0 {
		return ""
	}
	var in any
	in = filter[0]
	for _, f := range filter[1:] {
		switch v := f.(type) {
		case func(any) any:
			in = v(in)
		default:
			in = f
		}
		if err, iserr := in.(error); iserr {
			log.Print(err)
			return ""
		}
	}
	return fmt.Sprintf("%v", in)
}

// PipePrint prints the output (and a newline) of a Pipe logging any
// errors encountered.
func PipePrint(filter ...any) { fmt.Println(Pipe(filter...)) }

// Or returns the first non-zero value of the two provided. If this is
// the zero value for type [T], it returns that; otherwise, it returns
// this.
func Or[T comparable](this, that T) T {
	if this == *new(T) {
		return that
	}
	return this
}

// Fall iterates over a list of values, returning the first non-zero
// value. If [vals] contains only one element, it returns that value. It
// recursively calls [Or] to compare each element with the subsequent
// elements until it finds a non-zero value or reaches the end of the
// list.
func Fall[T comparable](vals ...T) T {
	if len(vals) == 1 {
		return vals[0]
	}
	return Or(vals[0], Fall(vals[1:]...))
}
