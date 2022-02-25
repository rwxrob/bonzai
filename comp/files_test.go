// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleFiles() {
	os.Chdir("testdata/files")
	defer os.Chdir("../..")
	fmt.Println(comp.Files(nil))
	//Output:
	// [bar blah come foo other]
}
