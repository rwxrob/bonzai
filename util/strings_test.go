package util_test

import (
	"github.com/rwxrob/bonzai/each"
	"github.com/rwxrob/bonzai/util"
)

func ExampleLines() {
	buf := `
some

thing 
here

mkay
`
	each.Print(util.Lines(buf))
	// Output:
	// something heremkay
}
