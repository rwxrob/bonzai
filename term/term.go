// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

/*
Package term provides traditional terminal escapes and interactions
including an Prompter/Responder specification for Expect-like,
line-based APIs, communications, and user interfaces.
*/
package term

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"unicode"

	"github.com/rwxrob/bonzai/term/esc"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	Reset      string
	Bright     string
	Bold       string
	Dim        string
	Italic     string
	Under      string
	Blink      string
	BlinkF     string
	Reverse    string
	Hide       string
	Strike     string
	BoldItalic string
	Black      string
	Red        string
	Green      string
	Yellow     string
	Blue       string
	Magenta    string
	Cyan       string
	White      string
	BBlack     string
	BRed       string
	BGreen     string
	BYellow    string
	BBlue      string
	BMagenta   string
	BCyan      string
	BWhite     string
	HBlack     string
	HRed       string
	HGreen     string
	HYellow    string
	HBlue      string
	HMagenta   string
	HCyan      string
	HWhite     string
	BHBlack    string
	BHRed      string
	BHGreen    string
	BHYellow   string
	BHBlue     string
	BHMagenta  string
	BHCyan     string
	BHWhite    string
	X          string
	B          string
	I          string
	U          string
	BI         string
)

// WinSizeStruct is the exact struct used by the ioctl system library.
type WinSizeStruct struct {
	Row, Col       uint16
	Xpixel, Ypixel uint16
}

// WinSize is 80x24 by default but is detected and set to a more
// accurate value at init() time on systems that support ioctl
// (currently) and can be updated with WinSizeUpdate on systems that
// support it. This value can be overriden by those wishing a more
// consistent value or who prefer not to fill the screen completely when
// displaying help and usage information.
var WinSize WinSizeStruct

var interactive bool

func init() {
	SetInteractive(DetectInteractive())
	EmphFromLess()
	WinSizeUpdate()
}

// SetInteractive forces the interactive internal state affecting output
// including calling AttrOn (true) or AttrOff (false).
func SetInteractive(to bool) {
	interactive = to
	if to {
		AttrOn()
	} else {
		AttrOff()
	}
}

// IsInteractive returns the internal interactive state set by
// SetInteractive. The default is that returned by DetectInteractive set
// at  init() time.
func IsInteractive() bool { return interactive }

// DetectInteractive returns true if the output is to an interactive
// terminal (not piped in any way).
func DetectInteractive() bool {
	if f, _ := os.Stdout.Stat(); (f.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

var attron bool

// IsAttrOn contains the state of the last AttrOn/AttrOff call.
func IsAttrOn() bool { return attron }

// AttrOff sets all the terminal attributes to zero values (empty strings).
// Note that this does not affect anything in the esc subpackage (which
// contains the constants from the VT100 specification). Sets the
// AttrAreOn bool to false.
func AttrOff() {
	attron = false
	Reset = ""
	Bright = ""
	Bold = ""
	Dim = ""
	Italic = ""
	Under = ""
	Blink = ""
	BlinkF = ""
	Reverse = ""
	Hide = ""
	Strike = ""
	BoldItalic = ""
	Black = ""
	Red = ""
	Green = ""
	Yellow = ""
	Blue = ""
	Magenta = ""
	Cyan = ""
	White = ""
	BBlack = ""
	BRed = ""
	BGreen = ""
	BYellow = ""
	BBlue = ""
	BMagenta = ""
	BCyan = ""
	BWhite = ""
	HBlack = ""
	HRed = ""
	HGreen = ""
	HYellow = ""
	HBlue = ""
	HMagenta = ""
	HCyan = ""
	HWhite = ""
	BHBlack = ""
	BHRed = ""
	BHGreen = ""
	BHYellow = ""
	BHBlue = ""
	BHMagenta = ""
	BHCyan = ""
	BHWhite = ""
	X = ""
	B = ""
	I = ""
	U = ""
	BI = ""
}

// AttrOn sets all the terminal attributes to zero values (empty strings).
// Note that this does not affect anything in the esc subpackage (which
// contains the constants from the VT100 specification). Sets the
// AttrAreOn bool to true.
func AttrOn() {
	attron = true
	Reset = esc.Reset
	Bright = esc.Bright
	Bold = esc.Bold
	Dim = esc.Dim
	Italic = esc.Italic
	Under = esc.Under
	Blink = esc.Blink
	BlinkF = esc.BlinkF
	Reverse = esc.Reverse
	Hide = esc.Hide
	Strike = esc.Strike
	Black = esc.Black
	Red = esc.Red
	Green = esc.Green
	Yellow = esc.Yellow
	Blue = esc.Blue
	Magenta = esc.Magenta
	Cyan = esc.Cyan
	White = esc.White
	BBlack = esc.BBlack
	BRed = esc.BRed
	BGreen = esc.BGreen
	BYellow = esc.BYellow
	BBlue = esc.BBlue
	BMagenta = esc.BMagenta
	BCyan = esc.BCyan
	BWhite = esc.BWhite
	HBlack = esc.HBlack
	HRed = esc.HRed
	HGreen = esc.HGreen
	HYellow = esc.HYellow
	HBlue = esc.HBlue
	HMagenta = esc.HMagenta
	HCyan = esc.HCyan
	HWhite = esc.HWhite
	BHBlack = esc.BHBlack
	BHRed = esc.BHRed
	BHGreen = esc.BHGreen
	BHYellow = esc.BHYellow
	BHBlue = esc.BHBlue
	BHMagenta = esc.BHMagenta
	BHCyan = esc.BHCyan
	BHWhite = esc.BHWhite
	X = esc.Reset
	B = esc.Bold
	I = esc.Italic
	U = esc.Under
	BI = esc.BoldItalic
}

// Print calls Println if IsInteractive, otherwise, Print. This works
// better with applications that would otherwise need to trim the
// tailing line return to set that as values for shell variables.
func Print(a ...any) (int, error) {
	if IsInteractive() {
		return fmt.Println(a...)
	}
	return fmt.Print(a...)
}

// Printf calls fmt.Printf directly but adds a line return if
// IsInteractive.
func Printf(format string, a ...any) (int, error) {
	if IsInteractive() {
		format += "\n"
	}
	return fmt.Printf(format, a...)
}

// Read reads a single line of input and chomps the \r?\n. Also see
// ReadHide.
func Read() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// TrapPanic recovers from any panic and more gracefully displays the
// panic by logging it before exiting. See log package for ways to alter
// the output of the embedded log.Println function.
var TrapPanic = func() {
	if r := recover(); r != nil {
		log.Println(r)
		os.Exit(1)
	}
}

// InOutFunc is a function that takes input as a string and simply
// responds based on that input.
type InOutFunc func(in string) string

// REPL starts a rudimentary, functional, read-evaluate print loop
// passing each line of prompted input to the respond function as it is
// entered, printing the response, and then prompting
// for another. No tab-completion is supported or planned. In this way,
// a REPL can be used to connect prompt and respond functions by passing
// input/output back and forth.
//
// The output from the prompt function will be printed directly to the
// terminal before prompting for more input. Most prompt function
// implementations will print a preceding line return.  Prompt function
// implementations must state explicitly any terminal requirements
// (plain text, markup, VT100 compatible, etc.)
//
// Response functions can be useful for encapsulating bots and other
// intelligent responders to any terminal input. In theory, one respond
// function can be connected to another. In that sense, respond
// functions are a rudimentary, single-line replacement for other API
// interactions (such as rest).
//
// Either prompt or respond functions may use panic or os.Exit to end
// the program, but panic is generally preferred since the REPL (or
// other caller) can trap it and exit gracefully. Panics within a REPL
// are generally sent directly to the user and therefore may break the
// all-lowercase convention normally observed for panic messages.
func REPL(prompt, respond InOutFunc) {
	defer TrapPanic()
	var input string
	for {
		p := prompt(input)
		fmt.Print(p)
		input = Read() // should block
		fmt.Print(respond(input))
	}
}

// ReadHide disables the cursor and echoing to the screen and reads
// a single line of input. Leading and trailing whitespace are removed.
// Also see Read.
func ReadHide() string {
	byt, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(byt))
}

// Prompt prints the given message if the terminal IsInteractive and
// reads the string by calling Read. The argument signature is identical
// as that passed to fmt.Printf().
func Prompt(form string, args ...any) string {
	if IsInteractive() {
		fmt.Printf(form, args...)
	}
	return Read()
}

// PromptHide prints the given message if the terminal IsInteractive
// and reads the string by calling ReadHide (which does not echo to
// the screen). The argument signature is identical and passed to to
// fmt.Printf().
func PromptHide(form string, args ...any) string {
	if IsInteractive() {
		fmt.Printf(form, args...)
	}
	return ReadHide()
}

// StripNonPrint remove non-printable runes, e.g. control characters in
// a string that is meant for consumption by terminals that support
// control characters.
func StripNonPrint(s string) string {
	return strings.Map(
		func(r rune) rune {
			if unicode.IsPrint(r) {
				return r
			}
			return -1
		}, s)
}

// EmphFromLess sets Italic, Bold, BoldItalic, and Under from the
// LESS_TERMCAP_us, _md, _mb, and _us environment variables
// respectively. This is a long used way to provide color to UNIX man
// pages dating back to initial color terminals. UNIX users frequently
// set these to provide color to man pages and more. Observes AttrAreOn
// and will simply return if set to false. EmphFromLess is called at
// package init() time automatically.
func EmphFromLess() {
	if !attron {
		return
	}
	var x string
	x = os.Getenv("LESS_TERMCAP_us")
	if x != "" {
		Italic = x
	}
	x = os.Getenv("LESS_TERMCAP_md")
	if x != "" {
		Bold = x
	}
	x = os.Getenv("LESS_TERMCAP_mb")
	if x != "" {
		BoldItalic = x
	}
	x = os.Getenv("LESS_TERMCAP_us")
	if x != "" {
		Under = x
	}
}
