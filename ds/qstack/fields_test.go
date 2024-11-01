package qstack_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/ds/qstack"
)

func ExampleFields() {

	text := `
      Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
      tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
      veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
      commodo consequat.

			Duis aute irure dolor in reprehenderit in voluptate velit esse
			cillum dolore eu fugiat nulla pariatur. Excepteur sint
      occaecat cupidatat non proident, sunt in culpa qui officia deserunt
      mollit anim id est laborum.
	`

	fields := qstack.Fields(text)

	fmt.Println(fields.Len)
	fmt.Println(fields.Items()[0])
	fmt.Println(fields.Items()[34])
	fmt.Println(fields.Items()[68])

	// Output:
	// 69
	// Lorem
	// commodo
	// laborum.
}

func ExampleFields_scan() {
	fields := qstack.Fields("some thing")
	fmt.Println(fields)
	for fields.Scan() {
		fmt.Print(fields.Current())
	}
	// Output:
	// ["some","thing"]
	// something
}
