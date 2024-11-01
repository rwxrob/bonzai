// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/comp"
)

func ExampleOpts_Complete() {
	foo := new(bonzai.Cmd)
	foo.Opts = `box`
	foo.Comp = comp.Opts
	foo.Add(`bar`)
	foo.Add(`blah`)

	// everything that begins with b (which is everything)
	fmt.Println(comp.Opts.Complete(foo, `b`))

	// everything that begins with bl (just blah)
	fmt.Println(comp.Opts.Complete(foo, `bl`))

	//Output:
	// [box]
	// []

}
