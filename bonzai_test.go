package bonzai_test

import (
	"fmt"
	"io"

	"github.com/rwxrob/bonzai"
)

func ExampleCmd_AliasSlice() {
	barCmd := &bonzai.Cmd{
		Name:  `bar`,
		Alias: `b|rab`,
		Call: func(_ *bonzai.Cmd, _ ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}
	fmt.Printf("%q", barCmd.AliasSlice())

	fooCmd := &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd},
	}

	fmt.Printf("%q", fooCmd.AliasSlice())

	// Output:
	// ["b" "rab"][]
}

func ExampleCmd_Can() {
	barCmd := &bonzai.Cmd{
		Name: `bar`,
		Call: func(_ *bonzai.Cmd, _ ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}

	fooCmd := &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd},
	}

	fmt.Println(fooCmd.Can(`bar`))

	// Output:
	// bar
}

func ExampleCmd_WithName() {
	barCmd := &bonzai.Cmd{
		Name: `bar`,
		Call: func(_ *bonzai.Cmd, _ ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}

	fooCmd := barCmd.WithName(`foo`)
	fmt.Println(barCmd.Name)
	barCmd.Call(barCmd)
	fmt.Println(fooCmd.Name)
	fooCmd.Call(fooCmd)

	// Output:
	// bar
	// i am bar
	// foo
	// i am bar
}

func ExampleCmd_Tree() {

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
	fmt.Println(Cmd.CmdTree())

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

func ExampleCmd_Hides() {

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
		//Cmds:  []*bonzai.Cmd{subFooCmd, subFooHiddenCmd},
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
