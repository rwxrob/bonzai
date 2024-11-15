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

func ExampleCode() {

	long := `
I'll use the {{code "{{code}}"}} thing instead with something like
{{code "go mod init"}} since cannot use backticks.`

	fmt.Println(funcs.Code(`go mod init`))
	out, _ := mark.Render(nil, funcs.Map, long)
	fmt.Println(out)

	// Output:
	// `go mod init`
	//
	// I'll use the `{{code}}` thing instead with something like
	// `go mod init` since cannot use backticks.

}
