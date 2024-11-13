package bonzai_test

import (
	"fmt"
	"io"

	"github.com/rwxrob/bonzai"
)

func ExampleCmd_Aliases() {
	var barCmd = &bonzai.Cmd{
		Name:  `bar`,
		Alias: `b|rab`,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}
	fmt.Printf("%q", barCmd.Aliases())

	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd},
	}

	fmt.Printf("%q", fooCmd.Aliases())

	// Output:
	// ["b" "rab" "bar"]["foo"]
}

func ExampleCmd_Can() {
	var barCmd = &bonzai.Cmd{
		Name: `bar`,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}

	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd},
	}

	fmt.Println(fooCmd.Can(`bar`))

	// Output:
	// bar
}

func ExampleCmd_WithName() {
	var barCmd = &bonzai.Cmd{
		Name: `bar`,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}

	fooCmd := barCmd.WithName(`foo`)
	fmt.Println(barCmd.Name)
	barCmd.Do(barCmd)
	fmt.Println(fooCmd.Name)
	fooCmd.Do(fooCmd)

	// Output:
	// bar
	// i am bar
	// foo
	// i am bar
}

func ExampleCmd_CmdTreeString() {
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
	}

	fmt.Print("# Synopsis\n\n")
	fmt.Println(Cmd.CmdTreeString())

	// Output:
	// # Synopsis
	//
	//     mycmd      ← my command short summary
	//       foo      ← foo this command
	//         subfoo ← under the foo command
	//       bar      ← bar this command
}

func ExampleCmd_Mark_noInteractiveTerminal() {
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

	out, _ := io.ReadAll(Cmd.Mark())
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

func ExampleCmd_AsHidden() {
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

	out, _ := io.ReadAll(Cmd.Mark())
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

func ExampleCmd_Run() {
	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am foo`)
			return nil
		},
	}

	var barCmd = &bonzai.Cmd{
		Name: `bar`,
		Cmds: []*bonzai.Cmd{fooCmd},
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am bar`)
			return nil
		},
	}

	var bazCmd = &bonzai.Cmd{
		Name: `baz`,
		Cmds: []*bonzai.Cmd{barCmd},
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am baz`)
			return nil
		},
	}

	fooCmd.Run()
	bazCmd.Run("bar")
	bazCmd.Run("bar", "foo")

	// Output:
	// I am foo
	// I am bar
	// I am foo
}

func ExampleErrInvalidVers() {
	var foo = &bonzai.Cmd{
		Name: `foo`,
		Vers: `this is a long version that is longer than 50 runes`,
	}

	err := foo.Run()
	fmt.Println(err)

	// Output:
	// Cmd.Vers length >50 for "foo": "this is a long version that is longer than 50 runes"
}

func ExampleErrInvalidShort() {
	var foo = &bonzai.Cmd{
		Name:  `foo`,
		Short: `this is a long short desc that is longer than 50 runes`,
	}

	err := foo.Run()
	fmt.Println(err)

	// Output:
	// Cmd.Short length >50 for "foo": "this is a long short desc that is longer than 50 runes"
}
