package bonzai_test

import (
	"fmt"

	bonzai "github.com/rwxrob/bonzai/pkg"
)

func ExampleCmd_AliasesSlice() {

	barCmd := &bonzai.Cmd{
		Name:    `bar`,
		Aliases: `b|rab`,
		Call: func(_ *bonzai.Cmd, _ ...string) error {
			fmt.Println(`i am bar`)
			return nil
		},
	}
	fmt.Printf("%q", barCmd.AliasesSlice())

	fooCmd := &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd},
	}

	fmt.Printf("%q", fooCmd.AliasesSlice())

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
