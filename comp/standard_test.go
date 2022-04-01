// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	Z "github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/fn/filt"
)

func ExampleStandard() {
	foo := new(Z.Cmd)
	foo.Params = []string{"box"}
	foo.Add("bar")
	foo.Add("blah")

	// if no args, we have to assume the command isn't finished yet
	fmt.Println(comp.Standard(foo))

	// we know it's not a command, but no prefix just yet
	// (usually this is when a space has been added after the command)
	fmt.Println(comp.Standard(foo, ""))

	// everything that begins with a (nothing)
	fmt.Println(comp.Standard(foo, `a`))

	// everything that begins with b (which is everything)
	fmt.Println(comp.Standard(foo, `b`))

	// everything that begins with bl (just blah)
	fmt.Println(comp.Standard(foo, `bl`))

	// give own completer for days of the week
	foo.Completer = func(cmd comp.Command, args ...string) []string {
		list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
		if len(args) == 0 {
			return list
		}
		return filt.HasPrefix(list, args[0])
	}
	fmt.Println(comp.Standard(foo, `t`))

	//Output:
	// []
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
	// [tue thu]

}
