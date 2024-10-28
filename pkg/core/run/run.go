// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rwxrob/bonzai/pkg/core/futil"
	"github.com/rwxrob/bonzai/pkg/core/mark"
)

// ExeName returns just the base name of the executable from
// [os.Executable] without the path or any suffix (ex: .exe). Note that
// the name may actually be a symbolic link. Use [RealExeName] if the
// resolved version is wanted.
func ExeName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return path, err
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return name, nil
}

// RealExeName returns same as [ExeName] but with all symbolic links
// resolved.
func RealExeName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return path, err
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return name, nil
}

func addExeName(base string) (string, error) {
	name, err := ExeName()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, name), nil
}

func addRealExeName(base string) (string, error) {
	name, err := RealExeName()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, name), nil
}

func ExeCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return addExeName(dir)
}

func RealExeCacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return addRealExeName(dir)
}

func ExeConfigDir() (string, error) {
	dir, err := futil.UserConfigDir()
	if err != nil {
		return "", err
	}
	return addExeName(dir)
}

func RealExeStateDir() (string, error) {
	dir, err := futil.UserStateDir()
	if err != nil {
		return "", err
	}
	return addExeName(dir)
}

func ExeStateDir() (string, error) {
	dir, err := futil.UserStateDir()
	if err != nil {
		return "", err
	}
	return addExeName(dir)
}

func ExeIsSymLink() (bool, error) {
	path, err := os.Executable()
	if err != nil {
		return false, err
	}
	return futil.IsSymLink(path)
}

func ExeIsHardLink() (bool, error) {
	path, err := os.Executable()
	if err != nil {
		return false, err
	}
	return futil.IsHardLink(path)
}

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
// system and then runs it with [exec.Command.Run]  exiting in a way that
// is supported across all architectures that Go supports. The [os.Stdin],
// [os.Stdout], and [os.Stderr] are connected directly to that of the calling
// program. Sometimes this is insufficient and the UNIX-specific SysExec
// is preferred.
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

func ShellIsBash() bool {
	return strings.Contains(os.Getenv("SHELL"), `bash`)
}

func ShellIsFish() bool {
	return len(os.Getenv("FISH_VERSION")) > 0
}

func ShellIsZsh() bool {
	return strings.Contains(os.Getenv("SHELL"), `zsh`)
}

func ShellIsPowerShell() bool {
	return len(os.Getenv(`PSModulePath`)) > 0
}

// ArgsFrom returns a list of field strings split on space (using
// [strings.Fields]) with an extra trailing special space item appended
// if the line has any trailing spaces at all signifying a definite word
// boundary and not a potential prefix. This is critical for resolving
// completion with [bonzai.Completer].
func ArgsFrom(line string) []string {
	args := []string{}
	if line == "" {
		return args
	}
	args = strings.Fields(line)
	if line[len(line)-1] == ' ' {
		args = append(args, "")
	}
	return args
}

// ArgsOrIn takes an slice or nil as argument and if the slice has any
// length greater than 0 returns all the argument joined together with
// a single space between them. Otherwise, will read standard input
// until end of file reached (Cntl-D). Returns empty string if error.
func ArgsOrIn(args []string) string {
	if len(args) == 0 {
		buf, err := io.ReadAll(os.Stdin)
		if err != nil {
			return ""
		}
		return string(buf)
	}
	return strings.Join(args, " ")
}

// AllowPanic disables TrapPanic stopping it from cleaning panic errors.
var AllowPanic = false

// TrapPanic recovers from any panic and more gracefully displays the
// panic by logging it before exiting with a return value of 1.
var TrapPanic = func() {
	if !AllowPanic {
		if r := recover(); r != nil {
			log.Println(r)
			os.Exit(1)
		}
	}
}

// ExitError prints err and exits with 1 return value unless DoNotExit
// has been set to true. Commands should usually never call ExitError
// themselves returning an error from their Method instead.
func ExitError(err ...interface{}) {
	prev := os.Stdout
	os.Stdout = os.Stderr
	previ := mark.IndentBy
	mark.IndentBy = 0

	switch e := err[0].(type) {

	case string:
		if len(e) > 1 {
			mark.Printf(e+"\n", err[1:]...)
		} else {
			mark.Print(e)
		}

	case error:
		out := fmt.Sprintf("%v", e)
		if len(out) > 0 {
			fmt.Println(strings.TrimSpace(mark.Sprint(out)))
		}
	}

	mark.IndentBy = previ
	os.Stdout = prev

	if !DoNotExit {
		os.Exit(1)
	}

}

// Exit calls os.Exit(0) unless DoNotExit has been set to true. Cmds
// should never call Exit themselves returning a nil error from their
// Methods instead.
func Exit() {
	if !DoNotExit {
		os.Exit(0)
	}
}

// DoNotExit effectively disables Exit and ExitError allowing the
// program to continue running, usually for test evaluation.
var DoNotExit bool

// ExitOff sets DoNotExit to false.
func ExitOff() { DoNotExit = true }

// ExitOn sets DoNotExit to true.
func ExitOn() { DoNotExit = false }
