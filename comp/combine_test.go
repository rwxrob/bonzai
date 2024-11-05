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
	foo := bonzai.Cmd{}
	foo.Opts = `box`
	foo.Comp = comp.CmdsOpts
	foo.Add(`bar`)
	foo.Add(`blah`)

	// we know it's not a command, but no prefix just yet
	// (usually this is when a space has been added after the command)
	fmt.Println(comp.CmdsOpts.Complete(foo, ""))

	// everything that begins with a (nothing)
	fmt.Println(comp.CmdsOpts.Complete(foo, `a`))

	// everything that begins with b (which is everything)
	fmt.Println(comp.CmdsOpts.Complete(foo, `b`))

	// everything that begins with bl (just blah)
	fmt.Println(comp.CmdsOpts.Complete(foo, `bl`))

	// Output:
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
}
