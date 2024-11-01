package redu_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/fn/redu"
)

func ExampleLongest() {

	set := []string{
		"some thing",
		"i'm the longest",
		"here",
	}
	fmt.Println(redu.Longest(set))

	// Output:
	// 15
}

func ExampleUnique() {

	set := []string{
		"some thing",
		"some thing",
		"i'm the longest",
		"here",
		"here",
		"here",
	}
	each.Println(redu.Unique(set))

	// Unordered output:
	// some thing
	// i'm the longest
	// here
}
