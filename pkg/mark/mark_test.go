package mark_test

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/pkg/cmd"
	"github.com/rwxrob/bonzai/pkg/term"
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
	//fmt.Println(blocks)
	fmt.Printf("%q\n", blocks[0])
	fmt.Printf("%q\n", blocks[1])

	//Output:
	// "* some thing\n* another thing"
	// "* another block\n* here"

}

func ExampleBlocks_numbered() {
	in := `

			1. some thing
			2. another thing

			1. another block
			2. here

			`

	blocks := Z.Blocks(in)
	fmt.Printf("%q\n", blocks[0])
	fmt.Printf("%q\n", blocks[1])

	//Output:
	// "1. some thing\n2. another thing"
	// "1. another block\n2. here"

}

func ExampleBlocks_paragraph() {

	in := `
			Simple paragraph
			here on multiple
			lines.

			And another   one here
			with just a bit more.

			`

	blocks := Z.Blocks(in)
	//fmt.Println(len(blocks))
	//fmt.Printf("%v", blocks)
	fmt.Printf("%q\n", blocks[0])
	fmt.Printf("%q\n", blocks[1])

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
			* Must have another block type first.

			     Now we can start
			     a Verbatim
			     block.
			     
			     Which can have blank lines, even.

			* And back to another bulleted list.

			`

	blocks := Z.Blocks(in)
	//fmt.Println(len(blocks))
	//fmt.Println(blocks)
	fmt.Printf("%q\n", blocks[0])
	fmt.Printf("%q\n", blocks[1])
	fmt.Printf("%q\n", blocks[2])

	//Output:
	// "* Must have another block type first."
	// "Now we can start\na Verbatim\nblock.\n\nWhich can have blank lines, even."
	// "* And back to another bulleted list."

}

// -------------------------- main BonzaiMark -------------------------

func ExampleEmph_under() {
	term.Under = `<under>`
	term.Reset = `<reset>`
	fmt.Println(Z.Emph("<UNDER>"))
	fmt.Println(Z.Emph("< UNDER >"))
	// Output:
	// <<under>UNDER<reset>>
	// < UNDER >
}

func ExampleEmph_boldItalic() {
	term.BoldItalic = `<bolditalic>`
	term.Reset = `<reset>`
	fmt.Println(Z.Emph("***BoldItalic***"))
	// Output:
	// <bolditalic>BoldItalic<reset>
}

func ExampleEmph_bold() {
	term.Bold = `<bold>`
	term.Reset = `<reset>`
	fmt.Println(Z.Emph("**Bold**"))
	fmt.Println(Z.Emph("** Bold **"))
	// Output:
	// <bold>Bold<reset>
	// ** Bold **
}

func ExampleEmph_italic() {
	term.Italic = `<italic>`
	term.Reset = `<reset>`
	fmt.Println(Z.Emph("*Italic*"))
	fmt.Println(Z.Emph("* Italic *"))
	// Output:
	// <italic>Italic<reset>
	// * Italic *
}

func ExampleEmph_code() {
	term.Under = `<code>`
	term.Reset = `<reset>`
	fmt.Println(Z.Emph("`Code`"))
	fmt.Println(Z.Emph("` Code `"))
	fmt.Println(Z.Emph("`.git`"))
	// Output:
	// <code>Code<reset>
	// ` Code `
	// <code>.git<reset>
}

func ExampleEmph_basics() {

	// Emph observes the term escapes
	// (see package documentation for more)

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

	// Emph observes the term escapes
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

/*
func ExampleWrap_again() {

	defer func() {
		term.Italic = `<italic>`
		term.Bold = `<bold>`
		term.BoldItalic = `<bolditalic>`
		term.Under = `<under>`
		term.Reset = `<reset>`
	}()

	term.Italic = esc.Italic
	term.Bold = esc.Bold
	term.BoldItalic = esc.BoldItalic
	term.Under = esc.Under
	term.Reset = esc.Reset
	Z.IndentBy = 0
	Z.Columns = 40

	cmd := &Z.Cmd{
		Name: `some`,
		MarkMap: template.FuncMap{
			"builddir":  func() string { return "a/build/dir" },
			"buildfile": func() string { return "build.yaml" },
		},

		Description: `
		   The {{cmd .Name}} command looks for a {{pre buildfile}} file in
		   the specified or current directory and runs the build command on
		   each building them all concurrently into the {{pre builddir}}
		   directory where they are ready for upload to GitHub as a release.
		   If an argument is passed it is expected to be an explicit path to
		   a different build directory. If no path is specified will assume
		   the current directory and recursively search all parents for
		   {{pre buildfile}} until found creating a {{pre builddir}} within that
		   directory for the newly built artifacts.  `,
	}

	fmt.Println("Output")
	//	fmt.Println(Z.Mark(cmd.Fill(cmd.Description)))
	fmt.Println(to.Wrapped(cmd.Fill(cmd.Description), Z.Columns))

	// Output:
	// some

}
*/
