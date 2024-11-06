// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp_test

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/json"
	"github.com/rwxrob/bonzai/vars"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleVars_Complete() {
	foo := bonzai.Cmd{}
	foo.Comp = comp.Vars

	defer vars.Data.Delete(`some`)
	defer vars.Data.Delete(`someother`)
	vars.Data.Set(`some`, `thing`)
	vars.Data.Set(`someother`, `awesome`)

	json.This{comp.Vars.Complete(foo, `s`)}.Print()
	json.This{comp.Vars.Complete(foo, `some`)}.Print()

	// Output:
	// ["some","someother"]
	// ["some","someother"]
}
