package filt_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/comp/filt"
)

func ExamplePrefix_Filter() {
	p := filt.Prefix("b")
	fmt.Println(p.Filter([]string{"foo", "bar", "baz"}))

	// Output:
	// [bar baz]
}

func ExampleMaxLen_Filter() {
	m := filt.MaxLen(3)
	fmt.Println(m.Filter([]string{"foo", "bar", "blah"}))

	// Output:
	// [foo bar]
}

func ExampleMinLen_Filter() {
	m := filt.MinLen(4)
	fmt.Println(m.Filter([]string{"foo", "bar", "blah"}))

	// Output:
	// [blah]
}
