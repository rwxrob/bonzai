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

package filter_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/filter"
)

func ExampleHasPrefix() {
	set := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(filter.HasPrefix(set, "t"))
	// Output:
	// [two three]
}

func ExampleMinus() {
	set := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(filter.Minus(set, []string{"two", "four", "six"}))
	// Output:
	// [one three five seven]
}

func ExamplePrintln() {
	set := []string{"doe", "ray", "mi"}
	filter.Println(set)
	bools := []bool{false, true, true}
	filter.Println(bools)
	// Output:
	// doe
	// ray
	// mi
	// false
	// true
	// true
}
