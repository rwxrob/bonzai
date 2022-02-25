// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmd"
)

func ExampleCmd_Seek() {

	hello := &bonzai.Cmd{
		Name:   `hello`,
		Params: []string{"there"},
		Call: func(_ *bonzai.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("hello %v\n", args[0])
				return nil
			}
			fmt.Println("hello")
			return nil
		},
	}

	hi := &bonzai.Cmd{
		Name:   `hi`,
		Params: []string{"there", "ya"},
		Call: func(_ *bonzai.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("hi %v\n", args[0])
				return nil
			}
			fmt.Println("hi")
			return nil
		},
	}

	yo := &bonzai.Cmd{
		Name: `yo`,
		Call: func(_ *bonzai.Cmd, args ...string) error {
			fmt.Println("yo")
			return nil
		},
	}

	salut := &bonzai.Cmd{
		Name:   `salut`,
		Params: []string{"la"},
		Call: func(_ *bonzai.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("salut %v\n", args[0])
				return nil
			}
			fmt.Println("salut")
			return nil
		},
	}

	french := &bonzai.Cmd{
		Name:     `french`,
		Aliases:  []string{"fr"},
		Commands: []*bonzai.Cmd{cmd.Help, salut},
	}

	greet := &bonzai.Cmd{
		Name:     `greet`,
		Commands: []*bonzai.Cmd{cmd.Help, yo, hi, hello, french},
	}

	cmd, args := greet.Seek(bonzai.ArgsFrom(`hi there`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french salut`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french salut `))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french h`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french help`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	// Output:
	// hi ["there"]
	// salut []
	// salut [""]
	// french ["h"]
	// help []
}

func ExampleCmd_CmdNames() {
	foo := new(bonzai.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.CmdNames())
	// Output:
	// [bar blah other]
}

func ExampleCmd_GetCommands() {
	foo := new(bonzai.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.GetCommands())
	// Output:
	// [bar blah other]
}

func ExampleCmd_GetParams() {
	foo := new(bonzai.Cmd)
	foo.Params = []string{"box", "bing", "and"}
	fmt.Println(foo.GetParams())
	// Output:
	// [box bing and]
}
