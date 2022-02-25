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
	fmt.Println(comp.File(nil, "b"))
	fmt.Println(comp.File(nil, "blah"))
	fmt.Println(comp.File(nil, "blah/"))
	//Output:
	// [bar/ blah/ come/ foo/ other/]
	// [bar/ blah/]
	// [blah/]
	// [dir1/ dir2/ file1 file2]
}
