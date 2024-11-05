package bonzai_test

import (
	"fmt"

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
