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
	out, _ := mark.Fill(help, funcs.Map, `The {{aka .}} command.`)
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
	out, _ := mark.Fill(nil, funcs.Map, long)
	fmt.Println(out)

	// Output:
	// `go mod init`
	//
	// I'll use the `{{code}}` thing instead with something like
	// `go mod init` since cannot use backticks.

}

func ExampleCommand() {
	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		Short: `foo this command`,
		Cmds:  []*bonzai.Cmd{subFooCmd},
	}

	var barCmd = &bonzai.Cmd{
		Name:  `bar`,
		Alias: `b`,
		Short: `bar this command`,
	}

	var Cmd = &bonzai.Cmd{
		Name:  `mycmd`,
		Alias: `my|cmd`,
		Short: `my command short summary`,
		Cmds:  []*bonzai.Cmd{fooCmd, barCmd},
		Def:   fooCmd,
	}

	Cmd.Seek(`foo`, `subfoo`) // required for default detection

	fmt.Print("# Commands\n\n")
	fmt.Println(funcs.Commands(Cmd))

	// Output:
	// # Commands
	//
	//     foo      ← foo this command (default)
	//       subfoo ← under the foo command
	//     bar      ← bar this command
}

func ExampleCommands_hidden() {
	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var hiddenCmd = &bonzai.Cmd{
		Name: `imhidden`,
		Cmds: []*bonzai.Cmd{{Name: `some`}, {Name: `other`}},
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		Short: `foo this command`,
		Cmds:  []*bonzai.Cmd{subFooCmd},
	}

	var barCmd = &bonzai.Cmd{
		Name:  `bar`,
		Alias: `b`,
		Short: `bar this command`,
	}

	var Cmd = &bonzai.Cmd{
		Name:  `mycmd`,
		Alias: `my|cmd`,
		Short: `my command short summary`,
		Cmds:  []*bonzai.Cmd{fooCmd, barCmd, hiddenCmd.AsHidden()},
		Def:   fooCmd,
	}

	Cmd.Seek(`foo`, `sssh`, `some`) // required for default detection

	fmt.Print("# Commands\n\n")
	fmt.Println(funcs.Commands(Cmd))

	// Output:
	// # Commands
	//
	//     foo      ← foo this command (default)
	//       subfoo ← under the foo command
	//     bar      ← bar this command

}
