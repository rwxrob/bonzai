// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Method defines the main code to execute for a command (Cmd). By
// convention the parameter list should be named "args" if there are
// args expected and underscore (_) if not. Methods must never write
// error output to anything but standard error and should almost always
// use the log package to do so.
type Method func(caller *Cmd, args ...string) error

// ----------------------- errors, exit, debug -----------------------

// DoNotExit effectively disables Exit and ExitError allowing the
// program to continue running, usually for test evaluation.
var DoNotExit bool

// ExitOff sets DoNotExit to false.
func ExitOff() { DoNotExit = true }

// ExitOn sets DoNotExit to true.
func ExitOn() { DoNotExit = false }

// Exit calls os.Exit(0) unless DoNotExit has been set to true. Cmds
// should never call Exit themselves returning a nil error from their
// Methods instead.
func Exit() {
	if !DoNotExit {
		os.Exit(0)
	}
}

// ExitError prints err and exits with 1 return value unless DoNotExit
// has been set to true. Commands should usually never call ExitError
// themselves returning an error from their Method instead.
func ExitError(err ...interface{}) {
	switch e := err[0].(type) {
	case string:
		if len(e) > 1 {
			log.Printf(e+"\n", err[1:]...)
		} else {
			log.Println(e)
		}
	case error:
		out := fmt.Sprintf("%v", e)
		if len(out) > 0 {
			log.Println(out)
		}
	}
	if !DoNotExit {
		os.Exit(1)
	}
}

// ArgsFrom returns a list of field strings split on space with an extra
// trailing special space item appended if the line has any trailing
// spaces at all signifying a definite word boundary and not a potential
// prefix.
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
