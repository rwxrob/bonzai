// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter

import "fmt"

// Println prints ever element of the set.
func Println[T P](set []T) {
	for _, i := range set {
		fmt.Println(i)
	}
}
