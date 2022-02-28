package term

import (
	"os"

	"github.com/rwxrob/bonzai/term/esc"
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

// IsTerminal returns true if the output is to an interactive terminal
// (not piped in any way). This is useful when detemining if an extra
// line return is needed to avoid making programs chomp the line returns
// unnecessarily.
func IsTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

var Reset = esc.Reset
var Bold = esc.Bold
var Faint = esc.Faint
var Italic = esc.Italic
var Underline = esc.Underline
var BlinkSlow = esc.BlinkSlow
var BlinkFast = esc.BlinkFast
var Reverse = esc.Reverse
var Concealed = esc.Concealed
var StrikeOut = esc.StrikeOut
var BoldItalic = esc.BoldItalic
var Black = esc.Black
var Red = esc.Red
var Green = esc.Green
var Yellow = esc.Yellow
var Blue = esc.Blue
var Magenta = esc.Magenta
var Cyan = esc.Cyan
var White = esc.White
var BBlack = esc.BBlack
var BRed = esc.BRed
var BGreen = esc.BGreen
var BYellow = esc.BYellow
var BBlue = esc.BBlue
var BMagenta = esc.BMagenta
var BCyan = esc.BCyan
var BWhite = esc.BWhite
var HBlack = esc.HBlack
var HRed = esc.HRed
var HGreen = esc.HGreen
var HYellow = esc.HYellow
var HBlue = esc.HBlue
var HMagenta = esc.HMagenta
var HCyan = esc.HCyan
var HWhite = esc.HWhite
var BHBlack = esc.BHBlack
var BHRed = esc.BHRed
var BHGreen = esc.BHGreen
var BHYellow = esc.BHYellow
var BHBlue = esc.BHBlue
var BHMagenta = esc.BHMagenta
var BHCyan = esc.BHCyan
var BHWhite = esc.BHWhite
var ClearScreen = esc.ClearScreen
var X = esc.X
var B = esc.B
var I = esc.I
var U = esc.U
var BI = esc.BI
var CS = esc.CS
