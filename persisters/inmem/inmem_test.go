package inmem_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/persisters/inmem"
)

func ExampleSetup() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %v\n", r)
		}
	}()

	p := new(inmem.Persister)
	p.Setup()                  // initialized internal map
	p.Set(`some`, `thing`)     // no panic
	fmt.Println(p.Get(`some`)) // "thing"

	// Output:
	// thing
}

func ExampleSetup_panicOnSetWithoutSetup() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %v\n", r)
		}
	}()

	p := new(inmem.Persister) // simple construction
	p.Get(`some`)             // doesn't panic because only nil map inserts panic
	p.Set(`some`, `thing`)    // panics
	p.Setup()                 // initialized internal map

	// Output:
	// panic: assignment to entry in nil map
}
