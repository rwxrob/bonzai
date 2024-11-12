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
	fmt.Println(comp.Aliases.Complete())
	fmt.Println(comp.Aliases.Complete(""))
	fmt.Println(comp.Aliases.Complete("a"))
	fmt.Println(comp.Aliases.Complete("b"))
	fmt.Println(comp.Aliases.Complete("bo"))

	// Output:
	// []
	// [box b foo]
	// []
	// [box b]
	// [box]
}
