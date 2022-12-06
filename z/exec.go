// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z

import (
	"fmt"
	"log"
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
// system and then runs it with exec.Command.Run  exiting in a way that
// is supported across all architectures that Go supports. The stdin,
// stdout, and stderr are connected directly to that of the calling
// program. Sometimes this is insufficient and the UNIX-specific SysExec
// is preferred. See exec.Command.Run for more about distinguishing
// ExitErrors from others.
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

// Out returns the standard output of the executed command as
// a string. Errors are logged but not returned.
func Out(args ...string) string {
	if len(args) == 0 {
		log.Println("missing name of executable")
		return ""
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Println(err)
		return ""
	}
	out, err := exec.Command(path, args[1:]...).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
