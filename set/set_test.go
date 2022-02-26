package set_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/set"
)

func ExampleMinus() {
	s := []string{
		"one", "two", "three", "four", "five", "six", "seven",
	}
	fmt.Println(set.Minus(s, []string{"two", "four", "six"}))
	// Output:
	// [one three five seven]
}
