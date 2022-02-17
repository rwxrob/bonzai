/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bonzai_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
)

func ExampleCmd_Seek() {

	hello := &bonzai.Cmd{
		Name:   `hello`,
		Params: []string{"there"},
		Method: func(args ...string) error {
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
		Method: func(args ...string) error {
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
		Method: func(args ...string) error {
			fmt.Println("yo")
			return nil
		},
	}

	salut := &bonzai.Cmd{
		Name:   `salut`,
		Params: []string{"la"},
		Method: func(args ...string) error {
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
	// salut [" "]
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
	foo.Params = []string{"box"}
	foo.Add("bar")
	foo.Add("blah")
	foo.Add("other")
	fmt.Println(foo.GetCommands())
	fmt.Println(foo.GetCommands("b"))
	// Output:
	// [bar blah other]
	// [bar blah]
}

func ExampleCmd_GetParams() {
	foo := new(bonzai.Cmd)
	foo.Params = []string{"box", "bing", "and"}
	foo.Add("bar")
	foo.Add("blah")
	fmt.Println(foo.GetParams())
	fmt.Println(foo.GetParams("b"))
	// Output:
	// [box bing and]
	// [box bing]
}
