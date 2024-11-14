package funcs_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/mark/funcs"
)

func ExampleAKA() {

	var help = &bonzai.Cmd{
		Name:  `help`,
		Alias: `h|-h|-help|--help|/?`,
	}

	fmt.Println(funcs.AKA(help))
	out, _ := mark.Render(help, funcs.Map, `The {{aka .}} command.`)
	fmt.Println(out)

	// Output:
	// `help` (`h`|`-h`|`-help`|`--help`|`/?`)
	// The `help` (`h`|`-h`|`-help`|`--help`|`/?`) command.
}
