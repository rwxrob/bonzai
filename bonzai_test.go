// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z_test

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai"
)

func ExampleArgsFrom() {
	fmt.Printf("%q\n", Z.ArgsFrom(`greet  hi french`))
	fmt.Printf("%q\n", Z.ArgsFrom(`greet hi   french `))
	// Output:
	// ["greet" "hi" "french"]
	// ["greet" "hi" "french" ""]
}

func ExampleArgsOrIn_read_Nil() {

	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin, _ = os.Open(`testdata/in`)

	fmt.Println(Z.ArgsOrIn(nil))

	// Output:
	// some thing
}

func ExampleArgsOrIn_read_Zero_Args() {

	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin, _ = os.Open(`testdata/in`)

	fmt.Println(Z.ArgsOrIn([]string{}))

	// Output:
	// some thing
}

func ExampleArgsOrIn_args_Joined() {

	fmt.Println(Z.ArgsOrIn([]string{"some", "thing"}))

	// Output:
	// some thing
}

func ExampleEsc() {
	fmt.Println(Z.Esc("|&;()<>![]"))
	fmt.Printf("%q", Z.Esc(" \n\r"))
	// Output:
	// \|\&\;\(\)\<\>\!\[\]
	// "\\ \\\n\\\r"
}

func ExampleEscAll() {
	list := []string{"so!me", "<here>", "other&"}
	fmt.Println(Z.EscAll(list))
	// Output:
	// [so\!me \<here\> other\&]
}
