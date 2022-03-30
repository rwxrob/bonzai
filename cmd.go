// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/conf"
	"github.com/rwxrob/fn/each"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/structs/qstack"
)

type Cmd struct {
	Name        string            `json:"name,omitempty"`
	Aliases     []string          `json:"aliases,omitempty"`
	Summary     string            `json:"summary,omitempty"`
	Usage       string            `json:"usage,omitempty"`
	Version     string            `json:"version,omitempty"`
	Updated     int               `json:"updated,omitempty"` // isosec
	Copyright   string            `json:"copyright,omitempty"`
	License     string            `json:"license,omitempty"`
	Description string            `json:"description,omitempty"`
	Site        string            `json:"site,omitempty"`
	Source      string            `json:"source,omitempty"`
	Issues      string            `json:"issues,omitempty"`
	Other       map[string]string `json:"issues,omitempty"`
	Commands    []*Cmd            `json:"commands,omitempty"`
	Params      []string          `json:"params,omitempty"`
	Hidden      []string          `json:"hide,omitempty"`

	Completer comp.Completer  `json:"-"`
	Conf      conf.Configurer `json:"-"`
	// Cache cache.Cacher `json:"-"`

	Root    *Cmd   `json:"-"`
	Caller  *Cmd   `json:"-"`
	Call    Method `json:"-"`
	MinArgs int    `json:"-"`

	_aliases map[string]*Cmd
}

func (x *Cmd) cacheAliases() {
	x._aliases = map[string]*Cmd{}
	if x.Commands == nil {
		return
	}
	for _, c := range x.Commands {
		if c.Aliases == nil {
			continue
		}
		for _, a := range c.Aliases {
			x._aliases[a] = c
		}
	}
}

// Run is for running a command within a specific runtime (shell) and
// performs completion if completion context is detected.  Otherwise, it
// executes the leaf Cmd returned from Seek calling its Method, and then
// Exits. Normally, Run is called from within main() to convert the Cmd
// into an actual executable program and normally it exits the program.
// Exiting can be controlled, however, with ExitOn/ExitOff when testing
// or for other purposes requiring multiple Run calls. Using Call
// instead will also just call the Cmd's Call Method without exiting.
// Note: Only bash runtime ("COMP_LINE") is currently supported, but
// others such a zsh and shell-less REPLs are planned.
func (x *Cmd) Run() {

	x.cacheAliases()

	// resolve bonzai.Aliases (if completion didn't replace them)
	if len(os.Args) > 1 {
		args := []string{os.Args[0]}
		alias := Aliases[os.Args[1]]
		if alias != nil {
			args = append(args, alias...)
			args = append(args, os.Args[2:]...)
			os.Args = args
		}
	}

	// bash completion context
	line := os.Getenv("COMP_LINE")
	if line != "" {
		var list []string
		lineargs := ArgsFrom(line)
		if len(lineargs) == 2 {
			list = append(list, maps.KeysWithPrefix(Aliases, lineargs[1])...)
		}
		cmd, args := x.Seek(lineargs[1:])
		if cmd.Completer == nil {
			list = append(list, comp.Standard(cmd, args...)...)
			if len(list) == 1 && len(lineargs) == 2 {
				if v, has := Aliases[list[0]]; has {
					fmt.Println(strings.Join(EscAll(v), " "))
					Exit()
				}
			}
			each.Println(list)
			Exit()
		}
		each.Println(cmd.Completer(cmd, args...))
		Exit()
	}

	// seek should never fail to return something, but ...
	cmd, args := x.Seek(os.Args[1:])
	if cmd == nil {
		ExitError(x.UsageError())
	}

	// default to first Command if no Call defined
	if cmd.Call == nil {
		if len(cmd.Commands) > 0 {
			cmd = cmd.Commands[0]
		} else {
			ExitError(x.Unimplemented())
		}
	}

	if len(args) < cmd.MinArgs {
		ExitError(cmd.UsageError())
	}

	// delegate
	if err := cmd.Call(cmd, args...); err != nil {
		ExitError(err)
	}
	Exit()
}

// UsageError returns an error with a single-line usage string.
func (x *Cmd) UsageError() error {
	return fmt.Errorf("usage: %v %v", x.Name, x.Usage)
}

// Unimplemented returns an error with a single-line usage string.
func (x *Cmd) Unimplemented() error {
	return fmt.Errorf("%q has not yet been implemented", x.Name)
}

// MissingConfig returns an error showing the expected configuration
// entry that is missing from the given path.
func (x *Cmd) MissingConfig(path string) error {
	return fmt.Errorf("missing config: %v", x.Branch()+"."+path)
}

// Add creates a new Cmd and sets the name and aliases and adds to
// Commands returning a reference to the new Cmd. The name must be
// first.
func (x *Cmd) Add(name string, aliases ...string) *Cmd {
	c := &Cmd{
		Name:    name,
		Aliases: aliases,
	}
	x.Commands = append(x.Commands, c)
	return c
}

// Cmd looks up a given Command by name or name from Aliases.
func (x *Cmd) Cmd(name string) *Cmd {
	if x.Commands == nil {
		return nil
	}
	for _, c := range x.Commands {
		if name == c.Name {
			return c
		}
	}
	if c, has := x._aliases[name]; has {
		return c
	}
	return nil
}

func (x *Cmd) CmdNames() []string {
	list := []string{}
	for _, c := range x.Commands {
		if c.Name == "" {
			continue
		}
		list = append(list, c.Name)
	}
	return list
}

// Param returns Param matching name if found, empty string if not.
func (x *Cmd) Param(p string) string {
	if x.Params == nil {
		return ""
	}
	for _, c := range x.Params {
		if p == c {
			return c
		}
	}
	return ""
}

// IsHidden returns true if the specified name is in the list of
// Hidden commands.
func (x *Cmd) IsHidden(name string) bool {
	if x.Hidden == nil {
		return false
	}
	for _, h := range x.Hidden {
		if h == name {
			return true
		}
	}
	return false
}

func (x *Cmd) Seek(args []string) (*Cmd, []string) {
	if args == nil || x.Commands == nil {
		return x, args
	}
	cur := x
	cur.Root = x
	cur.Conf = DefaultConfigurer
	n := 0
	for ; n < len(args); n++ {
		next := cur.Cmd(args[n])
		if next == nil {
			break
		}
		next.Caller = cur
		next.Root = x
		next.Conf = DefaultConfigurer
		cur = next
	}
	return cur, args[n:]
}

// Branch returns the dotted path to this command location within the
// parent tree. This is useful for associating configuration and other
// data specifically with this command. The branch path is determined by
// walking backward from current Caller up rather than depending on
// anything from the command line used to invoke the composing binary.
func (x *Cmd) Branch() string {
	callers := qstack.New[string]()
	callers.Unshift(x.Name)
	for p := x.Caller; p != nil; p = p.Caller {
		callers.Unshift(p.Name)
	}
	callers.Shift()
	return strings.Join(callers.Items(), ".")
}

// Q is a shorter version of x.Conf.Query(x.Root.Name,x.Branch()+"."+q)
// for convenience.
func (x *Cmd) Q(q string) string {
	return x.Conf.Query(x.Root.Name, "."+x.Branch()+"."+q)
}

// Log is currently short for log.Printf() but may be supplemented in
// the future to have more fine-grained control of logging.
func (x *Cmd) Log(format string, a ...any) {
	log.Printf(format, a...)
}

// TODO C for Cache lookups

// ---------------------- comp.Command interface ----------------------

// mostly to overcome cyclical imports

// GetName fulfills the comp.Command interface.
func (x *Cmd) GetName() string { return x.Name }

// GetCommands fulfills the comp.Command interface.
func (x *Cmd) GetCommands() []string { return x.CmdNames() }

// GetHidden fulfills the comp.Command interface.
func (x *Cmd) GetHidden() []string { return x.Hidden }

// GetParams fulfills the comp.Command interface.
func (x *Cmd) GetParams() []string { return x.Params }

// GetOther fulfills the comp.Command interface.
func (x *Cmd) GetOther() map[string]string { return x.Other }

// GetCompleter fulfills the Command interface.
func (x *Cmd) GetCompleter() comp.Completer { return x.Completer }

// GetCaller fulfills the comp.Command interface.
func (x *Cmd) GetCaller() comp.Command { return x.Caller }
