package check_test

import (
	"fmt"
	"reflect"

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

func ExampleIsNil() {
	var names []string
	var namesi interface{}
	namesi = names
	fmt.Println(names == nil)
	fmt.Println(namesi == nil)
	fmt.Println(reflect.ValueOf(namesi).IsNil())
	fmt.Println(check.IsNil(namesi))
	// Output:
	// true
	// false
	// true
	// true
}
