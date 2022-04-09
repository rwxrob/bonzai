// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z_test

import (
	"fmt"
	"os"

	Z "github.com/rwxrob/bonzai/z"
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
		Commands: []*Z.Cmd{salut},
	}

	greet := &Z.Cmd{
		Name:     `greet`,
		Commands: []*Z.Cmd{yo, hi, hello, french},
	}

	cmd, args := greet.Seek(Z.ArgsFrom(`hi there`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french salut`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(Z.ArgsFrom(`french salut `))
	fmt.Printf("%v %q\n", cmd.Name, args)

	// Output:
	// hi ["there"]
	// salut []
	// salut [""]
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

func ExampleCmd_GetCommandNames() {
	foo := new(Z.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.GetCommandNames())
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

func ExampleCmd_PathString() {
	Z.ExitOff()

	z := new(Z.Cmd)
	c := z.Add("some")
	//fmt.Print(z.Commands[0].Name)
	c = c.Add("thing")
	//fmt.Print(z.Commands[0].Commands[0].Name)
	c = c.Add("deep")
	//fmt.Print(z.Commands[0].Commands[0].Commands[0].Name)

	c.Call = func(x *Z.Cmd, _ ...string) error {
		fmt.Println(x.PathString())
		return nil
	}

	defer func() { args := os.Args; os.Args = args }()
	os.Args = []string{"z", "some", "thing", "deep"} // first exe name

	z.Run()

	// Output:
	// .some.thing.deep
}

func ExampleCmd_Names() {

	x := &Z.Cmd{
		Name:    `foo`,
		Aliases: []string{"f", "FOO"},
	}
	fmt.Println(x.Names())

	//Output:
	// [f FOO foo]
}

func ExampleCmd_UsageNames() {

	x := &Z.Cmd{
		Name:    `foo`,
		Aliases: []string{"f", "FOO"},
	}
	fmt.Println(x.UsageNames())

	//Output:
	// (f|FOO|foo)
}

func ExampleCmd_UsageParams() {

	x := &Z.Cmd{
		Name:   `foo`,
		Params: []string{"p1", "p2"},
	}
	fmt.Println(x.UsageParams())

	x.MinParm = 1
	fmt.Println(x.UsageParams())

	x.MaxParm = 1
	fmt.Println(x.UsageParams())

	//Output:
	// (p1|p2)?
	// (p1|p2)+
	// (p1|p2)
}

func ExampleCmd_UsageError_commands_with_Aliases() {
	x := &Z.Cmd{
		Name: `cmd`,
		Commands: []*Z.Cmd{
			&Z.Cmd{Name: "foo", Aliases: []string{"f"}},
			&Z.Cmd{Name: "bar"},
		},
	}
	fmt.Println(x.UsageError())
	// Output:
	// usage: cmd ((f|foo)|bar)

}

func ExampleCmd_UsageError_params_but_No_Call() {
	x := &Z.Cmd{
		Name:   `cmd`,
		Params: []string{"p1", "p2"},
		Commands: []*Z.Cmd{
			&Z.Cmd{Name: "foo", Aliases: []string{"f"}},
			&Z.Cmd{Name: "bar"},
		},
	}
	fmt.Println(x.UsageError())
	// Output:
	// usage: cmd {ERROR: Params without Call: p1, p2}
}

func ExampleCmd_UsageError_no_Call_nor_Commands() {
	x := &Z.Cmd{
		Name: `cmd`,
	}
	fmt.Println(x.UsageError())
	// Output:
	// usage: cmd {ERROR: neither Call nor Commands defined}

}

func ExampleCmd_UsageCmdTitles() {
	x := &Z.Cmd{
		Name:   `cmd`,
		Params: []string{"p1", "p2"},
		Commands: []*Z.Cmd{
			&Z.Cmd{
				Name:    "foo",
				Aliases: []string{"f"},
				Summary: "foo the things",
			},
			&Z.Cmd{
				Name:    "bar",
				Summary: "bar the things",
			},
			&Z.Cmd{
				Name: "nosum",
			},
		},
	}
	fmt.Println(x.UsageCmdTitles())
	// Output:
	// f|foo - foo the things
	// bar   - bar the things
	// nosum
}
