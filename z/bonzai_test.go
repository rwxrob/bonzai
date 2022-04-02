// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z_test

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/term"
	"github.com/rwxrob/term/esc"
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

/*
func ExampleInferredUsage() {

	// call method, params, aliases, commands, and hidden commands
	// TODO

	// [(p1|p2)|(h|help)]

	// Call with optional params and also Commands
	x := &Z.Cmd{
		Params:   []string{"p1", "p2"},
		Commands: []*Z.Cmd{&Z.Cmd{Name: "foo"}, &Z.Cmd{Name: "bar"}},
		Call:     func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))

	// no Call, Command required
	x = &Z.Cmd{
		Commands: []*Z.Cmd{&Z.Cmd{Name: "foo"}, &Z.Cmd{Name: "bar"}},
	}
	fmt.Println(Z.InferredUsage(x))

	// no Call with unused params (ERROR)
	x = &Z.Cmd{
		Params: []string{"p1", "p2"},
	}
	fmt.Println(Z.InferredUsage(x))

	// Call with optional Params, but no commands
	x = &Z.Cmd{
		Params: []string{"p1", "p2"},
		Call:   func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))

	// Call with optional Commands, but no params
	x = &Z.Cmd{
		Commands: []*Z.Cmd{&Z.Cmd{Name: "foo"}, &Z.Cmd{Name: "bar"}},
		Call:     func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	fmt.Println(Z.InferredUsage(x))

	// no Call, Commands, or Params (ERROR)
	x = &Z.Cmd{}
	fmt.Println(Z.InferredUsage(x))

	// Output:
	// [p1|p2|foo|bar]
	// (foo|bar)
	// {ERROR: Params without Call: p1, p2}
	// [p1|p2]
	// [foo|bar]
	// {ERROR: neither Call nor Commands defined}

}
*/

func ExampleEmphasize_emphForLess() {

	/*
	   export LESS_TERMCAP_mb="[35m" # magenta
	   export LESS_TERMCAP_md="[33m" # yellow
	   export LESS_TERMCAP_me="" # "0m"
	   export LESS_TERMCAP_se="" # "0m"
	   export LESS_TERMCAP_so="[34m" # blue
	   export LESS_TERMCAP_ue="" # "0m"
	   export LESS_TERMCAP_us="[4m"  # underline
	*/

	os.Setenv("LESS_TERMCAP_mb", esc.Magenta)
	os.Setenv("LESS_TERMCAP_md", esc.Yellow)
	os.Setenv("LESS_TERMCAP_me", esc.Reset)
	os.Setenv("LESS_TERMCAP_se", esc.Reset)
	os.Setenv("LESS_TERMCAP_so", esc.Blue)
	os.Setenv("LESS_TERMCAP_ue", esc.Reset)
	os.Setenv("LESS_TERMCAP_us", esc.Under)

	term.EmphFromLess()

	fmt.Printf("%q\n", Z.Emphasize("*italic*"))
	fmt.Printf("%q\n", Z.Emphasize("**bold**"))
	fmt.Printf("%q\n", Z.Emphasize("**bolditalic**"))
	fmt.Printf("%q\n", Z.Emphasize("<under>"))

	// Output:
	// "\x1b[4mitalic\x1b[0m"
	// "\x1b[33mbold\x1b[0m"
	// "\x1b[33mbolditalic\x1b[0m"
	// "<\x1b[4munder\x1b[0m>"
}

func ExampleEmphasize_disable_with_Term() {
	term.AttrOff()
	fmt.Printf("%q\n", Z.Emphasize("*italic*"))
	fmt.Printf("%q\n", Z.Emphasize("**bold**"))
	fmt.Printf("%q\n", Z.Emphasize("**bolditalic**"))
	fmt.Printf("%q\n", Z.Emphasize("<under>"))
	// Output:
	// "italic"
	// "bold"
	// "bolditalic"
	// "<under>"
}

func ExampleFormat_remove_Initial_Blanks() {
	fmt.Printf("%q\n", Z.Format("\n   \n\n  \n   some"))
	// Output:
	// "some"
}

func ExampleFormat_wrapping() {
	fmt.Println(Z.Format(`
Here is a bunch of stuff just to fill the line beyond 80 columns so that it will wrap when it is supposed to and right now
as well if there was a hard return in the middle of a line.
`))
	// Output:
	// Here is a bunch of stuff just to fill the line beyond 80 columns so that it will
	// wrap when it is supposed to and right now
	// as well if there was a hard return in the middle of a line.
}
