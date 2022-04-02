package help_test

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/help"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/term"
	"github.com/rwxrob/term/esc"
)

func ExampleFormat_remove_Initial_Blanks() {
	fmt.Printf("%q\n", help.Format("\n   \n\n  \n   some"))
	// Output:
	// "some"
}

func ExampleFormat_wrapping() {
	fmt.Println(help.Format(`
Here is a bunch of stuff just to fill the line beyond 80 columns so that it will wrap when it is supposed to and right now
as well if there was a hard return in the middle of a line.
`))
	// Output:
	// Here is a bunch of stuff just to fill the line beyond 80 columns so that it will
	// wrap when it is supposed to and right now
	// as well if there was a hard return in the middle of a line.
}

func ExampleEmphasize() {

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

	fmt.Printf("%q\n", help.Emphasize("*italic*"))
	fmt.Printf("%q\n", help.Emphasize("**bold**"))
	fmt.Printf("%q\n", help.Emphasize("**bolditalic**"))
	fmt.Printf("%q\n", help.Emphasize("<under>"))

	// Output:
	// "\x1b[4mitalic\x1b[0m"
	// "\x1b[33mbold\x1b[0m"
	// "\x1b[33mbolditalic\x1b[0m"
	// "<\x1b[4munder\x1b[0m>"
}

func ExampleCmd_name() {
	x := &Z.Cmd{
		Name: "foo",
	}
	help.ForTerminal(x, "name")
	// Output:
	// foo
}

func ExampleCmd_summary() {
	x := &Z.Cmd{
		Summary: `foo all the things`,
	}
	help.ForTerminal(x, "summary")
	// Output:
	// foo all the things
}

func ExampleCmd_other() {
	x := &Z.Cmd{
		Other: map[string]string{"foo": `some foo text`},
	}
	help.ForTerminal(x, "foo")
	// Output:
	// some foo text
}

func ExampleCmd_all_Commands_and_Params() {
	x := &Z.Cmd{
		Name:    `foo`,
		Params:  []string{"p1", "p2"},
		Version: "v0.1.0",
		Summary: `foo all the things`,
		Other:   map[string]string{"foo": `some foo text`},
		Call:    func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// some foo text
}
