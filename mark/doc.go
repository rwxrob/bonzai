/*
# Specification

	Grammar   <- Block*
	Block     <- Header / Bulleted / Numbered / Verbatim / Paragraph
	Header    <- '#'{1,6} SP (!EOB !LF rune)* EOB
	Bulleted  <- '* ' (!EOB rune)* EOB
	Numbered  <- '1. ' (!EOB rune)* EOB
	Verbatim  <- SP{4} (!EOB rune)* EOB
	Paragraph <- (!EOB rune)* EOB
	EOB       <- LF{2} / EOD
	EOD       <- # end of data stream

Blocks strips preceding and trailing white space and then checks the
first line for indentation (spaces or tabs) and strips that exact
indentation string from every line. It then breaks up the input into
blocks separated by one or more empty lines and applies basic
formatting to each as follows:

  - Bulleted List - beginning with *
  - Numbered List - beginning with 1.
  - Verbatim      - beginning with four spaces

Everything else is considered a "paragraph" and will be unwrapped
into a single long line (which is normally wrapped later).

If no blocks are parsed returns an empty slice of Block pointers.

Note that because of the nature of Verbatim's block's initial (4
space) token Verbatim blocks must never be first since the entire
input buffer is first dedented and the spaces would be grouped with the
indentation to be stripped. This is never a problem, however,
because Verbatim blocks never make sense as the first block in
a BonzaiMark document.
*/
package mark
