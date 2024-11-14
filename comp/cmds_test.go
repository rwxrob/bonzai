// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

func ExampleCmds_Complete() {

	foo := &bonzai.Cmd{
		Name: `foo`,
		Opts: `box`,
		Cmds: []*bonzai.Cmd{{Name: `bar`}, {Name: `blah`}},
		Comp: comp.Cmds,
	}

	foo.Comp.(bonzai.CmdCompleter).SetCmd(foo)

	// if no args, we have to assume the command isn't finished yet
	fmt.Println(foo.Comp.Complete())

	// we know it's not a command, but no prefix just yet
	// (usually this is when a space has been added after the command)
	fmt.Println(foo.Comp.Complete(""))

	// everything that begins with a (nothing)
	fmt.Println(foo.Comp.Complete(`a`))

	// everything that begins with b (which is everything)
	fmt.Println(foo.Comp.Complete(`b`))

	// everything that begins with bl (just blah)
	fmt.Println(foo.Comp.Complete(`bl`))

	// Output:
	// []
	// [bar blah]
	// []
	// [bar blah]
	// [blah]
}
