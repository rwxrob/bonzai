// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package fn

// Map executes an operator function provided on each item in the
// slice returning a new slice. If error handling is needed it should be
// handled within an enclosure within the function. This keeps
// signatures simple and functional.
func Map[I any, O any](slice []I, f func(in I) O) []O {
	list := []O{}
	for _, i := range slice {
		list = append(list, f(i))
	}
	return list
}
