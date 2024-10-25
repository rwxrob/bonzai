// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai_test

import (
	"fmt"
	"os"

	bonzai "github.com/rwxrob/bonzai/pkg"
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
		Call: func(x *bonzai.Cmd, args ...string) error {
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
		Alias:  []string{"fr"},
		Commands: []*bonzai.Cmd{salut},
	}

	greet := &bonzai.Cmd{
		Name:     `greet`,
		Commands: []*bonzai.Cmd{yo, hi, hello, french},
	}

	cmd, args := greet.Seek(bonzai.ArgsFrom(`hi there`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french salut`))
	fmt.Printf("%v %q\n", cmd.Name, args)

	cmd, args = greet.Seek(bonzai.ArgsFrom(`french salut `))
	fmt.Printf("%v %q\n", cmd.Name, args)

	// Output:
	// hi ["there"]
	// salut []
	// salut [""]
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

func ExampleCmd_GetCommandNames() {
	foo := new(bonzai.Cmd)
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.GetCommandNames())
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

func ExampleCmd_Path() {
	bonzai.ExitOff()

	z := new(bonzai.Cmd)
	c := z.Add("some")
	//fmt.Print(z.Commands[0].Name)
	c = c.Add("thing")
	//fmt.Print(z.Commands[0].Commands[0].Name)
	c = c.Add("deep")
	//fmt.Print(z.Commands[0].Commands[0].Commands[0].Name)

	c.Call = func(x *bonzai.Cmd, _ ...string) error {
		fmt.Println(x.Path())
		fmt.Println(x.Path(`and`, `some`, `more`))
		return nil
	}

	defer func() { args := os.Args; os.Args = args }()
	os.Args = []string{"z", "some", "thing", "deep"} // first exe name

	z.Run()

	// Output:
	// .some.thing.deep
	// .some.thing.deep.and.some.more
}

func ExampleCmd_Names() {

	x := &bonzai.Cmd{
		Name:    `foo`,
		Alias: []string{"-f", "@f", "f", "FOO"},
	}
	fmt.Println(x.Names())

	//Output:
	// [f FOO foo]
}

func ExampleCmd_UsageNames() {

	x := &bonzai.Cmd{
		Name:    `foo`,
		Alias: []string{"f", "FOO"},
	}
	fmt.Println(x.UsageNames())

	//Output:
	// (f|FOO|foo)
}

func ExampleCmd_UsageParams() {

	x := &bonzai.Cmd{
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

func ExampleCmd_UsageError_commands_with_Alias() {
	x := &bonzai.Cmd{
		Name: `cmd`,
		Commands: []*bonzai.Cmd{
			&bonzai.Cmd{Name: "foo", Alias: []string{"f"}},
			&bonzai.Cmd{Name: "bar"},
		},
	}
	fmt.Println(bonzai.IncorrectUsage{x})
	// Output:
	// usage: cmd ((f|foo)|bar)

}

func ExampleCmd_UsageError_params_but_No_Call() {
	x := &bonzai.Cmd{
		Name:   `cmd`,
		Params: []string{"p1", "p2"},
		Commands: []*bonzai.Cmd{
			&bonzai.Cmd{Name: "foo", Alias: []string{"f"}},
			&bonzai.Cmd{Name: "bar"},
		},
	}
	fmt.Println(bonzai.IncorrectUsage{x})
	// Output:
	// usage: cmd {ERROR: Params without Call: p1, p2}
}

func ExampleCmd_UsageError_no_Call_nor_Commands() {
	x := &bonzai.Cmd{
		Name: `cmd`,
	}
	fmt.Println(bonzai.IncorrectUsage{x})
	// Output:
	// usage: cmd {ERROR: neither Call nor Commands defined}

}

func ExampleCmd_UsageCmdTitles() {
	x := &bonzai.Cmd{
		Name:   `cmd`,
		Params: []string{"p1", "p2"},
		Hide: []string{"hidden"},
		Commands: []*bonzai.Cmd{
			&bonzai.Cmd{
				Name:    "foo",
				Alias: []string{"f"},
				Summary: "foo the things",
			},
			&bonzai.Cmd{
				Name:    "bar",
				Summary: "bar the things",
			},
			&bonzai.Cmd{
				Name: "nosum",
			},
			&bonzai.Cmd{
				Name:    "hidden",
				Summary: "not listed, but works",
			},
		},
	}
	fmt.Println(x.UsageCmdTitles())
	// Output:
	// f|foo - foo the things
	// bar   - bar the things
	// nosum
}

func ExampleCmd_UsageCmdShortcuts() {
	x := &bonzai.Cmd{
		Name: `cmd`,
		Shortcuts: bonzai.ArgMap{
			"foo": {"a", "long", "way", "to", "foo"},
			"bar": {"a", "long", "long", "way", "to", "bar"},
		},
	}
	fmt.Println(x.UsageCmdShortcuts())
	// Unordered Output:
	// foo - a long way to foo
	// bar - a long long way to bar
}
