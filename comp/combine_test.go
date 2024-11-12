// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleCmdsOpts_Complete() {
	// dynamic/runtime Cmd creation (normally a niche use case)
	foo := new(bonzai.Cmd)
	foo.Opts = `box`
	foo.Comp = comp.CmdsOpts
	foo.Add(`bar`)
	foo.Add(`blah`)
	comp.CmdsOpts.SetCmd(foo) // [bonzai.CmpCompleter]

	// we know it's not a command, but no prefix just yet
	// (usually this is when a space has been added after the command)
	fmt.Println(comp.CmdsOpts.Complete(""))

	// everything that begins with a (nothing)
	fmt.Println(comp.CmdsOpts.Complete(`a`))

	// everything that begins with b (which is everything)
	fmt.Println(comp.CmdsOpts.Complete(`b`))

	// everything that begins with bl (just blah)
	fmt.Println(comp.CmdsOpts.Complete(`bl`))

	// Output:
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
}
