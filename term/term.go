package term

import (
	"os"
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
