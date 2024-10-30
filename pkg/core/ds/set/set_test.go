// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package set_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/core/ds/set"
)

func ExampleMinus() {
	s := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	m := []string{"two", "four", "six"}
	fmt.Println(set.Minus(s, m))
	// Output:
	// [one three five seven]
}
