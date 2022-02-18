/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bonzai

import (
	"os"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/filter"
)

// Cmd is a struct the easier to use and read when creating
// implementations of the Command interface.
//
// Params
//
// Params require a Method. While Methods may receive any number of
// arguments, Params are a way of helping completion for regular
// parameters. Standard completion will not recursively complete
// multiple params, one param per completion.
type Cmd struct {
	Name        string         `json:"name,omitempty"`
	Aliases     []string       `json:"aliases,omitempty"`
	Summary     string         `json:"summary,omitempty"`
	Usage       string         `json:"usage,omitempty"`
	Version     string         `json:"version,omitempty"`
	Copyright   string         `json:"copyright,omitempty"`
	License     string         `json:"license,omitempty"`
	Description string         `json:"description,omitempty"`
	Site        string         `json:"site,omitempty"`
	Source      string         `json:"source,omitempty"`
	Issues      string         `json:"issues,omitempty"`
	Caller      *Cmd           `json:"-"`
	Commands    []*Cmd         `json:"commands,omitempty"`
	Params      []string       `json:"params,omitempty"`
	Hidden      []string       `json:"hide,omitempty"`
	Completer   comp.Completer `json:"-"`
	Method      Method         `json:"-"`
}

// Run detects the runtime (shell) and performs completion and exists if
// completion context is detected. Otherwise, it executes the leaf Cmd
// returned from Seek calling its Method, and then Exits. Normally, Run
// is called from within main() to convert the Cmd into an actual
// executable program and normally it exits the program. Exiting can be
// controlled, however, with ExitOn/ExitOff when testing or for other
// purposes requiring multiple Run calls.
//
// Note: Only bash runtime ("COMP_LINE") is currently supported, but
// others such a zsh and shell-less REPLs are planned.
func (x *Cmd) Run() {

	// TODO add completion for other runtimes
	line := os.Getenv("COMP_LINE")
	if line != "" {
		cmd, args := x.Seek(ArgsFrom(line)[1:])
		if cmd.Completer == nil {
			comp.Standard(cmd, args)
		}
		cmd.Completer(cmd, args)
		Exit()
	}

	log.Print("would execute")

	// TODO execute Command with left-over args if match

	// TODO call x.Method with args[1:] if x.Method not nil

	// TODO fail with help (if found) or generic usage error

}

func (x *Cmd) initRuntime() {
	/*
		if os.Getenv("BASH_VERSION") == "" {
			x.runtime = run.Bash{}
			return
		}
	*/
	return
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

func (x *Cmd) Cmd(name string) *Cmd {
	if x.Commands == nil {
		return nil
	}
	for _, c := range x.Commands {
		if name == c.Name {
			return c
		}
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
	if args == nil || len(args) == 0 || x.Commands == nil || len(x.Commands) == 0 {
		return x, args
	}
	cur := x
	n := 0
	for ; n < len(args); n++ {
		next := cur.Cmd(args[n])
		if next == nil {
			break
		}
		cur = next
	}
	return cur, args[n:]
}

// ---------------------- comp.Command interface ----------------------

// mostly to overcome cyclical imports

// GetName fulfills the comp.Command interface.
func (x *Cmd) GetName() string { return x.Name }

// GetCaller fulfills the comp.Command interface.
func (x *Cmd) GetCaller() comp.Command { return x.Caller }

// GetCommands fulfills the comp.Command interface.
func (x *Cmd) GetCommands() []string { return x.CmdNames() }

// GetHidden fulfills the comp.Command interface.
func (x *Cmd) GetHidden() []string { return x.Hidden }

// GetParams fulfills the comp.Command interface.
func (x *Cmd) GetParams() []string { return x.Params }

// GetCompleter fulfills the Command interface.
func (x *Cmd) GetCompleter() comp.Completer { return x.Completer }
