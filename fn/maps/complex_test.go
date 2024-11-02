// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package maps_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/fn/maps"
)

func ExamplePrefix() {
	fmt.Println(maps.Prefix([]string{"foo", "bar"}, "my"))
	// Output:
	// [myfoo mybar]
}

func ExampleKeys() {
	m1 := map[string]int{"two": 2, "three": 3, "one": 1}
	m2 := map[string]string{"two": "two", "three": "three", "one": "one"}
	fmt.Println(maps.Keys(m1))
	fmt.Println(maps.Keys(m2))
	// Output:
	// [one three two]
	// [one three two]
}

func ExampleKeysWithPrefix() {
	m1 := map[string]int{"two": 2, "three": 3, "one": 1}
	fmt.Println(maps.KeysWithPrefix(m1, "t"))
	// Output:
	// [three two]
}
