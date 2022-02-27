package each_test

import (
	"fmt"
	"log"
	"os"

	"github.com/rwxrob/bonzai/each"
)

func ExampleDo() {
	set1 := []string{"doe", "ray", "mi"}
	each.Do(set1, func(s string) { fmt.Print(s) })
	fmt.Println()
	f1 := func() { fmt.Print("one") }
	f2 := func() { fmt.Print("two") }
	f3 := func() { fmt.Print("three") }
	set2 := []func(){f1, f2, f3}
	each.Do(set2, func(f func()) { f() })
	// Output:
	// doeraymi
	// onetwothree
}

func ExamplePrintln() {
	set := []string{"doe", "ray", "mi"}
	each.Println(set)
	bools := []bool{false, true, true}
	each.Println(bools)
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
	each.Print(set)
	bools := []bool{false, true, true}
	each.Print(bools)
	// Output:
	// doeraymifalsetruetrue
}

func ExamplePrintf() {
	set := []string{"doe", "ray", "mi"}
	each.Printf(set, "sing %v\n")
	// Output:
	// sing doe
	// sing ray
	// sing mi
}

func ExampleLogf() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	set := []string{"doe", "ray", "mi"}
	each.Logf(set, "sing %v\n")
	// Output:
	// sing doe
	// sing ray
	// sing mi
}

func ExampleLog() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	set := []string{"doe", "ray", "mi"}
	each.Log(set)
	// Output:
	// doe
	// ray
	// mi
}
