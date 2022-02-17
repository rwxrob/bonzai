package bonzai

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Method defines the main code to execute for a command (Cmd). By
// convention the parameter list should be names "args" if there are
// args expected and "none" if not. Methods must never write error
// output to anything but standard error and should almost always use
// the log package to do so.
type Method func(args ...string) error

// ----------------------- errors, exit, debug -----------------------

// DEBUG is set when os.Getenv("CMDBOX_DEBUG") is set to anything.
// Produces verbose debugging logs to stderr to help cmdbox users
// develop robust tools.
//
var DEBUG bool

func init() {
	if os.Getenv("CMDBOX_DEBUG") != "" {
		DEBUG = true
	}
}

// DoNotExit is a utility for disabling any call to os.Exit via any
// caller in order to trap panics and use them for testing, etc.
//
var DoNotExit bool

// ExitOff sets DoNotExit to false.
func ExitOff() { DoNotExit = true }

// ExitOn sets DoNotExit to true.
func ExitOn() { DoNotExit = false }

// Exit calls os.Exit(0) unless DoNotExit has been set to true.
func Exit() {
	if !DoNotExit {
		os.Exit(0)
	}
}

// ExitError prints err and exits with 1 return value unless DoNotExit
// has been set to true.
func ExitError(err ...interface{}) {
	switch e := err[0].(type) {
	case string:
		if len(e) > 1 {
			log.Printf(e+"\n", err[1:])
		}
		log.Println(e)
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

// ArgsFrom returns a list of strings split on space with an extra
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
		args = append(args, " ")
	}
	return args
}
