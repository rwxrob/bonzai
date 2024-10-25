// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
)

func ExampleCmds_Complete() {
	foo := new(bonzai.Cmd)
	foo.Params = `box`
	foo.Add(`bar`)
	foo.Add(`blah`)
	/*
		// if no args, we have to assume the command isn't finished yet
		fmt.Println(comp.Cmds.Complete(foo))

		// we know it's not a command, but no prefix just yet
		// (usually this is when a space has been added after the command)
		fmt.Println(comp.Cmds.Complete(foo, ""))

		// everything that begins with a (nothing)
		fmt.Println(comp.Cmds.Complete(foo, `a`))

		// everything that begins with b (which is everything)
		fmt.Println(comp.Cmds.Complete(foo, `b`))

		// everything that begins with bl (just blah)
		fmt.Println(comp.Cmds.Complete(foo, `bl`))
	*/
	//Output:
	// []
	// [bar blah]
	// []
	// [bar blah]
	// [blah]

}
