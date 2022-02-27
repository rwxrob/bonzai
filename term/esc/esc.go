/*
Package esc contains commonly used ANSI terminal escape sequences for
different attributes and terminal control. These are meant to be printed
directly to any terminal supporting the curses standard. They should be
used directly whenever possible for best terminal performance.  Also see
the equivalent variables in the term package that are set to empty strings when support for them is not detected at init() time in order to maintain
the original stateful approach to terminal escape sequence support (first designed for the terminfo C library.)
*/
package esc

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Faint     = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	BlinkSlow = "\033[5m"
	BlinkFast = "\033[6m"
	Reverse   = "\033[7m"
	Concealed = "\033[8m"
	StrikeOut = "\033[9m"

	BoldItalic = "\033[1m\033[3m"

	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	BBlack   = "\033[40m"
	BRed     = "\033[41m"
	BGreen   = "\033[42m"
	BYellow  = "\033[43m"
	BBlue    = "\033[44m"
	BMagenta = "\033[45m"
	BCyan    = "\033[46m"
	BWhite   = "\033[47m"

	HBlack   = "\033[90m"
	HRed     = "\033[91m"
	HGreen   = "\033[92m"
	HYellow  = "\033[93m"
	HBlue    = "\033[94m"
	HMagenta = "\033[95m"
	HCyan    = "\033[96m"
	HWhite   = "\033[97m"

	BHBlack   = "\033[100m"
	BHRed     = "\033[101m"
	BHGreen   = "\033[102m"
	BHYellow  = "\033[103m"
	BHBlue    = "\033[104m"
	BHMagenta = "\033[105m"
	BHCyan    = "\033[106m"
	BHWhite   = "\033[107m"

	ClearScreen = "\033[2J\033[H"

	X  = "\033[0m"        // Reset
	B  = "\033[1m"        // Bold
	I  = "\033[3m"        // Italic
	U  = "\033[4m"        // Underline
	BI = "\033[1m\033[3m" // BoldItalic
	CS = "\033[2J\033[H"  // ClearScreen
)
