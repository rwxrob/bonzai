package loop_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/loop"
)

func ExampleAll() {
	set1 := []string{"doe", "ray", "mi"}
	loop.All(set1, func(s string) { fmt.Print(s) })
	fmt.Println()
	f1 := func() { fmt.Print("one") }
	f2 := func() { fmt.Print("two") }
	f3 := func() { fmt.Print("three") }
	set2 := []func(){f1, f2, f3}
	loop.All(set2, func(f func()) { f() })
	// Output:
	// doeraymi
	// onetwothree
}

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
