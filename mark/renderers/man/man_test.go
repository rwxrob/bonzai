package man_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/mark/renderers/man"
	"github.com/rwxrob/bonzai/term"
	"github.com/rwxrob/bonzai/to"
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

	blocks := man.Blocks(in)
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

	blocks := man.Blocks(in)
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

	blocks := man.Blocks(in)
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

	blocks := man.Blocks(in)
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
	fmt.Println(man.Emph("<UNDER>"))
	fmt.Println(man.Emph("< UNDER >"))
	// Output:
	// <<under>UNDER<reset>>
	// < UNDER >
}

func ExampleEmph_boldItalic() {
	term.BoldItalic = `<bolditalic>`
	term.Reset = `<reset>`
	fmt.Println(man.Emph("***BoldItalic***"))
	// Output:
	// <bolditalic>BoldItalic<reset>
}

func ExampleEmph_bold() {
	term.Bold = `<bold>`
	term.Reset = `<reset>`
	fmt.Println(man.Emph("**Bold**"))
	fmt.Println(man.Emph("** Bold **"))
	// Output:
	// <bold>Bold<reset>
	// ** Bold **
}

func ExampleEmph_italic() {
	term.Italic = `<italic>`
	term.Reset = `<reset>`
	fmt.Println(man.Emph("*Italic*"))
	fmt.Println(man.Emph("* Italic *"))
	// Output:
	// <italic>Italic<reset>
	// * Italic *
}

func ExampleEmph_code() {
	term.Under = `<code>`
	term.Reset = `<reset>`
	fmt.Println(man.Emph("`Code`"))
	fmt.Println(man.Emph("` Code `"))
	fmt.Println(man.Emph("`.git`"))
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

	fmt.Println(man.Emph("*ITALIC*"))
	fmt.Println(man.Emph("**BOLD**"))
	fmt.Println(man.Emph("***BOLDITALIC***"))
	fmt.Println(man.Emph("<UNDER>")) // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExampleWrap() {
	col := man.Columns
	man.Columns = 10
	fmt.Println(man.Wrap(`some thing here that is more than 10 characters`))
	man.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExampleIndent() {
	indent := man.IndentBy
	col := man.Columns
	man.Columns = 10
	man.Columns = 10
	man.IndentBy = 4
	fmt.Printf("%q", man.Indent("some\nthat is \n  indented"))
	man.IndentBy = indent
	man.Columns = col
	// Output:
	// "    some\n    that is \n      indented\n"
}

func ExampleInWrap() {
	defer func() { man.IndentBy = man.IndentBy }()
	indent := man.IndentBy
	col := man.Columns
	man.Columns = 10
	man.IndentBy = 4
	fmt.Printf("%q", man.InWrap("some\nthat is \n  indented"))
	man.IndentBy = indent
	man.Columns = col
	// Output:
	// "    some\n    that\n    is\n    indented\n"
}

func ExampleSprint_simple() {

	fmt.Print(man.Sprint(`**foo**`))

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
	fmt.Print(man.Sprint(in))
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
	fmt.Println(man.Emphf(`some *%v* thing`, "italic"))
	// Output:
	// some <italic>italic<reset> thing
}

func ExampleWrapf() {
	col := man.Columns
	man.Columns = 3
	fmt.Println(man.Wrapf(`some %v here`, 10))
	man.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExampleIndentf() {
	in := man.IndentBy
	man.IndentBy = 3
	fmt.Println(man.Indentf("-----\nindented by %v here", man.IndentBy))
	man.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExampleInWrapf() {
	in := man.IndentBy
	col := man.Columns
	man.IndentBy = 3
	man.Columns = 10
	fmt.Println(
		man.InWrapf("-----\nindented by %v here and wrapped at %v",
			man.IndentBy, man.Columns,
		))
	man.IndentBy = in
	man.Columns = col
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
	fmt.Print(man.Sprintf(in, "another", "paragraph"))
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

	man.PrintEmph("*ITALIC*\n")
	man.PrintEmph("**BOLD**\n")
	man.PrintEmph("***BOLDITALIC***\n")
	man.PrintEmph("<UNDER>\n") // keeps brackets

	// Output:
	// <italic>ITALIC<reset>
	// <bold>BOLD<reset>
	// <bolditalic>BOLDITALIC<reset>
	// <<under>UNDER<reset>>

}

func ExamplePrintWrap() {
	col := man.Columns
	man.Columns = 10
	man.PrintWrap(`some thing here that is more than 10 characters`)
	man.Columns = col
	// Output:
	// some thing
	// here that
	// is more
	// than 10
	// characters
}

func ExamplePrintInWrap() {
	defer func() { man.IndentBy = man.IndentBy }()
	indent := man.IndentBy
	col := man.Columns
	man.Columns = 10
	man.IndentBy = 4
	fmt.Println("-----")
	man.PrintInWrap("some\nthat is \n  indented")
	man.IndentBy = indent
	man.Columns = col
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
	man.Print(in)
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
	man.PrintEmphf(`some *%v* thing`, "italic")
	// Output:
	// some <italic>italic<reset> thing
}

func ExamplePrintfWrapf() {
	col := man.Columns
	man.Columns = 3
	man.PrintWrapf(`some %v here`, 10)
	man.Columns = col
	// Output:
	// some
	// 10
	// here
}

func ExamplePrintIndentf() {
	in := man.IndentBy
	man.IndentBy = 3
	man.PrintIndentf("-----\nindented by %v here", man.IndentBy)
	man.IndentBy = in
	// Output:
	// -----
	//    indented by 3 here
}

func ExamplePrintInWrapf() {
	in := man.IndentBy
	col := man.Columns
	man.IndentBy = 3
	man.Columns = 10
	man.PrintInWrapf("-----\nindented by %v here and wrapped at %v",
		man.IndentBy, man.Columns,
	)
	man.IndentBy = in
	man.Columns = col
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
	man.Printf(in, "another", "paragraph")
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
	man.IndentBy = 0
	man.Columns = 40

	cmd := &man.Cmd{
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
	//	fmt.Println(man.Mark(cmd.Fill(cmd.Description)))
	fmt.Println(to.Wrapped(cmd.Fill(cmd.Description), man.Columns))

	// Output:
	// some

}
*/
