// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package maps_test

import (
	"os"

	"github.com/rwxrob/bonzai/pkg/fn/each"
	"github.com/rwxrob/bonzai/pkg/fn/maps"
)

func ExampleMarkDirs() {
	entries, _ := os.ReadDir("testdata/markdirs")
	each.Println(maps.MarkDirs(entries))
	//Output:
	// dir1/
	// file1
}

func ExampleBase() {
	paths := []string{
		`some/thing /here`,
		`other/thing`,
		`foo`,
	}
	each.Println(maps.Base(paths))
	//Output:
	// here
	// thing
	// foo
}

func ExampleHashComment() {
	each.Println(maps.HashComment([]string{"foo", "bar"}))
	// Output:
	// # foo
	// # bar
}

func ExampleEscSpace() {
	each.Println(maps.EscSpace([]string{"some thing", "one other   thing"}))
	// Output:
	// some\ thing
	// one\ other\ \ \ thing
}

func ExampleTrimSpace() {
	each.Println(maps.TrimSpace([]string{
		"  some thing  ",
		" one other   thing",
		"ing-ing  ",
		"",
	}))
	// Output:
	// some thing
	// one other   thing
	// ing-ing
	//
}
