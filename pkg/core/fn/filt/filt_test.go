// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package filt_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/core/fn/each"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
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
	each.Println(filt.BaseHasPrefix(paths, "f"))
	// Output:
	// some/foo
	// some/foo1
}

func ExampleNotEmpty() {
	set := []string{
		"one", "", "two", "", "three",
	}
	fmt.Println(filt.NotEmpty(set))
	// Output:
	// [one two three]
}

func ExampleRemoveIndex() {

	type SomeType struct {
		Thing string
	}

	one := &SomeType{`one`}
	two := &SomeType{`two`}
	three := &SomeType{`three`}

	set := []*SomeType{one, two, three}
	oneref := set[0]
	threeref := set[2]

	nset := filt.RemoveIndex(set, 1)

	fmt.Println(nset[0] == oneref)
	fmt.Println(nset[1] == threeref)

	// Output:
	// true
	// true
}
