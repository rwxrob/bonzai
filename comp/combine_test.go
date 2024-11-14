// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/tr"
)

func ExampleCmdsOpts_Complete() {

	foo := &bonzai.Cmd{
		Name: `foo`,
		Opts: `box`,
		Comp: comp.CmdsOpts,
		Cmds: []*bonzai.Cmd{{Name: `bar`}, {Name: `blah`}},
	}

	foo.Comp.(bonzai.CmdCompleter).SetCmd(foo)

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

func ExampleCombine_transform() {

	foo := &bonzai.Cmd{
		Name: `foo`,
		Opts: `box`,
		Comp: comp.Combine{comp.CmdsOpts, tr.Prefix{`aprefix`}},
		Cmds: []*bonzai.Cmd{{Name: `bar`}, {Name: `blah`}},
	}

	foo.Comp.(bonzai.CmdCompleter).SetCmd(foo)

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
	// [aprefixbar aprefixblah aprefixbox]
	// []
	// [aprefixbar aprefixblah aprefixbox]
	// [aprefixblah]

}
