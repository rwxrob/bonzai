package mark_test

import (
	"fmt"
	"io"
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

func ExampleRenderString() {

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
	zmark := `
	{{hello}}, my name is {{.Name}} with {{.Count}}. Summary: {{.Summary}}`
	funcs := template.FuncMap{}
	funcs[`hello`] = func() string { return `Hello` }

	out, err := mark.Render(thing, funcs, zmark)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(out)

	// Output:
	// Hello, my name is Thing with 20. Summary: Thing 20

}

func ExampleCmdTree() {
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

	//  ─ │ ┐ ┘ ┌ └ ├ ┤ ┬ ┴ ┼

	Cmd.SetCallers()
	fmt.Print("# Synopsis\n\n")
	fmt.Println(mark.CmdTree(Cmd))

	// Output:
	// # Synopsis
	//
	//     mycmd      ← my command short summary
	//     ├─foo      ← foo this command
	//     │ └─subfoo ← under the foo command
	//     └─bar      ← bar this command
}

func ExampleUsage_noInteractiveTerminal() {
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
		Long: `
			Here is a long description.
			On multiple lines.`,
	}

	r, err := mark.Usage(Cmd)
	if err != nil {
		fmt.Println(err)
	}
	out, _ := io.ReadAll(r)
	fmt.Println(string(out))

	// Output:
	// # Usage
	//
	//     mycmd      ← my command short summary
	//       foo      ← foo this command
	//         subfoo ← under the foo command
	//       bar      ← bar this command
	//
	// Here is a long description.
	// On multiple lines.
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

	r, err := mark.Usage(Cmd)
	if err != nil {
		fmt.Println(err)
	}

	out, _ := io.ReadAll(r)
	fmt.Println(string(out))

	// Output:
	// # Usage
	//
	//     mycmd      ← my command short summary
	//       foo      ← foo this command
	//         subfoo ← under the foo command
	//       bar      ← bar this command
	//
	// Here is a long description.
	// On multiple lines.
}
