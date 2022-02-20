// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filter_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/filter"
)

func ExampleHasPrefix() {
	set := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(filter.HasPrefix(set, "t"))
	// Output:
	// [two three]
}

func ExampleMinus() {
	set := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(filter.Minus(set, []string{"two", "four", "six"}))
	// Output:
	// [one three five seven]
}
