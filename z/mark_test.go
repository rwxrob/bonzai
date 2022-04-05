package Z_test

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/term"
)

func ExampleEmph_basics() {

	// Emph observes the rwxrob/term escapes
	// (see package documentation for more)

	term.Italic = `<italic>`
	term.Bold = `<bold>`
	term.BoldItalic = `<bolditalic>`
	term.Under = `<under>`
	term.Reset = `<reset>`

	fmt.Println(Z.Emph("*ITALIC*"))
	fmt.Println(Z.Emph("**BOLD**"))
	fmt.Println(Z.Emph("***BOLDITALIC***"))
	fmt.Println(Z.Emph("<UNDER>")) // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExampleWrap() {
	defer func() { Z.Columns = Z.Columns }()
	Z.Columns = 10
	fmt.Println(Z.Wrap(`some thing here that is more than 10 characters`))
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExampleIndent() {
	defer func() { Z.IndentBy = Z.IndentBy }()
	Z.IndentBy = 4
	fmt.Printf("%q", Z.Indent("some\nthat is \n  indented"))
	// Output:
	// "    some\n    that is \n      indented\n"
}

func ExampleInWrap() {
	defer func() { Z.IndentBy = Z.IndentBy }()
	defer func() { Z.Columns = Z.Columns }()
	Z.IndentBy = 4
	Z.Columns = 10
	fmt.Printf("%q", Z.InWrap("some\nthat is \n  indented"))
	// Output:
	// "    some\n    that\n    is\n    indented\n"
}

func ExampleBlocks_bulleted() {
	in := `

			* some thing
			* another thing

			* another block
			* here

			`

	fmt.Println(Z.Blocks(in)[1])

	//Output:
	// * another block
	// * here
}

func ExampleBlocks_numbered() {
	in := `

			1. some thing
			2. another thing

			1. another block
			2. here

			`

	fmt.Println(Z.Blocks(in)[1])

	//Output:
	// 1. another block
	// 2. here
}

func ExampleBlocks_paragraph() {
	in := `
			Simple paragraph
			here on multiple
			lines

			And another one here
			with just a bit more.

			`

	fmt.Println(Z.Blocks(in)[1])

	// Output:
	// And another one here with just a bit more.
}

func ExampleBlocks_verbatim() {

	// Note that the following begins consistently with three tabs so that
	// dedenting works consistently. There are four spaces before Now and
	// the verbatim block. Notice that even the blank line within the
	// verbatim block must have the exact same indentation and spaced
	// verbatim prefix. (If using Vi/m try set :list to display them.)

	in := `
			Must have another block type first.

			     Now we can start
			     a Verbatim
			     block.
			     
			     Which can have blank lines, even.

			And back to a paragraph block.

			`

	fmt.Printf("%q\n", Z.Blocks(in)[0])
	fmt.Printf("%q\n", Z.Blocks(in)[1])
	fmt.Printf("%q\n", Z.Blocks(in)[2])

	//Output:
	// "Must have another block type first."
	// "Now we can start\na Verbatim\nblock.\n\nWhich can have blank lines, even."
	// "And back to a paragraph block."

}

// Now we can start
// a Verbatim
// block.
//
// Which can have blank lines, even.

/*
func ExampleBlocks() {
	in := `

		 Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
		 eiusmod tempor incididunt ut labore et dolore magna aliqua.

		 Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
		 nisi ut aliquip ex ea commodo consequat.

		 Duis aute irure dolor in reprehenderit in voluptate velit esse
		 cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
		 cupidatat non proident, sunt in culpa qui officia deserunt mollit
		 anim id est laborum.

		 `

	fmt.Println(Z.Blocks(in))

	// Output:
	// some thing
}
*/

/*
func ExampleFormat_remove_Initial_Blanks() {
	fmt.Printf("%q\n", Z.Format("\n   \n\n  \n   some"))
	// Output:
	// "some"
}

func ExampleFormat_wrapping() {
	fmt.Println(Z.Format(`
Here is a bunch of stuff just to fill the line beyond 80 columns so that it will wrap when it is supposed to and right now
as well if there was a hard return in the middle of a line.
`))
	// Output:
	// Here is a bunch of stuff just to fill the line beyond 80 columns so that it will
	// wrap when it is supposed to and right now
	// as well if there was a hard return in the middle of a line.
}
*/
