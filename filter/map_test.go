package filter_test

import "github.com/rwxrob/bonzai/filter"

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
