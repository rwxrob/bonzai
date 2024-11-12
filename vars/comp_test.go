// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package vars_test

import (
	"github.com/rwxrob/bonzai/json"

	"github.com/rwxrob/bonzai/vars"
)

func ExampleComp_Complete() {
	defer vars.Data.Delete(`some`)
	defer vars.Data.Delete(`someother`)
	vars.Data.Set(`some`, `thing`)
	vars.Data.Set(`someother`, `awesome`)

	json.This{This: vars.Comp.Complete(`s`)}.Print()
	json.This{This: vars.Comp.Complete(`some`)}.Print()

	// Output:
	// ["some","someother"]
	// ["some","someother"]
}
