package mark

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/rwxrob/bonzai/term/esc"
)

// OpenItalic opens italic emphasis. Default: ANSI italic.
var OpenItalic = esc.Italic

// CloseItalic closes italic emphasis. Default: ANSI reset.
var CloseItalic = esc.Reset

// OpenBold opens bold emphasis. Default: ANSI bold.
var OpenBold = esc.Bold

// CloseBold closes bold emphasis. Default: ANSI reset.
var CloseBold = esc.Reset

// OpenBoldItalic opens bold italic emphasis. Default: ANSI bold italic.
var OpenBoldItalic = esc.BoldItalic

// CloseBoldItalic closes bold italic emphasis. Default: ANSI reset.
var CloseBoldItalic = esc.Reset

// OpenUnderline open underline emphasis. Default: ANSI underline.
var OpenUnderline = esc.Underline

// CloseUnderline closes underline emphasis. Default: ANSI reset.
var CloseUnderline = esc.Reset

/*
func init() {
	if !term.IsTerminal() {
		reset = ""
		italic = ""
		bold = ""
		bolditalic = ""
		underline = ""
		return
	}
	emphFromLess()
}
*/

func emphFromLess() {
	var x string
	x = os.Getenv("LESS_TERMCAP_us")
	if x != "" {
		OpenItalic = x
	}
	x = os.Getenv("LESS_TERMCAP_md")
	if x != "" {
		OpenBold = x
	}
	x = os.Getenv("LESS_TERMCAP_mb")
	if x != "" {
		OpenBoldItalic = x
	}
	x = os.Getenv("LESS_TERMCAP_us")
	if x != "" {
		OpenUnderline = x
	}
}

// FormatWrapped takes a command documentation format string (an
// extremely limited version of Markdown that is also Godoc friendly)
// and transforms it as follows:
//
// * Initial and trailing blank lines are removed.
//
// * Indentation is removed - the number of spaces preceeding the first
//   word of the first line are ignored in every line (including raw
//   text blocks).
//
// * Raw text ignored - any line beginning with four or more spaces
//   (after convenience indentation is removed) will be kept as it is
//   exactly (code examples, etc.) but should never exceed 80 characters
//   (including the spaces).
//
// * Blocks are unwrapped - any non-blank (without three or less initial
//   spaces) will be trimmed line following a line will be joined to the
//   preceding line recursively (unless hard break).
//
// * Hard breaks kept - like Markdown any line that ends with two or
//   more spaces will automatically force a line return.
//
// * URL links argument names and anything else within angle brackets
//   (<url>), will trigger underline in both text blocks
//   and usage sections.
//
// * Italic, Bold, and BoldItalic inline emphasis using one, two, or
//   three stars respectively will be observed and cannot be intermixed
//   or intra-word.  Each opener must be preceded by a UNICODE space (or
//   nothing) and followed by a non-space rune. Each closer must be
//   preceded by a non-space rune and followed by a UNICODE space (or
//   nothing).
//
// For historic reasons the following environment variables will be
// observed if found (and also provide color support for the less pager
// utility):
//
//   * Italic      LESS_TERMCAP_so
//   * Bold        LESS_TERMCAP_md
//   * BoldItalic  LESS_TERMCAP_mb
//   * Underline   LESS_TERMCAP_us
//
func FormatWrapped(input string, indent, width int) (output string) {

	// this scanner could be waaaay more lexy
	// but suits the need and clear to read

	var strip int
	var blockbuf string

	// standard state machine approach
	inblock := false
	inraw := false
	inhard := false
	gotindent := false

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		txt := scanner.Text()
		trimmed := strings.TrimSpace(txt)

		// ignore blank lines
		if !(inraw || inblock) && len(trimmed) == 0 {
			continue
		}

		// infer the indent to strip for every line
		if !gotindent && len(trimmed) > 0 {
			for i, v := range txt {
				if v != ' ' {
					strip = i
					gotindent = true
					break
				}
			}
		}

		// strip convenience indent
		if len(txt) >= strip {
			txt = txt[strip:]
		}

		// raw block start
		if !inblock && !inraw && len(txt) > 4 && txt[0:4] == "    " && len(trimmed) > 0 {
			inraw = true
			output += "\n\n" + txt
			continue
		}

		// in raw block
		if inraw && len(txt) > 4 {
			output += "\n" + txt
			continue
		}

		// raw block end
		if inraw && len(trimmed) == 0 {
			inraw = false
			continue
		}

		// another block line, join it
		if inblock && len(trimmed) > 0 {
			if len(txt) >= 2 && txt[len(txt)-2:] == "  " {
				inhard = true
			}
			space := " "
			if inhard {
				space = "\n"
			}
			blockbuf += space + trimmed
			continue
		}

		// beginning of a new block
		if !inblock && len(trimmed) > 0 {
			inhard = false
			inblock = true
			if len(txt) >= 2 && txt[len(txt)-2:] == "  " {
				inhard = true
			}
			blockbuf = trimmed
			continue
		}

		// end block
		if inblock && len(trimmed) == 0 {
			inblock = false
			output += "\n\n" + Format(Wrap(blockbuf, width-strip-4))
			continue
		}
	}

	// flush last block
	if inblock {
		output += "\n\n" + Format(Wrap(blockbuf, width-strip-4))
	}
	output = Indent(strings.TrimSpace(output), indent)
	return
}

// Format replaces minimal Markdown-like syntax with *Italic*,
// **Bold**, ***BoldItalic***, and <bracketed>
func Format(buf string) string {
	// TODO add functional parser, for fun
	log.Println("would Format")
	return buf
}

// Indent indents each line the set number of spaces.
func Indent(buf string, spaces int) string {
	nbuf := ""
	scanner := bufio.NewScanner(strings.NewReader(buf))
	scanner.Scan()
	for n := 0; n < spaces; n++ {
		nbuf += " "
	}
	nbuf += scanner.Text()
	for scanner.Scan() {
		nbuf += "\n"
		for n := 0; n < spaces; n++ {
			nbuf += " "
		}
		nbuf += scanner.Text()
	}
	return nbuf
}

// Wrapped is the same as Format but without any emphasis.
func Wrapped(input string, indent, width int) (output string) {

	// this scanner could be waaaay more lexy
	// but suits the need and clear to read

	var strip int
	var blockbuf string

	// standard state machine approach
	inblock := false
	inraw := false
	inhard := false
	gotindent := false

	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		txt := scanner.Text()
		trimmed := strings.TrimSpace(txt)

		// ignore blank lines
		if !(inraw || inblock) && len(trimmed) == 0 {
			continue
		}

		// infer the indent to strip for every line
		if !gotindent && len(trimmed) > 0 {
			for i, v := range txt {
				if v != ' ' {
					strip = i
					gotindent = true
					break
				}
			}
		}

		// strip convenience indent
		if len(txt) >= strip {
			txt = txt[strip:]
		}

		// raw block start
		if !inblock && !inraw && len(txt) > 4 && txt[0:4] == "    " && len(trimmed) > 0 {
			inraw = true
			output += "\n\n" + txt
			continue
		}

		// in raw block
		if inraw && len(txt) > 4 {
			output += "\n" + txt
			continue
		}

		// raw block end
		if inraw && len(trimmed) == 0 {
			inraw = false
			continue
		}

		// another block line, join it
		if inblock && len(trimmed) > 0 {
			if len(txt) >= 2 && txt[len(txt)-2:] == "  " {
				inhard = true
			}
			space := " "
			if inhard {
				space = "\n"
			}
			blockbuf += space + trimmed
			continue
		}

		// beginning of a new block
		if !inblock && len(trimmed) > 0 {
			inhard = false
			inblock = true
			if len(txt) >= 2 && txt[len(txt)-2:] == "  " {
				inhard = true
			}
			blockbuf = trimmed
			continue
		}

		// end block
		if inblock && len(trimmed) == 0 {
			inblock = false
			output += "\n\n" + Wrap(blockbuf, width-strip-4)
			continue
		}
	}

	// flush last block
	if inblock {
		output += "\n\n" + Wrap(blockbuf, width-strip-4)
	}
	output = Indent(strings.TrimSpace(output), indent)
	return
}

// peekWord returns the runes up to the next space.
func peekWord(buf []rune, start int) []rune {
	word := []rune{}
	for _, r := range buf[start:] {
		if unicode.IsSpace(r) {
			break
		}
		word = append(word, r)
	}
	return word
}

// Wrap wraps the string to the given width using spaces to separate
// words. If passed a negative width will effectively join all words in
// the buffer into a single line with no wrapping.
func Wrap(buf string, width int) string {
	if width == 0 {
		return buf
	}
	nbuf := ""
	curwidth := 0
	for i, r := range []rune(buf) {
		// hard breaks always as is
		if r == '\n' {
			nbuf += "\n"
			curwidth = 0
			continue
		}
		if unicode.IsSpace(r) {
			// FIXME: don't peek every word, only after passed width
			// change the space to a '\n' in the buffer slice directly
			next := peekWord([]rune(buf), i+1)
			if width > 0 && (curwidth+len(next)+1) > width {
				nbuf += "\n"
				curwidth = 0
				continue
			}
		}
		nbuf += string(r)
		curwidth++
	}
	return nbuf
}
