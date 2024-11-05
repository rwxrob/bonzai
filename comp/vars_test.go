// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/json"
)

func ExampleVars_Complete() {
	foo := new(bonzai.Cmd)
	foo.Comp = comp.Vars

	defer bonzai.Vars.Delete(`some`)
	defer bonzai.Vars.Delete(`someother`)
	bonzai.Vars.Set(`some`, `thing`)
	bonzai.Vars.Set(`someother`, `awesome`)

	json.This{comp.Vars.Complete(foo, `s`)}.Print()
	json.This{comp.Vars.Complete(foo, `some`)}.Print()

	//Output:
	// ["some","someother"]
	// ["some","someother"]

}
