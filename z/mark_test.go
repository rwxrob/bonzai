package Z_test

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/term"
)

func init() {
	term.Italic = `<italic>`
	term.Bold = `<bold>`
	term.BoldItalic = `<bolditalic>`
	term.Under = `<under>`
	term.Reset = `<reset>`
}

func ExampleLines() {
	fmt.Printf("%q\n", Z.Lines("line one\nline two"))
	// Output:
	// ["line one" "line two"]
}

// --------------------------- start Blocks ---------------------------

func ExampleBlocks_bulleted() {
	in := `

			* some thing
			* another thing

			* another block
			* here

			*boldnotbullet*

			`

	blocks := Z.Blocks(in)
	fmt.Printf("%v %q\n", blocks[1].T, blocks[1])
	fmt.Printf("%v %q\n", blocks[2].T, blocks[2])

	//Output:
	// 3 "* another block\n* here"
	// 1 "*boldnotbullet*"
}

func ExampleBlocks_numbered() {
	in := `

			1. some thing
			2. another thing

			1. another block
			2. here

			`

	fmt.Printf("%q\n", Z.Blocks(in)[1])

	//Output:
	// "1. another block\n2. here"
}

func ExampleBlocks_paragraph() {

	// defer func() { scan.Trace = 0 }()
	// log.SetFlags(0)
	// scan.Trace = 1

	in := `
			Simple paragraph
			here on multiple
			lines.

			And another one here
			with just a bit more.

			`

	fmt.Printf("%q\n", Z.Blocks(in)[0])
	fmt.Printf("%q\n", Z.Blocks(in)[1])

	// Output:
	// "Simple paragraph here on multiple lines."
	// "And another one here with just a bit more."
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

// -------------------------- main BonzaiMark -------------------------

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
	col := Z.Columns
	Z.Columns = 10
	fmt.Println(Z.Wrap(`some thing here that is more than 10 characters`))
	Z.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExampleIndent() {
	indent := Z.IndentBy
	col := Z.Columns
	Z.Columns = 10
	Z.Columns = 10
	Z.IndentBy = 4
	fmt.Printf("%q", Z.Indent("some\nthat is \n  indented"))
	Z.IndentBy = indent
	Z.Columns = col
	// Output:
	// "    some\n    that is \n      indented\n"
}

func ExampleInWrap() {
	defer func() { Z.IndentBy = Z.IndentBy }()
	indent := Z.IndentBy
	col := Z.Columns
	Z.Columns = 10
	Z.IndentBy = 4
	fmt.Printf("%q", Z.InWrap("some\nthat is \n  indented"))
	Z.IndentBy = indent
	Z.Columns = col
	// Output:
	// "    some\n    that\n    is\n    indented\n"
}

func ExampleMark_simple() {

	fmt.Print(Z.Mark(`**foo**`))

	//Output:
	// <bold>foo<reset>

}

func ExampleMark() {

	in := `
			Must have *another* block before verbatim:

			     Now we can start
			     a Verbatim
			     block.

			     Which can have blank lines, even.

			And back to a paragraph block.

			* **foo**
			* bar

			And a numbered list

			1. Something
			2. here

			That's really it.

			`

	fmt.Println("----------------------")
	fmt.Print(Z.Mark(in))
	fmt.Println("----------------------")

	//Output:
	// ----------------------
	//        Must have <italic>another<reset> block before verbatim:
	//
	//            Now we can start
	//            a Verbatim
	//            block.
	//
	//            Which can have blank lines, even.
	//
	//        And back to a paragraph block.
	//
	//        * <bold>foo<reset>
	//        * bar
	//
	//        And a numbered list
	//
	//        1. Something
	//        2. here
	//
	//        That's really it.
	//
	// ----------------------

}

// ------------------------ Sprintf variations ------------------------

func ExampleEmphf() {
	fmt.Println(Z.Emphf(`some *%v* thing`, "italic"))
	// Output:
	// some <italic>italic<reset> thing
}

func ExampleWrapf() {
	col := Z.Columns
	Z.Columns = 3
	fmt.Println(Z.Wrapf(`some %v here`, 10))
	Z.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExampleIndentf() {
	in := Z.IndentBy
	Z.IndentBy = 3
	fmt.Println(Z.Indentf("-----\nindented by %v here", Z.IndentBy))
	Z.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExampleInWrapf() {
	in := Z.IndentBy
	col := Z.Columns
	Z.IndentBy = 3
	Z.Columns = 10
	fmt.Println(
		Z.InWrapf("-----\nindented by %v here and wrapped at %v",
			Z.IndentBy, Z.Columns,
		))
	Z.IndentBy = in
	Z.Columns = col
	// -----
	//    indented
	//    by 3
	//    here
	//    and
	//    wrapped
	//    at 10
}

func ExampleMarkf() {

	in := `
			Must have *%v* block before verbatim:

			     Now we can start
			     a Verbatim
			     block.

			     Which can have blank lines, even.

			And back to a %v block.

			* **foo**
			* bar

			And a numbered list

			1. Something
			2. here

			That's really it.

			`

	fmt.Println("----------------------")
	fmt.Print(Z.Markf(in, "another", "paragraph"))
	fmt.Println("----------------------")

	//Output:
	// ----------------------
	//        Must have <italic>another<reset> block before verbatim:
	//
	//            Now we can start
	//            a Verbatim
	//            block.
	//
	//            Which can have blank lines, even.
	//
	//        And back to a paragraph block.
	//
	//        * <bold>foo<reset>
	//        * bar
	//
	//        And a numbered list
	//
	//        1. Something
	//        2. here
	//
	//        That's really it.
	//
	// ----------------------

}

// ------------------------- Print variations -------------------------

func ExamplePrintEmph_basics() {

	// Emph observes the rwxrob/term escapes
	// (see package documentation for more)

	term.Italic = `<italic>`
	term.Bold = `<bold>`
	term.BoldItalic = `<bolditalic>`
	term.Under = `<under>`
	term.Reset = `<reset>`

	Z.PrintEmph("*ITALIC*\n")
	Z.PrintEmph("**BOLD**\n")
	Z.PrintEmph("***BOLDITALIC***\n")
	Z.PrintEmph("<UNDER>\n") // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExamplePrintWrap() {
	col := Z.Columns
	Z.Columns = 10
	Z.PrintWrap(`some thing here that is more than 10 characters`)
	Z.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExamplePrintInWrap() {
	defer func() { Z.IndentBy = Z.IndentBy }()
	indent := Z.IndentBy
	col := Z.Columns
	Z.Columns = 10
	Z.IndentBy = 4
	fmt.Println("-----")
	Z.PrintInWrap("some\nthat is \n  indented")
	Z.IndentBy = indent
	Z.Columns = col
	// Output:
	// -----
	//     some
	//     that
	//     is
	//     indented
}

func ExamplePrintMark() {

	in := `
			Must have *another* block before verbatim:

			     Now we can start
			     a Verbatim
			     block.

			     Which can have blank lines, even.

			And back to a paragraph block.

			* **foo**
			* bar

			And a numbered list

			1. Something
			2. here

			That's really it.

			`

	fmt.Println("----------------------")
	Z.PrintMark(in)
	fmt.Println("----------------------")

	//Output:
	// ----------------------
	//        Must have <italic>another<reset> block before verbatim:
	//
	//            Now we can start
	//            a Verbatim
	//            block.
	//
	//            Which can have blank lines, even.
	//
	//        And back to a paragraph block.
	//
	//        * <bold>foo<reset>
	//        * bar
	//
	//        And a numbered list
	//
	//        1. Something
	//        2. here
	//
	//        That's really it.
	//
	// ----------------------

}

// --------------------- Print(Sprintf) variations --------------------

func ExamplePrintEmphf() {
	Z.PrintEmphf(`some *%v* thing`, "italic")
	// Output:
	// some <italic>italic<reset> thing
}

func ExamplePrintfWrapf() {
	col := Z.Columns
	Z.Columns = 3
	Z.PrintWrapf(`some %v here`, 10)
	Z.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExamplePrintIndentf() {
	in := Z.IndentBy
	Z.IndentBy = 3
	Z.PrintIndentf("-----\nindented by %v here", Z.IndentBy)
	Z.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExamplePrintInWrapf() {
	in := Z.IndentBy
	col := Z.Columns
	Z.IndentBy = 3
	Z.Columns = 10
	Z.PrintInWrapf("-----\nindented by %v here and wrapped at %v",
		Z.IndentBy, Z.Columns,
	)
	Z.IndentBy = in
	Z.Columns = col
	// -----
	//    indented
	//    by 3
	//    here
	//    and
	//    wrapped
	//    at 10
}

func ExamplePrintMarkf() {

	in := `
			Must have *%v* block before verbatim:

			     Now we can start
			     a Verbatim
			     block.

			     Which can have blank lines, even.

			And back to a %v block.

			* **foo**
			* bar

			And a numbered list

			1. Something
			2. here

			That's really it.

			`

	fmt.Println("----------------------")
	Z.PrintMarkf(in, "another", "paragraph")
	fmt.Println("----------------------")

	//Output:
	// ----------------------
	//        Must have <italic>another<reset> block before verbatim:
	//
	//            Now we can start
	//            a Verbatim
	//            block.
	//
	//            Which can have blank lines, even.
	//
	//        And back to a paragraph block.
	//
	//        * <bold>foo<reset>
	//        * bar
	//
	//        And a numbered list
	//
	//        1. Something
	//        2. here
	//
	//        That's really it.
	//
	// ----------------------

}
