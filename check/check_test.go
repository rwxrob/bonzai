package check_test

import (
	"fmt"
	"reflect"

	"github.com/rwxrob/bonzai/check"
)

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

func ExampleIsDir() {
	fmt.Println(check.IsDir("testdata"))
	fmt.Println(check.IsDir("nothing"))
	// Output:
	// true
	// false
}
