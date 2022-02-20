// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/filter"
)

func ExampleStandard() {
	foo := new(bonzai.Cmd)
	foo.Params = []string{"box"}
	foo.Add("bar")
	foo.Add("blah")

	// all commands and params since nothing specified
	args := []string{}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with a (nothing)
	args = []string{`a`}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with b (which is everything)
	args = []string{`b`}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with bl (just blah)
	args = []string{`bl`}
	fmt.Println(comp.Standard(foo, args))

	// give own completer for days of the week
	foo.Completer = func(cmd comp.Command, args []string) []string {
		list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
		if len(args) == 0 {
			return list
		}
		return filter.HasPrefix(list, args[0])
	}
	args = []string{`t`}
	fmt.Println(comp.Standard(foo, args))

	//Output:
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
	// [tue thu]

}
