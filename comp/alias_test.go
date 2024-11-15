package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

func ExampleAliases_Complete() {

	foo := &bonzai.Cmd{
		Name:  `foo`,
		Alias: `box|b`,
		Comp:  comp.Aliases,
	}

	foo.Comp.(bonzai.CmdCompleter).SetCmd(foo)

	// if no args, we have to assume the command isn't finished yet
	fmt.Println(foo.Comp.Complete())
	fmt.Println(foo.Comp.Complete(""))
	fmt.Println(foo.Comp.Complete("a"))
	fmt.Println(foo.Comp.Complete("b"))
	fmt.Println(foo.Comp.Complete("bo"))

	// Output:
	// []
	// [box b foo]
	// []
	// [box b]
	// [box]
}
