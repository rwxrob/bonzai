// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

/*
Package esc contains commonly used VT100 ANSI terminal escape sequences
for different attributes and terminal control. These are meant to be
printed directly to any terminal supporting the VT100 standard. They
should be used directly whenever possible for best terminal performance.
Also see the equivalent variables in the term package that are set to
empty strings when support for them is not detected at init() time in
order to maintain the original stateful approach to terminal escape
sequence support (first designed for the terminfo C library.)
*/
package esc

const (
	QCode   = "\033[c"  // "\033[{code}0c"
	QStatus = "\033[5n" // "\033[0n" ok "\033[3n" not ok
	QPos    = "\033[6n" // "\033[{row};{col}R"

	ResetDevice = "\033c" // reset to defaults
	LineWrapOn  = "\033[7h"
	LineWrapOff = "\033[7l"

	DefFont = "\033("
	AltFont = "\033)"

	TopLeft = "\033[H" // "\033[{row};{col}H "\033[{row};{col}f"
	Up      = "\033[A" // "\033[{count}A"
	Down    = "\033[B" // "\033[{count}B"
	Forward = "\033[C" // "\033[{count}C"
	Back    = "\033[D" // "\033[{count}D"

	Save      = "\033[s" // save position
	Unsave    = "\033[u" // restore position
	SaveAll   = "\0337"  // save position and attributes
	UnsaveAll = "\0338"  // restore position and attributes

	ScrollOn   = "\033[r" // "\033[{start}:{end}r"
	ScrollDown = "\033D"
	ScrollUp   = "\033M"

	SetTab      = "\033H"
	ClearTab    = "\033[g"
	ClearAllTab = "\033[3g"

	EraseLineE   = "\033[K"
	EraseLineS   = "\033[1K"
	EraseLine    = "\033[2K"
	EraseScreenE = "\033[J"
	EraseScreenS = "\033[1J"
	EraseScreen  = "\033[2J"
	Clear        = EraseScreen + TopLeft + Reset

	// unsupported on most modern terminals
	PrintScreen = "\033[i"
	PrintLine   = "\033[1i"
	PrintLogOn  = "\033[4i"
	PrintLogOff = "\033[5i"

	// SetKey = "\033[{key};"string"p"

	// rest are attributes, the most common

	Reset      = "\033[0m"
	Bright     = "\033[1m"
	Bold       = Bright
	Dim        = "\033[2m"
	Italic     = "\033[3m" // usually Reverses instead
	Under      = "\033[4m"
	Blink      = "\033[5m"
	BlinkF     = "\033[6m" // usually not supported
	Reverse    = "\033[7m"
	Hide       = "\033[8m"
	Strike     = "\033[9m" // modern support
	BoldItalic = Bold + Italic

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

	CursorOn  = "\033[?25h"
	CursorOff = "\033[?25l"

	AltBufOn  = "\033[?1049h"
	AltBufOff = "\033[?1049l"
)
