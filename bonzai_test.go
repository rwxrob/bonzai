package bonzai_test

import (
	"fmt"

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

func ExampleCmd_WalkDeep() {

	var barCmd = &bonzai.Cmd{Name: `bar`}

	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd, barCmd.WithName(`bar2`)},
	}

	var Cmd = &bonzai.Cmd{
		Name: `top`,
		Cmds: []*bonzai.Cmd{fooCmd, fooCmd.WithName(`foo2`)},
	}

	Cmd.SetCallers() // no Run/Exec/Seek, so explicit

	names := []string{} // enclosed

	aggregate := func(x *bonzai.Cmd) error {
		names = append(names, fmt.Sprintf("%v-%v", x.Name, x.Level()))
		return nil
	}

	errors := []error{} // enclosed
	onerror := func(err error) {
		errors = append(errors, err)
	}

	Cmd.WalkDeep(aggregate, onerror)
	fmt.Println(names)

	// Output:
	// [top-0 foo-1 bar-2 bar2-2 foo2-1 bar-2 bar2-2]

}

func ExampleCmd_WalkWide() {

	var barCmd = &bonzai.Cmd{Name: `bar`}

	var fooCmd = &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{barCmd, barCmd.WithName(`bar2`)},
	}

	var Cmd = &bonzai.Cmd{
		Name: `top`,
		Cmds: []*bonzai.Cmd{fooCmd, fooCmd.WithName(`foo2`)},
	}

	Cmd.SetCallers() // no Run/Exec/Seek, so explicit

	names := []string{} // enclosed

	aggregate := func(x *bonzai.Cmd) error {
		names = append(names, fmt.Sprintf("%v-%v", x.Name, x.Level()))
		return nil
	}

	errors := []error{} // enclosed
	onerror := func(err error) {
		errors = append(errors, err)
	}

	Cmd.WalkWide(aggregate, onerror)
	fmt.Println(names)

	// Output:
	// [top-0 foo-1 foo2-1 bar-2 bar2-2 bar-2 bar2-2]

}
