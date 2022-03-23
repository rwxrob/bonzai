// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/util"
	"github.com/rwxrob/fn/each"
)

func ExampleArgsFrom() {
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet  hi french`))
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet hi   french `))
	// Output:
	// ["greet" "hi" "french"]
	// ["greet" "hi" "french" ""]
}

func ExampleFiles() {
	each.Println(util.Files("testdata/files"))
	// Output:
	// testdata/files/bar
	// testdata/files/blah
	// testdata/files/dir1/
	// testdata/files/foo
	// testdata/files/other
	// testdata/files/some
}

func ExampleFiles_spaces() {
	each.Println(util.Files("testdata/files/dir1"))
	// Output:
	// testdata/files/dir1/some\ thing
}

func ExampleFiles_empty() {
	os.Chdir("testdata/files")
	defer os.Chdir("../..")
	each.Println(util.Files(""))
	// Output:
	// bar
	// blah
	// dir1/
	// foo
	// other
	// some
}

func ExampleFiles_not_Directory() {
	fmt.Println(util.Files("none"))
	// Output:
	// []
}
