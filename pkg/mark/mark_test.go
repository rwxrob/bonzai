package mark_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/pkg/mark"
	"github.com/rwxrob/bonzai/pkg/term"
	"github.com/rwxrob/bonzai/pkg/to"
)

func init() {
	term.Italic = `<italic>`
	term.Bold = `<bold>`
	term.BoldItalic = `<bolditalic>`
	term.Under = `<under>`
	term.Reset = `<reset>`
}

func ExampleLines() {
	fmt.Printf("%q\n", to.Lines("line one\nline two"))
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

	blocks := mark.Blocks(in)
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

	blocks := mark.Blocks(in)
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

	blocks := mark.Blocks(in)
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

	blocks := mark.Blocks(in)
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
	fmt.Println(mark.Emph("<UNDER>"))
	fmt.Println(mark.Emph("< UNDER >"))
	// Output:
	// <<under>UNDER<reset>>
	// < UNDER >
}

func ExampleEmph_boldItalic() {
	term.BoldItalic = `<bolditalic>`
	term.Reset = `<reset>`
	fmt.Println(mark.Emph("***BoldItalic***"))
	// Output:
	// <bolditalic>BoldItalic<reset>
}

func ExampleEmph_bold() {
	term.Bold = `<bold>`
	term.Reset = `<reset>`
	fmt.Println(mark.Emph("**Bold**"))
	fmt.Println(mark.Emph("** Bold **"))
	// Output:
	// <bold>Bold<reset>
	// ** Bold **
}

func ExampleEmph_italic() {
	term.Italic = `<italic>`
	term.Reset = `<reset>`
	fmt.Println(mark.Emph("*Italic*"))
	fmt.Println(mark.Emph("* Italic *"))
	// Output:
	// <italic>Italic<reset>
	// * Italic *
}

func ExampleEmph_code() {
	term.Under = `<code>`
	term.Reset = `<reset>`
	fmt.Println(mark.Emph("`Code`"))
	fmt.Println(mark.Emph("` Code `"))
	fmt.Println(mark.Emph("`.git`"))
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

	fmt.Println(mark.Emph("*ITALIC*"))
	fmt.Println(mark.Emph("**BOLD**"))
	fmt.Println(mark.Emph("***BOLDITALIC***"))
	fmt.Println(mark.Emph("<UNDER>")) // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExampleWrap() {
	col := mark.Columns
	mark.Columns = 10
	fmt.Println(mark.Wrap(`some thing here that is more than 10 characters`))
	mark.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExampleIndent() {
	indent := mark.IndentBy
	col := mark.Columns
	mark.Columns = 10
	mark.Columns = 10
	mark.IndentBy = 4
	fmt.Printf("%q", mark.Indent("some\nthat is \n  indented"))
	mark.IndentBy = indent
	mark.Columns = col
	// Output:
	// "    some\n    that is \n      indented\n"
}

func ExampleInWrap() {
	defer func() { mark.IndentBy = mark.IndentBy }()
	indent := mark.IndentBy
	col := mark.Columns
	mark.Columns = 10
	mark.IndentBy = 4
	fmt.Printf("%q", mark.InWrap("some\nthat is \n  indented"))
	mark.IndentBy = indent
	mark.Columns = col
	// Output:
	// "    some\n    that\n    is\n    indented\n"
}

func ExampleSprint_simple() {

	fmt.Print(mark.Sprint(`**foo**`))

	//Output:
	// <bold>foo<reset>

}

func ExampleSprint() {

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
	fmt.Print(mark.Sprint(in))
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
	fmt.Println(mark.Emphf(`some *%v* thing`, "italic"))
	// Output:
	// some <italic>italic<reset> thing
}

func ExampleWrapf() {
	col := mark.Columns
	mark.Columns = 3
	fmt.Println(mark.Wrapf(`some %v here`, 10))
	mark.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExampleIndentf() {
	in := mark.IndentBy
	mark.IndentBy = 3
	fmt.Println(mark.Indentf("-----\nindented by %v here", mark.IndentBy))
	mark.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExampleInWrapf() {
	in := mark.IndentBy
	col := mark.Columns
	mark.IndentBy = 3
	mark.Columns = 10
	fmt.Println(
		mark.InWrapf("-----\nindented by %v here and wrapped at %v",
			mark.IndentBy, mark.Columns,
		))
	mark.IndentBy = in
	mark.Columns = col
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
	fmt.Print(mark.Sprintf(in, "another", "paragraph"))
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

	mark.PrintEmph("*ITALIC*\n")
	mark.PrintEmph("**BOLD**\n")
	mark.PrintEmph("***BOLDITALIC***\n")
	mark.PrintEmph("<UNDER>\n") // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExamplePrintWrap() {
	col := mark.Columns
	mark.Columns = 10
	mark.PrintWrap(`some thing here that is more than 10 characters`)
	mark.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExamplePrintInWrap() {
	defer func() { mark.IndentBy = mark.IndentBy }()
	indent := mark.IndentBy
	col := mark.Columns
	mark.Columns = 10
	mark.IndentBy = 4
	fmt.Println("-----")
	mark.PrintInWrap("some\nthat is \n  indented")
	mark.IndentBy = indent
	mark.Columns = col
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
	mark.Print(in)
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
	mark.PrintEmphf(`some *%v* thing`, "italic")
	// Output:
	// some <italic>italic<reset> thing
}

func ExamplePrintfWrapf() {
	col := mark.Columns
	mark.Columns = 3
	mark.PrintWrapf(`some %v here`, 10)
	mark.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExamplePrintIndentf() {
	in := mark.IndentBy
	mark.IndentBy = 3
	mark.PrintIndentf("-----\nindented by %v here", mark.IndentBy)
	mark.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExamplePrintInWrapf() {
	in := mark.IndentBy
	col := mark.Columns
	mark.IndentBy = 3
	mark.Columns = 10
	mark.PrintInWrapf("-----\nindented by %v here and wrapped at %v",
		mark.IndentBy, mark.Columns,
	)
	mark.IndentBy = in
	mark.Columns = col
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
	mark.Printf(in, "another", "paragraph")
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
	mark.IndentBy = 0
	mark.Columns = 40

	cmd := &mark.Cmd{
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
	//	fmt.Println(mark.Mark(cmd.Fill(cmd.Description)))
	fmt.Println(to.Wrapped(cmd.Fill(cmd.Description), mark.Columns))

	// Output:
	// some

}
*/
