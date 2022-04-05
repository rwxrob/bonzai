package Z

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	"github.com/rwxrob/scan"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
)

// IndentBy is the number of spaces to indent in Indent. Default is 7.
// Bonzai command tree creator can change this for every composite
// command imported their application in this one place.
var IndentBy = 7

// Columns is the number of bytes (not runes) at which Wrap will wrap.
// Default is 80. Bonzai command tree creator can change this for every
// composite command imported their application in this one place.
var Columns = 80

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
// Note that because of the nature of Verbatim's block's initial (4
// space) token Verbatim blocks must never be first since the entire
// input buffer is first dedented and the spaces would grouped with the
// indentation to be stripped. This is never a problem, however,
// because Verbatim blocks never make sense as the first block in
// a BonzaiMark document. This simplicity and clarity of 4-space tokens
// far outweighs the advantages of alternatives (such as fences).
func Blocks(in string) []*Block {

	var blocks []*Block
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
						blocks = append(blocks, &Block{Bulleted, s.Buf[m:s.Pos]})
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
						blocks = append(blocks, &Block{Numbered, s.Buf[m:s.Pos]})
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

			var block []byte
			for s.Scan() {

				if s.Rune == '\n' {

					// add in indented lines
					if s.Is(string(pre)) {
						block = append(block, '\n')
						s.Pos += len(pre)
						continue
					}

					// end of the block
					blocks = append(blocks, &Block{Verbatim, block})
					continue MAIN
				}

				block = append(block, []byte(string(s.Rune))...)
			}

		case '\n', '\r', '\t': // inconsequential white space
			continue

		default: // paragraph
			var block []byte
			block = append(block, []byte(string(s.Rune))...)
			for s.Scan() {
				switch s.Rune {
				case '\n', '\r':
					block = append(block, ' ')
				default:
					block = append(block, []byte(string(s.Rune))...)
				}
				if s.Is("\n\n") {
					blocks = append(blocks, &Block{Paragraph, block})
					s.Scan()
					s.Scan()
					continue MAIN
				}
			}

		}

	}
	return blocks
}

// Mark renders BonzaiMark emphasis spans specifically for
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

// Mark parses the input as a string of BonzaiMark, multiple blocks with
// optional emphasis (see Blocks and Emph).
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
			out += InWrap(Emph(string(block.V))) + "\n"
		case Bulleted:
			out += Indent(Emph(string(block.V))) + "\n"
		case Numbered:
			out += Indent(Emph(string(block.V))) + "\n"
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
