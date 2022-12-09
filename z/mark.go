package Z

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	"github.com/rwxrob/pegn/scanner"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// IndentBy is the number of spaces to indent in Indent. Default is 7.
// Bonzai command tree creator can change this for every composite
// command imported their application in this one place.
var IndentBy = 7

// Columns is the number of bytes (not runes) at which Wrap will wrap.
// By default detects the terminal width (if possible) otherwise keeps
// 80 standard. Bonzai command tree creator can change this for every
// composite command imported their application in this one place.
var Columns = int(term.WinSize.Col)

// Lines returns the string converted into a slice of lines.
func Lines(in string) []string { return to.Lines(in) }

const (
	Paragraph = iota + 1
	Numbered
	Bulleted
	Verbatim
)

type Block struct {
	T int
	V []byte
}

// String fulfills the fmt.Stringer interface.
func (s *Block) String() string { return string(s.V) }

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
// If no blocks are parsed returns an empty slice of Block pointers ([]
// *Block).
//
// Note that because of the nature of Verbatim's block's initial (4
// space) token Verbatim blocks must never be first since the entire
// input buffer is first dedented and the spaces would grouped with the
// indentation to be stripped. This is never a problem, however,
// because Verbatim blocks never make sense as the first block in
// a BonzaiMark document. This simplicity and clarity of 4-space tokens
// far outweighs the advantages of alternatives (such as fences).
//
//  PEGN Specification
//
//        Grammar     <-- Block*
//        Block       <-- Bulleted / Numbered / Verbatim / Paragraph
//        Bulleted    <-- '* ' (!EOB unipoint)* EOB
//        Numbered    <-- '1. ' (!EOB unipoint)* EOB
//        Verbatim    <-- '    ' (!EOB unipoint)* EOB
//        Paragraph   <-- (!EOB unipoint)* EOB
//        EOB          <- LF{2} / EOD
//        EOD          <- # end of data stream
//
func Blocks(in string) []*Block {
	var blocks []*Block

	in = to.Dedented(in) // also trims initial white space

	s := scanner.New(in)
	//s.Trace++

	for s.Scan() {

		// Bulleted
		if s.Is("* ") {
			var beg, end int
			beg = s.E - 1
			for s.Scan() {
				if s.Is("\n\n") {
					end = s.E - 1
					s.E++
					break
				}
				end = s.E
			}
			blocks = append(blocks, &Block{Bulleted, s.Buf[beg:end]})
			continue
		}

		// Numbered
		if s.Is("1. ") {
			var beg, end int
			beg = s.E - 1
			for s.Scan() {
				if s.Is("\n\n") {
					end = s.E - 1
					s.E++
					break
				}
				end = s.E
			}
			blocks = append(blocks, &Block{Numbered, s.Buf[beg:end]})
			continue
		}

		// Verbatim
		if ln := s.Match(begVerbatim); ln >= 4 {
			var beg, end int
			beg = s.B
			for s.Scan() {
				if s.Is("\n\n") {
					end = s.E - 1
					s.E++
					break
				}
				end = s.E
			}
			dedented := to.Dedented(string(s.Buf[beg:end]))
			blocks = append(blocks, &Block{Verbatim, []byte(dedented)})
			continue
		}

		// Paragraph (default)
		if !unicode.IsSpace(s.R) {
			var beg, end int
			beg = s.B
			for s.Scan() {
				if s.Is("\n\n") {
					end = s.E - 1
					s.E++
					break
				}
				end = s.E
			}
			blocks = append(blocks, &Block{Paragraph,
				[]byte(to.Words(string(s.Buf[beg:end])))})
			continue
		}

	}

	return blocks
}

// don't expose these until mark has own package

var begVerbatim = regexp.MustCompile(`^ {4,}`)
var ws = regexp.MustCompile(`^[\s\r\n]+`)
var begUnder = regexp.MustCompile(`^<\p{L}`)
var endUnder = regexp.MustCompile(`^\p{L}>`)
var begBoldItalic = regexp.MustCompile(`^\*{3}\p{L}`)
var endBoldItalic = regexp.MustCompile(`^\p{L}\*{3}`)
var begBold = regexp.MustCompile(`^\*{2}\p{L}`)
var endBold = regexp.MustCompile(`^\p{L}\*{2}`)
var begItalic = regexp.MustCompile(`^\*\p{L}`)
var endItalic = regexp.MustCompile(`^\p{L}\*`)
var begCode = regexp.MustCompile(`^` + "`" + `\S`)
var endCode = regexp.MustCompile(`^\S` + "`")

// Emph renders BonzaiMark emphasis spans specifically for
// VT100-compatible terminals (which almost all are today):
//
//     *Italic*
//     **Bold**
//     ***BoldItalic***
//     <Under> (keeping brackets)
//     `Code`
//
// See Mark for block formatting and term for terminal rendering.
func Emph[T string | []byte | []rune](buf T) string {
	var nbuf []rune

	s := scanner.New(buf)

	for s.Scan() {

		// <under>
		if s.Match(begUnder) > 0 {
			nbuf = append(nbuf, '<')
			nbuf = append(nbuf, []rune(term.Under)...)
			for s.Scan() {
				if s.Match(endUnder) > 0 {
					nbuf = append(nbuf, s.R)
					nbuf = append(nbuf, []rune(term.Reset)...)
					nbuf = append(nbuf, '>')
					s.E++
					break
				}
				nbuf = append(nbuf, s.R)
			}
			continue
		}

		// ***BoldItalic***
		if s.Match(begBoldItalic) > 0 {
			s.Scan()
			s.Scan()
			nbuf = append(nbuf, []rune(term.BoldItalic)...)
			for s.Scan() {
				if s.Match(endBoldItalic) > 0 {
					nbuf = append(nbuf, s.R)
					nbuf = append(nbuf, []rune(term.Reset)...)
					s.E += 3
					break
				}
				nbuf = append(nbuf, s.R)
			}
			continue
		}

		// **Bold**
		if s.Match(begBold) > 0 {

			s.E += 1
			nbuf = append(nbuf, []rune(term.Bold)...)
			for s.Scan() {
				if s.Match(endBold) > 0 {
					nbuf = append(nbuf, s.R)
					s.E += 2
					nbuf = append(nbuf, []rune(term.Reset)...)
					break
				}
				nbuf = append(nbuf, s.R)
			}
			continue
		}

		// *Italic*
		if s.Match(begItalic) > 0 {
			nbuf = append(nbuf, []rune(term.Italic)...)
			for s.Scan() {
				if s.Match(endItalic) > 0 {
					nbuf = append(nbuf, s.R)
					nbuf = append(nbuf, []rune(term.Reset)...)
					s.E++
					break
				}
				nbuf = append(nbuf, s.R)
			}
			continue
		}

		// `Code`
		if s.Match(begCode) > 0 {
			nbuf = append(nbuf, []rune(term.Under)...)
			for s.Scan() {
				if s.Match(endCode) > 0 {
					nbuf = append(nbuf, s.R)
					nbuf = append(nbuf, []rune(term.Reset)...)
					s.E++
					break
				}
				nbuf = append(nbuf, s.R)
			}
			continue
		}

		nbuf = append(nbuf, s.R)

	} // end main scan loop

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

// Mark parses the input as a string of BonzaiMark, multiple blocks with
// optional emphasis (see Blocks and Emph) and applies IndentBy and
// Columns wrapping to it.
func Mark(in string) string {

	if in == "" {
		return ""
	}

	blocks := Blocks(in)
	if len(blocks) == 0 {
		return ""
	}

	var out string

	for _, block := range blocks {
		switch block.T {
		case Paragraph:
			out += Emph(InWrap(string(block.V))) + "\n"
		case Bulleted:
			out += Emph(Indent(string(block.V))) + "\n"
		case Numbered:
			out += Emph(Indent(string(block.V))) + "\n"
		case Verbatim:
			out += to.Indented(Indent(string(block.V)), 4) + "\n"
		default:
			panic("unknown block type: " + strconv.Itoa(block.T))
		}

	}

	return out
}

// Emphf calls fmt.Sprintf on the string before passing it to Emph.
func Emphf(a string, f ...any) string {
	return Emph(fmt.Sprintf(a, f...))
}

// Indentf calls fmt.Sprintf on the string before passing it to Indent.
func Indentf(a string, f ...any) string {
	return Indent(fmt.Sprintf(a, f...))
}

// Wrapf calls fmt.Sprintf on the string before passing it to Wrap.
func Wrapf(a string, f ...any) string {
	return Wrap(fmt.Sprintf(a, f...))
}

// InWrapf calls fmt.Sprintf on the string before passing it to InWrap.
func InWrapf(a string, f ...any) string {
	return InWrap(fmt.Sprintf(a, f...))
}

// Markf calls fmt.Sprintf on the string before passing it to Mark.
func Markf(a string, f ...any) string {
	return Mark(fmt.Sprintf(a, f...))
}

// PrintEmph passes string to Emph and prints it.
func PrintEmph(a string) { fmt.Print(Emph(a)) }

// PrintWrap passes string to Wrap and prints it.
func PrintWrap(a string) { fmt.Print(Wrap(a)) }

// PrintIndent passes string to Indent and prints it.
func PrintIndent(a string) { fmt.Print(Indent(a)) }

// PrintInWrap passes string to InWrap and prints it.
func PrintInWrap(a string) { fmt.Print(InWrap(a)) }

// PrintMark passes string to Mark and prints it.
func PrintMark(a string) { fmt.Print(Mark(a)) }

// PrintEmphf calls fmt.Sprintf on the string before passing it to Emph
// and then printing it.
func PrintEmphf(a string, f ...any) {
	fmt.Print(Emph(fmt.Sprintf(a, f...)))
}

// PrintWrapf calls fmt.Sprintf on the string before passing it to Wrap
// and then printing it.
func PrintWrapf(a string, f ...any) {
	fmt.Print(Wrap(fmt.Sprintf(a, f...)))
}

// PrintIndentf calls fmt.Sprintf on the string before passing it to
// Indent and then printing it.
func PrintIndentf(a string, f ...any) {
	fmt.Print(Indent(fmt.Sprintf(a, f...)))
}

// PrintInWrapf calls fmt.Sprintf on the string before passing it to
// InWrap and then printing it.
func PrintInWrapf(a string, f ...any) {
	fmt.Print(InWrap(fmt.Sprintf(a, f...)))
}

// PrintMarkf calls fmt.Sprintf on the string before passing it to Mark
// and then printing it.
func PrintMarkf(a string, f ...any) {
	fmt.Print(Mark(fmt.Sprintf(a, f...)))
}
