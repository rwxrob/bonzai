package futil_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/futil"
)

/*
func ExampleMultipart_String() {
	out := futil.Multipart{
		`stdout`:  "some standard output\n\non multiple lines",
		`stderr`:  "some standard err on single line",
		`exitval`: "-1",
	}
	fmt.Print(out)
	// Output:
	// ignored
}
*/

/*
func ExampleMultipart_UnmarshalText() {
	out := futil.Multipart{Map: map[string]string{`dummy`: `just checking`}}
	buf := `NNQLO9MP27BRECLC6CED8QC2RGHQPHRL stdout
some standard output

on multiple lines
NNQLO9MP27BRECLC6CED8QC2RGHQPHRL stderr
some standard err on single line
NNQLO9MP27BRECLC6CED8QC2RGHQPHRL exitval
-1
NNQLO9MP27BRECLC6CED8QC2RGHQPHRL end
`
	err := out.UnmarshalText([]byte(buf))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(out)
	// Output:
	// ignored
}
*/

func ExampleMultipart_UnmarshalText_explicit() {
	out := futil.Multipart{
		Delimiter: `IMMADELIM`,
		Map:       map[string]string{`dummy`: `just checking`},
	}
	buf := `
random
ignored
lines
here
IMMADELIM stdout
some standard output

on multiple lines
IMMADELIM stderr
some standard err on single line
IMMADELIM exitval
-1
IMMADELIM break
`
	err := out.UnmarshalText([]byte(buf))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(out)
	// Unordered Output:
	// IMMADELIM stdout
	// some standard output
	//
	// on multiple lines
	// IMMADELIM stderr
	// some standard err on single line
	// IMMADELIM exitval
	// -1
	// IMMADELIM break

}
