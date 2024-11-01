// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

package term_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/term"
	"github.com/rwxrob/bonzai/term/esc"
)

func ExampleRed() {
	term.Red = "<red>"
	term.Reset = "<reset>"
	fmt.Println(term.Red + "simply red" + term.Reset)
	defer term.AttrOn()
	term.AttrOff()
	fmt.Println(term.Red + "simply red" + term.Reset)
	// Output:
	// <red>simply red<reset>
	// simply red
}

func ExampleStripNonPrint() {
	some := esc.Bold + "not bold" + esc.Reset
	fmt.Println(term.StripNonPrint(some))
	// Output;
	// not bold
}

func ExampleEmphFromLess() {

	/*
	   export LESS_TERMCAP_mb="[35m" # magenta
	   export LESS_TERMCAP_md="[33m" # yellow
	   export LESS_TERMCAP_me="" # "0m"
	   export LESS_TERMCAP_se="" # "0m"
	   export LESS_TERMCAP_so="[34m" # blue
	   export LESS_TERMCAP_ue="" # "0m"
	   export LESS_TERMCAP_us="[4m"  # underline
	*/

	os.Setenv("LESS_TERMCAP_mb", esc.Magenta)
	os.Setenv("LESS_TERMCAP_md", esc.Yellow)
	os.Setenv("LESS_TERMCAP_me", esc.Reset)
	os.Setenv("LESS_TERMCAP_se", esc.Reset)
	os.Setenv("LESS_TERMCAP_so", esc.Blue)
	os.Setenv("LESS_TERMCAP_ue", esc.Reset)
	os.Setenv("LESS_TERMCAP_us", esc.Under)

	term.EmphFromLess()

	fmt.Printf("%q\n", term.Italic+"italic"+term.Reset)
	fmt.Printf("%q\n", term.Bold+"bold"+term.Reset)
	fmt.Printf("%q\n", term.BoldItalic+"bolditalic"+term.Reset)
	fmt.Printf("%q\n", term.Under+"under"+term.Reset)

	// Output:
	// "\x1b[4mitalic\x1b[0m"
	// "\x1b[33mbold\x1b[0m"
	// "\x1b[35mbolditalic\x1b[0m"
	// "\x1b[4munder\x1b[0m"

}

func ExampleREPL() {
	defer term.TrapPanic()

	// both are enclosed in prompt/respond functions
	var history []string
	hcount := 1

	prompt := func(_ string) string {
		if hcount > 3 {
			panic("All done.")
		}
		return fmt.Sprintf("%v> ", hcount)
	}

	respond := func(in string) string {
		hcount++
		history = append(history, in)
		return "okay"
	}

	term.REPL(prompt, respond)
}
