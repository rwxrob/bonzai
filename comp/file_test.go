// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleFile() {
	os.Chdir("testdata/file")
	defer os.Chdir("../..")

	fmt.Println(comp.File(nil))
	fmt.Println(comp.File(nil, ""))
	fmt.Println(comp.File(nil, "fo"))
	fmt.Println(comp.File(nil, "foo"))
	fmt.Println(comp.File(nil, "bar"))
	fmt.Println(comp.File(nil, "bar/"))
	fmt.Println(comp.File(nil, "bar/fo"))
	fmt.Println(comp.File(nil, "bar/foo"))
	fmt.Println(comp.File(nil, "com"))
	fmt.Println(comp.File(nil, "come/"))
	fmt.Println(comp.File(nil, "b"))
	fmt.Println(comp.File(nil, "blah"))
	fmt.Println(comp.File(nil, "blah/"))
	fmt.Println(comp.File(nil, "blah/f"))
	fmt.Println(comp.File(nil, "blah/file1"))
	fmt.Println(comp.File(nil, "blah/file1", "blah/file1"))

	//Output:
	// [bar blah come foo foo.go other]
	// [bar blah come foo foo.go other]
	// [foo foo.go]
	// [foo foo.go]
	// [bar/foo bar/foo.go bar/other]
	// [bar/foo bar/foo.go bar/other]
	// [bar/foo bar/foo.go]
	// [bar/foo bar/foo.go]
	// [come/one]
	// [come/one]
	// [bar blah]
	// [blah/dir1 blah/dir2 blah/file1 blah/file2]
	// [blah/dir1 blah/dir2 blah/file1 blah/file2]
	// [blah/file1 blah/file2]
	// [blah/file1]
	// []
}
