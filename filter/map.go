// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter

import (
	"fmt"
	"sort"
)

// Println prints ever element of the set.
func Println[T P](set []T) {
	for _, i := range set {
		fmt.Println(i)
	}
}

// Keys returns the keys in lexicographically sorted order.
func Keys[T any](m map[string]T) []string {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
		sort.Strings(keys)
	}
	return keys
}
