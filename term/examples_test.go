package term_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/term"
)

func ExampleRed() {
	term.Red = "<red>"
	term.Reset = "<reset>"
	fmt.Println(term.Red + "simply red" + term.Reset)
	term.AttrOff()
	fmt.Println(term.Red + "simply red" + term.Reset)
	// Output:
	// <red>simply red<reset>
	// simply red
}
