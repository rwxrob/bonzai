// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z_test

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/inc/help"
)

func ExampleCmd_Seek() {

	hello := &Z.Cmd{
		Name:   `hello`,
		Params: []string{"there"},
		Call: func(_ *Z.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("hello %v\n", args[0])
				return nil
			}
			fmt.Println("hello")
			return nil
		},
	}

	hi := &Z.Cmd{
		Name:   `hi`,
		Params: []string{"there", "ya"},
		Call: func(_ *Z.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("hi %v\n", args[0])
				return nil
			}
			fmt.Println("hi")
			return nil
		},
	}

	yo := &Z.Cmd{
		Name: `yo`,
		Call: func(x *Z.Cmd, args ...string) error {
			fmt.Println("yo")
			return nil
		},
	}

	salut := &Z.Cmd{
		Name:   `salut`,
		Params: []string{"la"},
		Call: func(_ *Z.Cmd, args ...string) error {
			if len(args) > 0 {
				fmt.Printf("salut %v\n", args[0])
				return nil
			}
			fmt.Println("salut")
			return nil
		},
	}

	french := &Z.Cmd{
		Name:     `french`,
		Aliases:  []string{"fr"},
		Commands: []*Z.Cmd{help.Cmd, salut},
	}

	greet := &Z.Cmd{
		Name:     `greet`,
		Commands: []*Z.Cmd{help.Cmd, yo, hi, hello, french},
	}

	cmd, args := greet.Seek(Z.ArgsFrom(`hi there`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french salut`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french salut `))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french h`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french help`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	// Output:
	// hi ["there"]
	// salut []
	// salut [""]
	// french ["h"]
	// help []
}

func ExampleCmd_CmdNames() {
	foo := new(Z.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.CmdNames())
	// Output:
	// [bar blah other]
}

func ExampleCmd_GetCommands() {
	foo := new(Z.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.GetCommands())
	// Output:
	// [bar blah other]
}

func ExampleCmd_GetParams() {
	foo := new(Z.Cmd)
	foo.Params = []string{"box", "bing", "and"}
	fmt.Println(foo.GetParams())
	// Output:
	// [box bing and]
}

func ExampleCmd_Branch() {
	Z.ExitOff()

	z := new(Z.Cmd)
	c := z.Add("some")
	//fmt.Print(z.Commands[0].Name)
	c = c.Add("thing")
	//fmt.Print(z.Commands[0].Commands[0].Name)
	c = c.Add("deep")
	//fmt.Print(z.Commands[0].Commands[0].Commands[0].Name)

	c.Call = func(x *Z.Cmd, _ ...string) error {
		fmt.Println(x.Branch())
		return nil
	}

	defer func() { args := os.Args; os.Args = args }()
	os.Args = []string{"z", "some", "thing", "deep"}

	z.Run()

	// Output:
	// some.thing.deep
}
