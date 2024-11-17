package bonzai_test

import (
	"fmt"
	"os"

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
	// Short length >50 or not lower for "foo": "this is a long short desc that is longer than 50 runes"
}

func ExampleValidate() {

	var after = &bonzai.Cmd{
		Name: `after`,
		Cmds: []*bonzai.Cmd{&bonzai.Cmd{Name: `some`}},
	}

	var foo = &bonzai.Cmd{
		Name:  `foo`,
		Short: `this is a long short desc that is longer than 50 runes`,
		Cmds:  []*bonzai.Cmd{after},
	}

	var Cmd = &bonzai.Cmd{
		Name: `main`,
		Cmds: []*bonzai.Cmd{foo},
	}

	err := foo.Validate()
	fmt.Println(err)

	err = Cmd.Validate() // does not check anything but self
	fmt.Println(err)

	// Output:
	// Short length >50 or not lower for "foo": "this is a long short desc that is longer than 50 runes"
	// <nil>
}

func ExampleErrInvalidArg() {
	var foo = &bonzai.Cmd{
		Name:     `foo`,
		Short:    `just oooo`,
		RegxArgs: `^o+$`,
		Cmds:     []*bonzai.Cmd{&bonzai.Cmd{Name: `foo`}},
	}

	err := foo.Run(`fooo`)
	fmt.Println(err)

	// Output:
	// arg #1 must match: ^o+$
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

	names := []string{} // enclosed

	aggregate := func(level int, x *bonzai.Cmd) error {
		names = append(names, fmt.Sprintf("%v-%v", x.Name, level))
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

	names := []string{} // enclosed

	aggregate := func(level int, x *bonzai.Cmd) error {
		names = append(names, fmt.Sprintf("%v-%v", x.Name, level))
		return nil
	}

	errors := []error{} // enclosed
	onerror := func(err error) {
		errors = append(errors, err)
	}

	Cmd.WalkWide(aggregate, onerror)
	fmt.Println(names)

	// Output:
	// [top-0 foo-1 foo2-1 bar-1 bar2-1 bar-1 bar2-1]

}

func ExampleValidate_callableDefOnly() {

	var Cmd = &bonzai.Cmd{
		Name: `foo`,
		Def:  &bonzai.Cmd{Name: `something`},
	}

	fmt.Println(Cmd.Validate())

	// Output:
	// <nil>
}

func ExampleSeek() {

	var after = &bonzai.Cmd{
		Name: `after`,
		Cmds: []*bonzai.Cmd{&bonzai.Cmd{Name: `some`}},
	}

	var foo = &bonzai.Cmd{
		Name: `foo`,
		Cmds: []*bonzai.Cmd{after},
	}

	var Cmd = &bonzai.Cmd{
		Name: `main`,
		Cmds: []*bonzai.Cmd{foo},
	}

	cmd, args := Cmd.Seek(`foo`, `after`, `some`, `arg1`, `arg2`)
	fmt.Println(`cmd:`, cmd)
	fmt.Println(`args:`, args)

	// Output:
	// cmd: some
	// args: [arg1 arg2]

}

func ExampleCmd_Env() {
	defer os.Unsetenv(`SOME`)

	var Cmd = &bonzai.Cmd{
		Name: `foo`,
		Env:  map[string]string{`SOME`: `thing`},
		Def:  &bonzai.Cmd{Name: `nothing`},
	}

	Cmd.Run()
	fmt.Println(os.Getenv(`SOME`))

	// Output:
	// thing
}

func ExampleCmd_Init() {

	printname := func(x *bonzai.Cmd, _ ...string) error {
		fmt.Print(x.Name + " ")
		return nil
	}

	var Cmd = &bonzai.Cmd{
		Name: `foo`,
		Init: printname,
		Cmds: []*bonzai.Cmd{
			{Name: `nothing`, Init: printname},
			{Name: `other`, Init: printname},
		},
	}

	Cmd.Run(`nothing`)
	Cmd.Run(`other`)

	// Output:
	// foo nothing foo other
}
