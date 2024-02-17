// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z_test

import (
	"fmt"
	"os"

	Z "github.com/vekio/bonzai/z"
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
	// input
	//
	// with blank line
}

func ExampleArgsOrIn_read_Zero_Args() {

	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin, _ = os.Open(`testdata/in`)

	fmt.Println(Z.ArgsOrIn([]string{}))

	// Output:
	// some thing
	// input
	//
	// with blank line
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

func ExampleInferredUsage_optional_Param() {
	x := &Z.Cmd{
		Params: []string{"p1", "p2"},
		Call:   func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// (p1|p2)?
}

func ExampleInferredUsage_min_One_Param() {
	x := &Z.Cmd{
		Params:  []string{"p1", "p2"},
		MinParm: 1,
		Call:    func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// (p1|p2)+
}

func ExampleInferredUsage_min_3_Param() {
	x := &Z.Cmd{
		Params:  []string{"p1", "p2"},
		MinParm: 3,
		Call:    func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// (p1|p2){3,}
}

func ExampleInferredUsage_commands() {
	x := &Z.Cmd{
		Commands: []*Z.Cmd{
			&Z.Cmd{Name: "foo", Aliases: []string{"f"}},
			&Z.Cmd{Name: "bar"},
		},
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// ((f|foo)|bar)
}

func ExampleInferredUsage_commands_and_Params() {
	x := &Z.Cmd{
		Params: []string{"p1", "p2"},
		Commands: []*Z.Cmd{
			&Z.Cmd{Name: "foo", Aliases: []string{"f"}},
			&Z.Cmd{Name: "bar"},
		},
		Call: func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// ((p1|p2)?|((f|foo)|bar))
}

func ExampleInferredUsage_error_No_Call_or_Command() {
	x := &Z.Cmd{
		Params: []string{"p1", "p2"},
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// {ERROR: neither Call nor Commands defined}
}

func ExampleInferredUsage_error_Params_without_Call() {
	x := &Z.Cmd{
		Params: []string{"p1", "p2"},
		Commands: []*Z.Cmd{
			&Z.Cmd{Name: "foo", Aliases: []string{"f"}},
			&Z.Cmd{Name: "bar"},
		},
	}
	fmt.Println(Z.InferredUsage(x))
	// Output:
	// {ERROR: Params without Call: p1, p2}
}

func ExampleUsageGroup() {
	fmt.Println(Z.UsageGroup([]string{"", "foo", "", "bar", "with space"}, 1, 1))
	fmt.Printf("%q\n", Z.UsageGroup([]string{"", ""}, 1, 1))
	fmt.Println(Z.UsageGroup([]string{"one"}, 1, 1))
	// Output:
	// (foo|bar|with space)
	// ""
	// one
}
