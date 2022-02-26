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

func ExamplePrint() {
	set := []string{"doe", "ray", "mi"}
	loop.Print(set)
	bools := []bool{false, true, true}
	loop.Print(bools)
	// Output:
	// doeraymifalsetruetrue
}

func ExamplePrintf() {
	set := []string{"doe", "ray", "mi"}
	loop.Printf(set, "sing %v\n")
	// Output:
	// sing doe
	// sing ray
	// sing mi
}
