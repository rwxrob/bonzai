package term

import (
	"os"

	"github.com/rwxrob/bonzai/term/esc"
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
	Hidden     string
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

// IsInteractive returns true if the output is to an interactive terminal
// (not piped in any way). This is useful when determining if an extra
// line return is needed to avoid making programs chomp the line returns
// unnecessarily.
func IsInteractive() bool {
	if f, _ := os.Stdout.Stat(); (f.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

// MustVT100 calls bonzai.ExitError if terminal does not support minimal
// VT100 terminal standard.
func MustVT100() {
	// TODO
}

func init() {
	if IsInteractive() {
		AttrOn()
	}
}

// AttrOff sets all the terminal attributes to zero values (empty strings).
// Note that this does not affect anything in the esc subpackage (which
// contains the constants from the VT100 specification).
func AttrOff() {
	Reset = ""
	Bright = ""
	Bold = ""
	Dim = ""
	Italic = ""
	Under = ""
	Blink = ""
	BlinkF = ""
	Reverse = ""
	Hidden = ""
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
// contains the constants from the VT100 specification).
func AttrOn() {
	Reset = esc.Reset
	Bright = esc.Bright
	Bold = esc.Bold
	Dim = esc.Dim
	Italic = esc.Italic
	Under = esc.Under
	Blink = esc.Blink
	BlinkF = esc.BlinkF
	Reverse = esc.Reverse
	Hidden = esc.Hidden
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
