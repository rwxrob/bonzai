// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package vars_test

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/json"

	"github.com/rwxrob/bonzai/vars"
)

func ExampleComp_Complete() {
	foo := bonzai.Cmd{}
	foo.Comp = vars.Comp

	defer vars.Data.Delete(`some`)
	defer vars.Data.Delete(`someother`)
	vars.Data.Set(`some`, `thing`)
	vars.Data.Set(`someother`, `awesome`)

	json.This{This: vars.Comp.Complete(foo, `s`)}.Print()
	json.This{This: vars.Comp.Complete(foo, `some`)}.Print()

	// Output:
	// ["some","someother"]
	// ["some","someother"]
}
