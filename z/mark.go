package Z

import (
	"log"
	"regexp"
	"unicode"

	"github.com/rwxrob/scan"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// This file contains BonzaiMark.

// IndentBy is the number of spaces to indent in Indent. Default is 7.
// Bonzai command tree creator can change this for every composite
// command imported their application in this one place.
var IndentBy = 7

// Columns is the number of bytes (not runes) at which Wrap will wrap.
// Default is 80. Bonzai command tree creator can change this for every
// composite command imported their application in this one place.
var Columns = 80

// Emph renders BonzaiMark emphasis spans specifically for
// VT100-compatible terminals (which almost all are today):
//
//     *Italic*
//     **Bold**
//     ***BoldItalic***
//     <under> (keeping brackets)
//
// See Mark for block formatting and rwxrob/term for terminal rendering.
func Emph(buf string) string {
	var nbuf []rune
	var opentok, closetok bool
	var otok, ctok string
	prev := ' '

	for i := 0; i < len([]rune(buf)); i++ {
		r := []rune(buf)[i]

		if r == '<' {
			nbuf = append(nbuf, '<')
			nbuf = append(nbuf, []rune(term.Under)...)
			for {
				i++
				r = rune(buf[i])
				if r == '>' {
					i++
					break
				}
				nbuf = append(nbuf, r)
			}
			nbuf = append(nbuf, []rune(term.Reset)...)
			nbuf = append(nbuf, '>')
			i--
			continue
		}

		if r != '*' {

			if opentok {
				tokval := " "
				if !unicode.IsSpace(r) {
					switch otok {
					case "*":
						tokval = term.Italic
					case "**":
						tokval = term.Bold
					case "***":
						tokval = term.BoldItalic
					}
				} else {
					tokval = otok
				}
				nbuf = append(nbuf, []rune(tokval)...)
				opentok = false
				otok = ""
			}

			if closetok {
				nbuf = append(nbuf, []rune(term.Reset)...) // practical, not perfect
				ctok = ""
				closetok = false
			}

			prev = r
			nbuf = append(nbuf, r)
			continue
		}

		// everything else for '*'
		if unicode.IsSpace(prev) || opentok {
			opentok = true
			otok += string(r)
			continue
		}

		// only closer conditions remain
		if !unicode.IsSpace(prev) {
			closetok = true
			ctok += string(r)
			continue
		}

		// nothing special
		closetok = false
		nbuf = append(nbuf, r)
	}

	// for tokens at the end of a block
	if closetok {
		nbuf = append(nbuf, []rune(term.Reset)...)
	}

	return string(nbuf)
}

// Wrap wraps to Columns width.
func Wrap(in string) string { w, _ := to.Wrapped(in, Columns); return w }

// Indent indents the number of spaces set by IndentBy.
func Indent(in string) string { return to.Indented(in, IndentBy) }

// InWrap combines both Wrap and Indent.
func InWrap(in string) string {
	w, _ := to.Wrapped(in, Columns-IndentBy)
	return to.Indented(w, IndentBy)
}

// ---------------------------- finish this ---------------------------

// TODO finish the full Mark implementation and documentation

/*
Mark renders BonzaiMark markup as formatted VT100-compatible terminal
output, wrapped and indented as specified by the package variables
IndentBy and Columns.

BonzaiMark

BonzaiMark is a minimal subset of CommonMark (Markdown) suitable for
rendering to terminals, web pages, PDF, or plain text. The syntax is
deliberately simple and easy to read (much like standard GoDoc) but rich
enough to look well rendered in documents as well as the terminal.

Stripped Indentation

Unlike CommonMark, BonzaiMark ignores any number of blank lines or
whitespace before the first line and uses the initial characters
preceding the first line (tabs or spaces) as a basis for what to strip
from every subsequent line of the document. All trailing white space and
blank lines are also discarded. This allows BonzaiMark to be included in
source code in very readable ways (preferably with backtick string
literals).

    some := `
		    Here is a paragraph
				that will have the initial spaces
				stripped.
		`

Structure: Blocks of Spans

Like CommonMark, every BonzaiMark document consists of one or more
blocks which contain one or more spans of the following type:

    plain
    *italic*
    **bold**
    ***bolditalic***
    <under> (brackets remain)

Unlike CommonMark, spans may not contain any other span type.

While the underline format is not supported in CommonMark, underlining is frequently used in place of italic for most terminals. Angle brackets are, however, supported in CommonMark.

The following limited block types are supported. (All others from CommonMark are not):

    * Paragraph
		* Verbatim
		* Numbered List
		* Bulleted List
		* Numbered Long List Item
		* Bulleted Long List Item

Unlike CommonMark, all blocks must be separated by two or more line returns.

Paragraph Blocks

Paragraph blocks are the most common. They consist of one or more of spans.

Verbatim Blocks

A Verbatim block will be included exactly as typed. It begins with the
first line that has four or more spaces followed by a non-whitespace
character. The block continues until the next block is detected. (All
other blocks must begin on the first column (after stripped
indentation).

Lists

Like CommonMark there are long lists and compact lists. Long lists are
actually multiple consecutive long list item blocks whereas compact
lists consist of list items that are on one line after another (no
double-return block separator).

Lists are either numbered or bulleted and never exceed one level (no
nesting). Both list types may contain any number of paragraph blocks but
most will only contain the one. When multiple paragraph blocks are
wanted the subsequent paragraph block must begin on the same column as
the first character of the first line of the first list item paragraph;
they must line up.

    1. **Keep list item paragraphs lined up**

		   This is a second paragraph block under the same list item because
			 it lines up with the first line of the first paragraph

blank lines -- over multiple consecutive lines but each line after the
first must line up exactly with the first character of the first list
item line to be considered still a part of the list item.


Numbered Lists

Numbered lists always begin with a go integer and a dot (.).
Conventionally a 1. is used for everything so that document maintainers
can quickly reorganize when needed without a tool for renumbered. Nested
lists are not supported. Each item in the list must be on immediate
subsequent lines. Numbered lists must always be rendered with Arabic
numerals.

Bulleted Lists

Bulleted lists must begin with a single asterisk (*) followed by
a single space. No other bullet type from CommonMark is supported.
Nested lists are not supported. Each item in a list must be on an
immediately subsequent line. List items may contain any number of spans
over multiple consecutive lines but each line after the first must line
up exactly with the first character of the first list item line.

    * This is a
		  list item

Only Inline Links

Only explicit link URLs are supported. The must always be wrapped with
angle brackets (<>). Technically inline links are a span of type "under"
which also gives them an underline emphasis on the terminal.

No Escapes

There is no support for escaping anything in BonzaiMark. (CommonMark allows the placement of a backslash to remove any special meaning.) Therefore, most authors will use verbatim blocks when it is necessary to use the reserved BonzaiMark tokens in other ways.



Soft and Hard Line Endings

Like CommonMark lines that follow other lines immediately are
effectively joined together unless there are two or more spaces at the
end of the line (a hard return). This is after any indentation has been
removed (see Stripped Indentation).



of spaces for the first line of indentation.
Any line beginning with at least four spaces (after trimming
indentation) will be kept verbatim.

Emphasis will be applied as possible if the following markup is
detected:

Note that the format of the emphasis might not always be as
specifically named. For example, most terminal do not support italic
fonts and so will instead underline *italic* text, so (as specified
in HTML5 for <i>, for example) these format names should be taken to
mean their semantic equivalents.

For terminal rendering details see the rwxrob/term package.
*/

// Mark

func Mark(in string) string {
	if in == "" {
		return ""
	}

	//var out string
	blocks := Blocks(in)
	log.Print(blocks)

	//out := to.Dedented(markup)
	//out, _ = to.Wrapped(out, 80)
	//out = Emph(out)
	//return out
	return ""
}

// Blocks strips preceding and trailing white space and then checks the
// first line for indentation (spaces or tabs) and strips that exact
// indentation string from every line. It then breaks up the input into
// blocks separated by one or more empty lines and applies basic
// formatting to each as follows:
//
//     If is one of the following leave alone with no wrapping:
//
//     * Bulleted List - beginning with *
//     * Numbered List - beginning with 1.
//     * Verbatim      - beginning with four spaces
//
//     Everything else is considered a "paragraph" and will be unwrapped
//     into a single long line (which is normally wrapped later).
//
// For now, these blocks are added as is, but plans are to eventually
// add support for short and long lists much like CommonMark.
//
// Note that because of the nature of Verbatim's block's initial (4
// space) token Verbatim blocks must never be first since the entire
// input buffer is first dedented and the spaces would grouped with the
// indentation to be stripped. This is never a problem, however,
// because Verbatim blocks never make sense as the first block in
// a BonzaiMark document. This simplicity and clarity of 4-space tokens
// far outweighs the advantages of alternatives (such as fences).
func Blocks(in string) []string {

	var blocks []string
	verbpre := regexp.MustCompile(` {4,}`)
	s := scan.R{Buf: []byte(to.Dedented(in))}

MAIN:
	for s.Scan() {

		switch s.Rune {

		case '*': // bulleted list
			if s.Is(" ") {
				m := s.Pos - 1
				for s.Scan() {
					if s.Is("\n\n") {
						blocks = append(blocks, string(s.Buf[m:s.Pos]))
						s.Pos += 2
						continue MAIN
					}
				}
			}

		case '1': // numbered list
			if s.Is(". ") {
				m := s.Pos - 1
				for s.Scan() {
					if s.Is("\n\n") {
						blocks = append(blocks, string(s.Buf[m:s.Pos]))
						s.Pos += 2
						continue MAIN
					}
				}
			}

		case ' ': // verbatim
			s.Pos -= 1
			ln := s.Match(verbpre)
			s.Pos++

			if ln < 0 {
				continue
			}
			pre := s.Buf[s.Pos-1 : s.Pos+ln-1]
			s.Pos += len(pre) - 1

			block := []rune{}
			for s.Scan() {

				if s.Rune == '\n' {

					// add in indented lines
					if s.Is(string(pre)) {
						block = append(block, '\n')
						s.Pos += len(pre)
						continue
					}

					// end of the block
					blocks = append(blocks, string(block))
					continue MAIN
				}

				block = append(block, s.Rune)
			}

		case '\n', '\r', '\t': // inconsequential white space
			continue

		default: // paragraph
			block := []rune{s.Rune}
			for s.Scan() {
				switch s.Rune {
				case '\n', '\r':
					block = append(block, ' ')
				default:
					block = append(block, s.Rune)
				}
				if s.Is("\n\n") {
					blocks = append(blocks, string(block))
					s.Scan()
					s.Scan()
					continue MAIN
				}
			}

		}

	}
	return blocks
}
