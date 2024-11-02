package main

import (
	"fmt"

	"github.com/rwxrob/bonzai/term"
)

func main() {
	fmt.Print("Enter sample text: ")
	p := term.Read()
	fmt.Println("You entered: ", p)
	fmt.Print("Enter sample hidden text: ")
	p = term.ReadHide()
	fmt.Println()
	fmt.Println("You entered while hidden: ", p)

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
