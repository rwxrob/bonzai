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

func ExampleCmd_Title() {

	x := &bonzai.Cmd{
		Name:  `mycmd`,
		Alias: `my|cmd`,
		Short: `my command short summary`,
	}

	fmt.Println(x.Title())

	// Output:
	// mycmd (my|cmd) - my command short summary
}

func ExampleCmd_Title_noShort() {

	x := &bonzai.Cmd{
		Name:  `mycmd`,
		Alias: `my|cmd`,
	}

	fmt.Println(x.Title())

	// Output:
	// mycmd (my|cmd)
}

func ExampleCmd_Title_noShortNoAlias() {

	x := &bonzai.Cmd{
		Name: `mycmd`,
	}

	fmt.Println(x.Title())

	// Output:
	// mycmd
}

func ExampleCmd_Title_noName() {

	x := &bonzai.Cmd{
		Alias: `my|cmd`,
	}

	fmt.Println(x.Title())

	// Output:
	// NONAME (my|cmd)
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

	fmt.Println("# Synopsis")
	fmt.Println(Cmd.CmdTree())

	// Output:
	// # Synopsis
	//     foo (f)       - foo this command
	//       subfoo (sf) - under the foo command
	//     bar (b)       - bar this command

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
			On multiple lines.
		`,
	}

	//_, err := io.Copy(os.Stdout, Cmd.Mark())

	out, err := io.ReadAll(Cmd.Mark())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%q\n", out)

	// Output:
	// "# Name\n\nmycmd (my|cmd) - my command short summary\n\n# Synopsis\n\n    foo (f)       - foo this command\n      subfoo (sf) - under the foo command\n    bar (b)       - bar this command\n\n# Description\n\nHere is a long description.\t\nOn multiple lines.\n\t\t"

}
