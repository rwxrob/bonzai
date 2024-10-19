// Copyright 2022 Robert S. Muhlestein
// SPDX-License-Identifier: Apache-2.0

//go:build !aix && !js && !nacl && !plan9 && !windows && !android && !solaris

package term

import (
	"syscall"
	"unsafe"
)

// WinSizeUpdate makes a SYS_IOCTL syscall to get the information for
// the current terminal. This returns nothing but zeros unless the
// terminal is interactive (standard output is a terminal). Consider
// gdamore/tcell if more reliable dimensions are needed.
func WinSizeUpdate() {
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&WinSize)))
}
