package check_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/check"
)

func ExampleBlank() {
	fmt.Println(check.Blank(""))
	fmt.Println(check.Blank(nil))
	fmt.Println(check.Blank([]string{""}))
	fmt.Println(check.Blank([][]byte{}))
	// and now for true
	fmt.Println(check.Blank("some"))
	fmt.Println(check.Blank([]string{"some"}))
	fmt.Println(check.Blank([][]byte{{'a'}}))

	// Output:
	// false
	// false
	// false
	// false
	// true
	// true
	// true
}
