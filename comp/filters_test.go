package comp_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/rwxrob/bonzai/comp"
)

func ExampleWithMaxLen() {
	fooCmd := &bonzai.Cmd{Name: "foo"}
	barCmd := &bonzai.Cmd{Name: "bar"}
	blahCmd := &bonzai.Cmd{Name: "blah"}
	cmd := bonzai.Cmd{
		Name: "cmd",
		Cmds: []*bonzai.Cmd{fooCmd, barCmd, blahCmd},
	}

	fmt.Println(comp.WithMaxLen(3, comp.Cmds).Complete(cmd))

	// Output:
	// [foo bar]
}

func ExampleWithMinLen() {
	fooCmd := &bonzai.Cmd{Name: "foo"}
	barCmd := &bonzai.Cmd{Name: "bar"}
	blahCmd := &bonzai.Cmd{Name: "blah"}
	cmd := bonzai.Cmd{
		Name: "cmd",
		Cmds: []*bonzai.Cmd{fooCmd, barCmd, blahCmd},
	}

	fmt.Println(comp.WithMinLen(4, comp.Cmds).Complete(cmd))

	// Output:
	// [blah]
}

func ExampleWithPrefix() {
	fooCmd := &bonzai.Cmd{Name: "foo"}
	barCmd := &bonzai.Cmd{Name: "bar"}
	blahCmd := &bonzai.Cmd{Name: "blah"}
	cmd := bonzai.Cmd{
		Name: "cmd",
		Cmds: []*bonzai.Cmd{fooCmd, barCmd, blahCmd},
	}

	fmt.Println(comp.WithPrefix("b", comp.Cmds).Complete(cmd))
	fmt.Println(comp.WithPrefix("b", comp.Cmds).Complete(cmd, "a"))

	// Output:
	// [bar blah]
	// [bar]
}
