package loop_test

import "github.com/rwxrob/bonzai/loop"

func ExamplePrintln() {
	set := []string{"doe", "ray", "mi"}
	loop.Println(set)
	bools := []bool{false, true, true}
	loop.Println(bools)
	// Output:
	// doe
	// ray
	// mi
	// false
	// true
	// true
}
