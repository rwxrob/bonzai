package term_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/z/term"
)

func ExampleWinSize() {
	fmt.Println(term.WinSize)
	// always zeros because not an interactive terminal

	// Output:
	// {0 0 0 0}
}
