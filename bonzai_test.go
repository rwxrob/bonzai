// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
)

func ExampleArgsFrom() {
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet  hi french`))
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet hi   french `))
	// Output:
	// ["greet" "hi" "french"]
	// ["greet" "hi" "french" " "]
}
