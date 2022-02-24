// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/filter"
)

func ExamplePrintln() {
	set := []string{"doe", "ray", "mi"}
	filter.Println(set)
	bools := []bool{false, true, true}
	filter.Println(bools)
	// Output:
	// doe
	// ray
	// mi
	// false
	// true
	// true
}

func ExampleKeys() {
	m1 := map[string]int{"two": 2, "three": 3, "one": 1}
	m2 := map[string]string{"two": "two", "three": "three", "one": "one"}
	fmt.Println(filter.Keys(m1))
	fmt.Println(filter.Keys(m2))
	// Output:
	// [one three two]
	// [one three two]
}
