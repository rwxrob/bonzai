// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filt_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/filt"
	"github.com/rwxrob/bonzai/loop"
)

func ExampleHasPrefix() {
	set := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(filt.HasPrefix(set, "t"))
	// Output:
	// [two three]
}

func ExampleBaseHasPrefix() {
	paths := []string{
		"some/foo",
		"some/foo1",
		"some/",
		"some/blah",
	}
	loop.Println(filt.BaseHasPrefix(paths, "f"))
	// Output:
	// some/foo
	// some/foo1
}
