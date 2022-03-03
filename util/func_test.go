package util_test

import (
	"github.com/rwxrob/bonzai/each"
	"github.com/rwxrob/bonzai/fn"
	"github.com/rwxrob/bonzai/util"
)

func Foo() {}

func ExampleFuncName() {

	f1 := func() {}
	f2 := func() {}

	// Foo is defined outside of the ExampleFuncName

	each.Println(fn.Map([]any{f1, f2, Foo, util.Files}, util.FuncName))

	// Output:
	// func1
	// func2
	// Foo
	// Files
}
