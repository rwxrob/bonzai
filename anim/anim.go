package anim

import (
	"fmt"
	"syscall"

	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/term/esc"
)

// SimpleAnimationScreen conveniently sets up an alternate screen buffer
// clears it, turns off the cursor and traps any interrupts so that the
// screen and cursor are restored and the program exits. This is useful
// when making simple ASCII animations without needing a full terminal
// animation package.
func SimpleAnimationScreen() error {
	run.Trap(func() {
		fmt.Print(esc.Clear)     // Clear terminal
		fmt.Print(esc.AltBufOff) // Show cursor
		fmt.Print(esc.CursorOn)  // Show cursor
		run.Exit()
	}, syscall.SIGINT, syscall.SIGTERM)
	fmt.Print(esc.CursorOff)
	fmt.Print(esc.AltBufOn)
	return nil
}
