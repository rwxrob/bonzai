// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Exec (not to be confused with Execute) will check for the existence
// of the first argument as an executable on the system and then execute
// it using syscall.Exec(), which replaces the currently running program
// with the new one in all respects (stdin, stdout, stderr, process ID,
// signal handling, etc).
//
// Note that although this is exceptionally faster and cleaner than
// calling any of the os/exec variations it may be less compatible with
// different operating systems.
func Exec(args ...string) error {
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

// Run checks for existence of first argument as an executable on the
// system and then runs it without exiting in a way that is supported
// across different operating systems. The stdin, stdout, and stderr are
// connected directly to that of the calling program. Use more specific
// exec alternatives if intercepting stdout and stderr are desired. Also
// see Exec.
func Run(args ...string) error {
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
