package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleAliases_Complete() {
	foo := bonzai.Cmd{}
	foo.Alias = `box|b`

	// if no args, we have to assume the command isn't finished yet
	fmt.Println(comp.Aliases.Complete(foo))
	fmt.Println(comp.Aliases.Complete(foo, ""))
	fmt.Println(comp.Aliases.Complete(foo, "a"))
	fmt.Println(comp.Aliases.Complete(foo, "b"))
	fmt.Println(comp.Aliases.Complete(foo, "bo"))

	// Output:
	// []
	// [box b]
	// []
	// [box b]
	// [box]
}
