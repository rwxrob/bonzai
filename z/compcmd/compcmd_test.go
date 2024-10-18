// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package compcmd_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/bonzai/z/compcmd"
	"github.com/rwxrob/bonzai/z/fn/filt"
)

// give own completer for days of the week
type weekcomp struct{}

func NewWeekComp() *weekcomp { return new(weekcomp) }

func (weekcomp) Complete(cmd bonzai.Command, args ...string) []string {
	list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	if len(args) == 0 {
		return list
	}
	return filt.HasPrefix(list, args[0])
}

func ExampleStandard() {
	foo := new(Z.Cmd)
	foo.Params = []string{"box"}
	foo.Add("bar")
	foo.Add("blah")

	// if no args, we have to assume the command isn't finished yet
	fmt.Println(compcmd.New().Complete(foo))

	// we know it's not a command, but no prefix just yet
	// (usually this is when a space has been added after the command)
	fmt.Println(compcmd.New().Complete(foo, ""))

	// everything that begins with a (nothing)
	fmt.Println(compcmd.New().Complete(foo, `a`))

	// everything that begins with b (which is everything)
	fmt.Println(compcmd.New().Complete(foo, `b`))

	// everything that begins with bl (just blah)
	fmt.Println(compcmd.New().Complete(foo, `bl`))

	/* (note this has to happen outside of block because of receiver)
	// give own completer for days of the week
	type weekcomp struct{}

	func NewWeekComp() *weekcomp{ return new(weekcomp) }

	func (weekcomp) Complete(cmd bonzai.Command, args ...string) []string {
		list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
		if len(args) == 0 {
			return list
		}
		return filt.HasPrefix(list, args[0])
	}
	*/

	fmt.Println(NewWeekComp().Complete(foo, `t`))

	//Output:
	// []
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
	// [tue thu]

}
