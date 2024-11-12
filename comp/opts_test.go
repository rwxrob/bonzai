// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

func ExampleOpts_Complete() {

	foo := &bonzai.Cmd{
		Name: `foo`,
		Opts: `box`,
		Comp: comp.Opts,
		Cmds: []*bonzai.Cmd{{Name: `bar`}, {Name: `blah`}},
	}

	foo.Comp.(bonzai.CmdCompleter).SetCmd(foo)

	// everything that begins with b (which is everything)
	fmt.Println(comp.Opts.Complete(`b`))

	// everything that begins with bl (just blah)
	fmt.Println(comp.Opts.Complete(`bl`))

	// Output:
	// [box]
	// []
}
