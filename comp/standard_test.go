/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/filter"
)

func ExampleStandard() {
	foo := new(bonzai.Cmd)
	foo.Params = []string{"box"}
	foo.Add("bar")
	foo.Add("blah")

	// all commands and params since nothing specified
	args := []string{}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with a (nothing)
	args = []string{`a`}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with b (which is everything)
	args = []string{`b`}
	fmt.Println(comp.Standard(foo, args))

	// everything that begins with bl (just blah)
	args = []string{`bl`}
	fmt.Println(comp.Standard(foo, args))

	// give own completer for days of the week
	foo.Completer = func(cmd comp.Command, args []string) []string {
		list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
		if len(args) == 0 {
			return list
		}
		return filter.HasPrefix(list, args[0])
	}
	args = []string{`t`}
	fmt.Println(comp.Standard(foo, args))

	//Output:
	// [bar blah box]
	// []
	// [bar blah box]
	// [blah]
	// [tue thu]

}
