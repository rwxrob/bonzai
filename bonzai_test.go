package bonzai_test

import (
	"fmt"
	"sync"

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

func ExampleCmd_SeekInit() {

	var foo = &bonzai.Cmd{
		Name: `foo`,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am foo`)
			return nil
		},
	}

	var bar = &bonzai.Cmd{
		Name: `bar`,
		Cmds: []*bonzai.Cmd{foo},
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am bar`)
			return nil
		},
	}

	var baz = &bonzai.Cmd{
		Name:   `baz`,
		Cmds:   []*bonzai.Cmd{bar},
		NoArgs: true,
		Do: func(*bonzai.Cmd, ...string) error {
			fmt.Println(`I am baz`)
			return nil
		},
	}

	fmt.Println(baz.SeekInit(`bar`))
	fmt.Println(baz.SeekInit(`bar`, `arg1`))

	// Output:
	// bar [] <nil>
	// bar [arg1] <nil>

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
	// developer error: Cmd.Vers (foo) length must be less than 50 runes

}

func ExampleErrInvalidShort() {
	var foo = &bonzai.Cmd{
		Name:  `foo`,
		Short: `this is a long short desc that is longer than 50 runes`,
	}

	err := foo.Run()
	fmt.Println(err)

	// Output:
	// developer error: Cmd.Short (foo) length must be less than 50 runes and must begin with a lowercase letter

}

func ExampleCmd_Def() {

	var foo = &bonzai.Cmd{
		Name: `foo`,
		Do:   bonzai.Nothing,
	}

	var Cmd = &bonzai.Cmd{
		Name: `mycmd`,
		Def:  foo,
	}

	x, args, err := Cmd.SeekInit(`arg1`, `arg2`)
	fmt.Println(x, args, err)
	fmt.Println(x.Caller())

	// Output:
	// foo [arg1 arg2] <nil>
	// mycmd

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
	// developer error: Cmd.Short (foo) length must be less than 50 runes and must begin with a lowercase letter
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
	// usage error: arg #1 must match: ^o+$
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

func ExampleCmd_Init() {

	printname := func(x *bonzai.Cmd, _ ...string) error {
		fmt.Print(x.Name + " ")
		return nil
	}

	var Cmd = &bonzai.Cmd{
		Name: `foo`,
		Init: printname,
		Cmds: []*bonzai.Cmd{
			{
				Name: `nothing`,
				Init: printname,
				Do:   bonzai.Nothing,
			},
			{
				Name: `other`,
				Init: printname,
				Do:   bonzai.Nothing,
			},
		},
	}

	if err := Cmd.Run(`nothing`); err != nil {
		fmt.Println(err)
	}

	if err := Cmd.Run(`other`); err != nil {
		fmt.Println(err)
	}

	// Output:
	// foo nothing foo other
}

func ExampleWalkDeep_hidden() {
	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var sssh = &bonzai.Cmd{
		Name: `sssh`,
		Do:   bonzai.Nothing,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		Short: `foo this command`,
		Cmds:  []*bonzai.Cmd{subFooCmd, sssh.AsHidden()},
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

	printname := func(level int, x *bonzai.Cmd) error {
		if x.IsHidden() {
			return nil
		}
		for range level {
			fmt.Print("  ")
		}
		fmt.Println(x.Name)
		return nil
	}

	Cmd.WalkDeep(printname, nil)

	// Output:
	// mycmd
	//   foo
	//     subfoo
	//   bar

}

func ExampleCmd_PathDashed() {
	var subFooCmd = &bonzai.Cmd{
		Name:  `subfoo`,
		Alias: `sf`,
		Short: `under the foo command`,
	}

	var sssh = &bonzai.Cmd{
		Name: `sssh`,
		Do:   bonzai.Nothing,
	}

	var fooCmd = &bonzai.Cmd{
		Name:  `foo`,
		Alias: `f`,
		Short: `foo this command`,
		Cmds:  []*bonzai.Cmd{subFooCmd, sssh.AsHidden()},
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

	cmd, args, err := Cmd.SeekInit(`foo`, `sssh`)
	fmt.Println(cmd, args, err)
	fmt.Println(cmd.PathDashed())

	// Output:
	// sssh [] <nil>
	// foo-sssh

}

type InMem struct {
	sync.Mutex
	m map[string]string
}

func (p *InMem) Setup() error {
	p.m = make(map[string]string)
	return nil
}
func (p *InMem) Set(k, v string) {
	p.Lock()
	defer p.Unlock()
	fmt.Println(`setting:`, k, v)
	p.m[k] = v
}
func (p *InMem) Get(k string) string {
	p.Lock()
	defer p.Unlock()
	fmt.Println(`getting:`, k)
	return p.m[k]
}

func ExampleCmd_Persist() {

	/*
		type InMem struct {
			sync.Mutex
			m map[string]string
		}

		func (p *InMem) Setup() error {
			p.m = make(map[string]string)
			return nil
		}
		func (p *InMem) Set(k, v string) {
			p.Lock()
			defer p.Unlock()
			fmt.Println(`setting:`, k, v)
			p.m[k] = v
		}
		func (p *InMem) Get(k string) string {
			p.Lock()
			defer p.Unlock()
			fmt.Println(`getting:`, k)
			return p.m[k]
		}
	*/

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`, P: true}},
		Pers: new(InMem),
		Do:   bonzai.Nothing,
	}

	// persistence doesn't get setup without
	Cmd.SeekInit(`cmd`)

	fmt.Println(Cmd.Get(`key`))
	Cmd.Set(`key`, `other`)
	fmt.Println(Cmd.Get(`key`))

	// Output:
	// getting: key
	// setting: key value
	// value
	// setting: key other
	// getting: key
	// other

}

func ExampleCmd_Vars_inMemory() {

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Pers: new(InMem),
		Do:   bonzai.Nothing,
	}

	Cmd.Run()

	fmt.Println(Cmd.Get(`key`))
	Cmd.Set(`key`, `other`)
	fmt.Println(Cmd.Get(`key`))

	// Output:
	// value
	// other

}

func ExampleCmd_Vars_undeclaredPanic() {

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Pers: new(InMem),
		Do:   bonzai.Nothing,
	}

	Cmd.Run() // vars don't get cached without

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	fmt.Println(Cmd.Get(`undeclared`))

	// Output:
	// Recovered from panic: not declared in Vars: undeclared

}

func ExampleCmd_Var() {

	var subCmd = &bonzai.Cmd{
		Name: `subcmd`,
		Vars: bonzai.Vars{{K: `key`, V: `other`}},
		Do:   bonzai.Nothing,
	}

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Cmds: []*bonzai.Cmd{subCmd},
	}

	err := Cmd.Run(`subcmd`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Cmd.Get(`key`))
	fmt.Println(subCmd.Get(`key`))

	// Output:
	// value
	// other

}

func ExampleCmd_Vars_sameInherited() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	var subCmd = &bonzai.Cmd{
		Name: `subcmd`,
		Vars: bonzai.Vars{{I: `key`}},
		Do:   bonzai.Nothing,
	}

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Cmds: []*bonzai.Cmd{subCmd},
	}

	err := Cmd.Run(`subcmd`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(subCmd.Get(`key`))

	// Output:
	// value

}

func ExampleCmd_Vars_failedInherit() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	var subCmd = &bonzai.Cmd{
		Name: `subcmd`,
		Vars: bonzai.Vars{{I: `goodluck`}},
		Do:   bonzai.Nothing,
	}

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Cmds: []*bonzai.Cmd{subCmd},
	}

	err := Cmd.Run(`subcmd`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(subCmd.Get(`key`))

	// Output:
	// Recovered from panic: failed to find inherited Var: goodluck
}

func ExampleCmd_Vars_inheritDeep() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	var subsubCmd = &bonzai.Cmd{
		Name: `subsubcmd`,
		Vars: bonzai.Vars{{I: `key`}},
		Do:   bonzai.Nothing,
	}

	var subCmd = &bonzai.Cmd{
		Name: `subcmd`,
		Cmds: []*bonzai.Cmd{subsubCmd},
	}

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Cmds: []*bonzai.Cmd{subCmd},
	}

	err := Cmd.Run(`subcmd`, `subsubcmd`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(subsubCmd.Get(`key`))

	// Output:
	// value
}

func ExampleCmd_Vars_inheritDeepSet() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()

	var subsubCmd = &bonzai.Cmd{
		Name: `subsubcmd`,
		Vars: bonzai.Vars{{I: `key`}},
		Do:   bonzai.Nothing,
	}

	var subCmd = &bonzai.Cmd{
		Name: `subcmd`,
		Cmds: []*bonzai.Cmd{subsubCmd},
	}

	var Cmd = &bonzai.Cmd{
		Name: `cmd`,
		Vars: bonzai.Vars{{K: `key`, V: `value`}},
		Cmds: []*bonzai.Cmd{subCmd},
	}

	err := Cmd.Run(`subcmd`, `subsubcmd`)
	if err != nil {
		fmt.Println(err)
	}
	subsubCmd.Set(`key`, `newvalue`)
	fmt.Println(subsubCmd.Get(`key`))

	// Output:
	// newvalue
}
