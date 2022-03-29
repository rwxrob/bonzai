// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai"
)

func ExampleArgsFrom() {
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet  hi french`))
	fmt.Printf("%q\n", bonzai.ArgsFrom(`greet hi   french `))
	// Output:
	// ["greet" "hi" "french"]
	// ["greet" "hi" "french" ""]
}

func ExampleArgsOrIn_read_Nil() {

	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin, _ = os.Open(`testdata/in`)

	fmt.Println(bonzai.ArgsOrIn(nil))

	// Output:
	// some thing
}

func ExampleArgsOrIn_read_Zero_Args() {

	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin, _ = os.Open(`testdata/in`)

	fmt.Println(bonzai.ArgsOrIn([]string{}))

	// Output:
	// some thing
}

func ExampleArgsOrIn_args_Joined() {

	fmt.Println(bonzai.ArgsOrIn([]string{"some", "thing"}))

	// Output:
	// some thing
}
