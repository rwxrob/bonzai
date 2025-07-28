package mark_test

import (
	"fmt"
	"text/template"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
)

type Thing struct {
	Name  string
	Count int
}

func (a Thing) Summary() string {
	return fmt.Sprintf("%v %v", a.Name, a.Count)
}

func ExampleFill() {

	/* cannot declare type with method within function, but this is it

	type Thing struct {
		Name  string
		Count int
	}

	func (a Thing) Summary() string {
		return fmt.Sprintf("%v %v", a.Name, a.Count)
	}

	*/

	thing := Thing{`Thing`, 20}
	tmpl := `
	{{hello}}, my name is {{.Name}} with {{.Count}}. Summary: {{.Summary}}`
	funcs := template.FuncMap{}
	funcs[`hello`] = func() string { return `Hello` }

	out, err := mark.Fill(thing, funcs, tmpl)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out)

	// Output:
	// Hello, my name is Thing with 20. Summary: Thing 20

}

func ExampleUsage_withHiddenCmds() {
	var subFooHiddenCmd = &bonzai.Cmd{
		Name:  `iamhidden`,
		Short: `i am hidden`,
	}

	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		Short: `foo this command`,
		Cmds:  []*bonzai.Cmd{subFooCmd, subFooHiddenCmd.AsHidden()},
		// Cmds:  []*bonzai.Cmd{subFooCmd, subFooHiddenCmd},
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
		Long: `
			Here is a long description.
			On multiple lines.`,
	}

	Cmd.Run()
	out, _ := mark.Bonzai(Cmd)
	fmt.Println(string(out))

	// Output:
	// # NAME
	//
	// `mycmd` - my command short summary
	//
	// # USAGE
	//
	//     my|cmd|mycmd COMMAND
	//
	// # COMMANDS
	//
	//     foo      - foo this command
	//       subfoo - under the foo command
	//     bar      - bar this command
	//
	// # DESCRIPTION
	//
	// Here is a long description.
	// On multiple lines.

}

func ExampleUsage_missingShort() {
	var subFooHiddenCmd = &bonzai.Cmd{
		Name:  `iamhidden`,
		Short: `i am hidden`,
	}

	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		//Short: `foo this command`,
		Cmds: []*bonzai.Cmd{subFooCmd, subFooHiddenCmd.AsHidden()},
		// Cmds:  []*bonzai.Cmd{subFooCmd, subFooHiddenCmd},
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
		Long: `
			Here is a long description.
			On multiple lines.`,
	}

	Cmd.Run()
	out, _ := mark.Bonzai(Cmd)
	fmt.Println(out)

	// Output:
	// # NAME
	//
	// `mycmd` - my command short summary
	//
	// # USAGE
	//
	//     my|cmd|mycmd COMMAND
	//
	// # COMMANDS
	//
	//     foo
	//       subfoo - under the foo command
	//     bar      - bar this command
	//
	// # DESCRIPTION
	//
	// Here is a long description.
	// On multiple lines.

}

func ExampleUsage_middle() {
	var subFooHiddenCmd = &bonzai.Cmd{
		Name:  `iamhidden`,
		Short: `i am hidden`,
	}

	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		//Short: `foo this command`,
		Cmds: []*bonzai.Cmd{subFooCmd, subFooHiddenCmd.AsHidden()},
		// Cmds:  []*bonzai.Cmd{subFooCmd, subFooHiddenCmd},
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
		Long: `
			Here is a long description.
			On multiple lines.`,
	}

	_ = Cmd

	Cmd.Run()
	out, err := mark.Bonzai(fooCmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out)

	// Output:
	// # NAME
	//
	// `foo`
	//
	// # USAGE
	//
	//     f|foo COMMAND
	//
	// # COMMANDS
	//
	//     subfoo - under the foo command

}

func ExampleUsage_longFirstName() {

	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		//Short: `a foo`,
		Do: func(*bonzai.Cmd, ...string) error {
			return nil
		},
	}

	var Cmd = &bonzai.Cmd{
		Name:  `help-test`,
		Alias: `h|ht`,
		Short: `just a help test`,
		Opts:  `some|-y|--yaml`,
		Cmds:  []*bonzai.Cmd{fooCmd, fooCmd.WithName(`foo2`)},
		Def:   fooCmd,
	}

	Cmd.Run(`foo`) // for default
	out, _ := mark.Bonzai(Cmd)
	fmt.Println(string(out))

	// Output:
	// # NAME
	//
	// `help-test` - just a help test
	//
	// # USAGE
	//
	//     h|ht|help-test COMMAND|some|-y|--yaml
	//
	// # COMMANDS
	//
	//     foo - (default)
	//     foo2
	//
}

func ExampleAKA() {

	var help = &bonzai.Cmd{
		Name:  `help`,
		Alias: `h|-h|-help|--help|/?`,
	}

	fmt.Println(mark.AKA(help))
	out, _ := mark.Fill(help, mark.Map, `The {{aka .}} command.`)
	fmt.Println(out)

	// Output:
	// `help` (`h`|`-h`|`-help`|`--help`|`/?`)
	// The `help` (`h`|`-h`|`-help`|`--help`|`/?`) command.
}

func ExampleCode() {

	long := `
I'll use the {{code "{{code}}"}} thing instead with something like
{{code "go mod init"}} since cannot use backticks.`

	fmt.Println(mark.Code(`go mod init`))
	out, _ := mark.Fill(nil, mark.Map, long)
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
	fmt.Println(mark.Commands(Cmd))

	// Output:
	// foo      - foo this command (default)
	//   subfoo - under the foo command
	// bar      - bar this command
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
	fmt.Println(mark.Commands(Cmd))

	// Output:
	// foo      - foo this command (default)
	//   subfoo - under the foo command
	// bar      - bar this command

}

func ExampleSummary() {

	var Cmd = &bonzai.Cmd{
		Name:  `mycmd`,
		Short: `my command short summary`,
		Do:    bonzai.Nothing,
	}

	fmt.Print(mark.Summary(Cmd))

	// Output:
	// `mycmd` - my command short summary
}
