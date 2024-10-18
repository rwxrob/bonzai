// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package mapf_test

import (
	"os"

	"github.com/rwxrob/bonzai/pkg/fn"
	"github.com/rwxrob/bonzai/pkg/fn/each"
	"github.com/rwxrob/bonzai/pkg/fn/mapf"
)

func ExampleMarkDirs() {
	entries, _ := os.ReadDir("testdata/markdirs")
	each.Println(fn.Map(entries, mapf.MarkDirs))
	//Output:
	// dir1/
	// file1
}

func ExampleHashComment() {
	each.Println(fn.Map([]string{"foo", "bar"}, mapf.HashComment))
	// Output:
	// # foo
	// # bar
}

func ExampleEscSpace() {
	s := []string{"one here", "and another    one"}
	each.Println(fn.Map(s, mapf.EscSpace))
	// Output:
	// one\ here
	// and\ another\ \ \ \ one
}
