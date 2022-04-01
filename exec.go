// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// SysExec will check for the existence of the first argument as an
// executable on the system and then execute it using syscall.Exec(),
// which replaces the currently running program with the new one in all
// respects (stdin, stdout, stderr, process ID, signal handling, etc).
// Generally speaking, this is only available on UNIX variations.  This
// is exceptionally faster and cleaner than calling any of the os/exec
// variations, but it can make your code far be less compatible
// with different operating systems.
func SysExec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	// exits the program unless there is an error
	return syscall.Exec(path, args, os.Environ())
}

// Exec checks for existence of first argument as an executable on the
// system and then runs it without exiting in a way that is supported
// across all architectures that Go supports. The stdin, stdout, and stderr are
// connected directly to that of the calling program. Sometimes this is
// insufficient and the UNIX-specific SysExec is preferred. For example,
// when handing over control to a terminal editor such as Vim.
func Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
