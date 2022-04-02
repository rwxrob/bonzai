// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/conf"
	"github.com/rwxrob/fn/each"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/fn/redu"
	"github.com/rwxrob/structs/qstack"
)

type Cmd struct {
	Name        string   `json:"name,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Usage       string   `json:"usage,omitempty"`
	Version     string   `json:"version,omitempty"`
	Copyright   string   `json:"copyright,omitempty"`
	License     string   `json:"license,omitempty"`
	Description string   `json:"description,omitempty"`
	Site        string   `json:"site,omitempty"`
	Source      string   `json:"source,omitempty"`
	Issues      string   `json:"issues,omitempty"`
	Commands    []*Cmd   `json:"commands,omitempty"`
	Params      []string `json:"params,omitempty"`
	Hidden      []string `json:"hidden,omitempty"`

	Other map[string]string `json:"other,omitempty"`

	Completer comp.Completer      `json:"-"`
	Conf      conf.Configurer     `json:"-"`
	UsageFunc func(x *Cmd) string `json:"-"`

	Root    *Cmd   `json:"-"`
	Caller  *Cmd   `json:"-"`
	Call    Method `json:"-"`
	MinArgs int    `json:"-"` // minimum number of args required (including parms)
	MinParm int    `json:"-"` // minimum number of params required
	MaxParm int    `json:"-"` // maximum number of params required

	_aliases map[string]*Cmd
}

// Names returns the Name and any Aliases grouped such that the Name is
// always last.
func (x *Cmd) Names() []string {
	var names []string
	names = append(names, x.Aliases...)
	names = append(names, x.Name)
	return names
}

// UsageNames returns single name, joined Names with bar (|) and wrapped
// in parentheses, or empty string if no names.
func (x *Cmd) UsageNames() string { return UsageGroup(x.Names(), 1, 1) }

// UsageParams returns the Params in UsageGroup notation.
func (x *Cmd) UsageParams() string {
	return UsageGroup(x.Params, x.MinParm, x.MaxParm)
}

// UsageCmdNames returns the Names for each of its Commands joined, if
// more than one, with usage regex notation.
func (x *Cmd) UsageCmdNames() string {
	var names []string
	for _, n := range x.Commands {
		names = append(names, n.UsageNames())
	}
	return UsageGroup(names, 1, 1)
}

// Title returns a dynamic field of Name and Summary combined (if
// exists). If the Name field of the commands is not defined will return
// a "{ERROR}".
func (x *Cmd) Title() string {
	if x.Name == "" {
		return "{ERROR: Name is empty}"
	}
	switch {
	case len(x.Summary) > 0 && len(x.Version) > 0:
		return x.Name + " (" + x.Version + ")" + " - " + x.Summary
	case len(x.Summary) > 0:
		return x.Name + " - " + x.Summary
	default:
		return x.Name
	}
}

// Legal returns a single line with the combined values of the
// Name, Version, Copyright, and License. If Version is empty or nil an
// empty string is returned instead. Legal() is used by the
// version builtin command to aggregate all the version information into
// a single output.
func (x *Cmd) Legal() string {
	switch {
	case len(x.Copyright) > 0 && len(x.License) == 0 && len(x.Version) == 0:
		return x.Name + " " + x.Copyright
	case len(x.Copyright) > 0 && len(x.License) > 0 && len(x.Version) > 0:
		return x.Name + " (" + x.Version + ") " +
			x.Copyright + "\nLicense " + x.License
	case len(x.Copyright) > 0 && len(x.License) > 0:
		return x.Name + " " + x.Copyright + "\nLicense " + x.License
	case len(x.Copyright) > 0 && len(x.Version) > 0:
		return x.Name + " (" + x.Version + ") " + x.Copyright
	case len(x.Copyright) > 0:
		return x.Name + "\n" + x.Copyright
	default:
		return ""
	}
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

	// resolve Z.Aliases (if completion didn't replace them)
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

// UsageError returns an error with a single-line usage string. The word
// "usage" can be changed by assigning Z.UsageText to something else.
// The commands own UsageFunc will be used if defined. If undefined, the
// Z.DefaultUsageFunc will be used instead (which can also be assigned
// to something else if needed).
func (x *Cmd) UsageError() error {
	usage := x.UsageFunc
	if usage == nil {
		usage = DefaultUsageFunc
	}
	return fmt.Errorf("%v: %v %v", UsageText, x.Name, usage(x))
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

// Resolve looks up a given Command by name or name from Aliases.
func (x *Cmd) Resolve(name string) *Cmd {
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

// CmdNames returns the names of every Command.
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

// UsageCmdTitles returns a single string with the titles of each
// subcommand indented and with a maximum title signature length for
// justification.  Hidden commands are not included. Note that the order
// of the Commands is preserved (not necessarily alphabetic).
func (x *Cmd) UsageCmdTitles() string {
	var set []string
	var summaries []string
	for _, c := range x.Commands {
		set = append(set, strings.Join(c.Names(), "|"))
		summaries = append(summaries, c.Summary)
	}
	longest := redu.Longest(set)
	var buf string
	for n := 0; n < len(set); n++ {
		if len(summaries[n]) > 0 {
			buf += fmt.Sprintf(`%-`+strconv.Itoa(longest)+"v - %v\n", set[n], summaries[n])
		} else {
			buf += fmt.Sprintf(`%-`+strconv.Itoa(longest)+"v\n", set[n])
		}
	}
	return buf
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
		next := cur.Resolve(args[n])
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
