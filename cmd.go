// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/rwxrob/bonzai/ds/qstack"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/is"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

type Cmd struct {
	Name  string // ex: delete
	Alias string // ex: rm|d|del
	Opts  string // ex: mon|wed|fri

	// Work done by this command
	Init Method // run-time initialization/validation
	Call Method // if nil, Def must be set

	// Delegation to subcommands
	Def  *Cmd   // default [Cmd] if no Call and no matching Cmds
	Cmds []*Cmd // compiled/composed commands in no particular order
	Hide string // disable completion: ex: old|defunct

	// Minimal when [Cmd.Docs] is overkill
	Usage string
	Vers  string
	Short string
	Long  string

	// Faster than lots of "if" conditions in [Cmd.Call]. Consider
	// [Cmd.Init] when more complex argument validation is needed.
	MinArgs   int    // min
	MaxArgs   int    // max
	NumArgs   int    // exact, also used for NoArg (0)
	MatchArgs string // PCRE/Go regular expression (requires Usage)

	// Self-completion support: complete -C foo foo
	Comp Completer

	// Template commands/functions to be added (or overwrite) the internal
	// [FuncMap] collection of template commands used by the [Cmd.Fill]
	// command. These apply to this [Cmd] only and will not be available
	// to subcommands. To share between commands (including subcommands)
	// assign the same FuncMap to all of them.
	FuncMap template.FuncMap

	// Following are never assigned in declarations but are instead
	// set at [Run] time for use in [Call] methods:
	Caller *Cmd // delegation

	aliases  []string        // see [CacheAlias]
	opts     []string        // see [CacheOpts]
	hidden   []string        // see [CacheHide]
	cmdAlias map[string]*Cmd // see [CacheCmdAlias]
}

// Method defines the main code to execute for a command [Cmd.Call]. By
// convention the parameter list should be named "args" and the caller
// "x". If either is unused an underscore should be used instead.
type Method func(x *Cmd, args ...string) error

func (x Cmd) WithName(name string) *Cmd {
	x.Name = name
	return &x
}

// Names returns slice of all [Alias] and the [Name] as the last item.
func (x *Cmd) Names() []string {
	var names []string
	for _, alias := range x.AliasSlice() {
		if len(alias) == 0 {
			continue
		}
		names = append(names, alias)
	}
	names = append(names, x.Name)
	return names
}

// IsValidName is assigned a function that returns a boolean
// for the given name. See [is.AllLatinASCIILower] for an example. Note
// that if this is changed certain characters may break the
// creation of multicall binary links and bash completion.
var IsValidName = is.AllLatinASCIILowerWithDashes

// CacheCmdAlias splits the [Cmd.Alias] for each [Cmd] in
// [Cmds] with its respective [Cmd.AliasSlice] and assigns them
// the [Cmd.CmdAliasMap] cache map. This is primarily used for bash
// tab completion support in [Run] and use as a multicall binary. If
// [Cmds] is nil or [Name] is empty silently returns.
func (x *Cmd) CacheCmdAlias() {
	x.cmdAlias = map[string]*Cmd{}
	if x.Cmds == nil || len(x.Name) == 0 {
		return
	}
	for _, c := range x.Cmds {
		aliases := c.AliasSlice()
		if len(aliases) == 0 {
			continue
		}
		for _, a := range aliases {
			x.cmdAlias[a] = c
		}
	}
}

// CmdAliasMap calls [CacheCmdAlias] to update cache if it
// is nil and then returns it. [Hide] is not applied.
func (x *Cmd) CmdAliasMap() map[string]*Cmd {
	if x.cmdAlias == nil {
		x.CacheCmdAlias()
	}
	return x.cmdAlias
}

// CacheOpts updates the [opts] cache by splitting [Opts]. Remember
// to call this whenever dynamically altering the value at
// runtime.
func (x *Cmd) CacheOpts() {
	if len(x.Opts) > 0 {
		x.opts = strings.Split(x.Opts, `|`)
		return
	}
	x.opts = []string{}
}

// OptsSlice updates the [params] internal cache ([CacheOpts]) and
// returns it as a slice.
func (x *Cmd) OptsSlice() []string {
	if x.opts == nil {
		x.CacheOpts()
	}
	return x.opts
}

// CacheAlias updates the [aliases] cache by splitting [Alias]
// and adding the [Name] to the end. Remember to call this whenever
// dynamically altering the value at runtime.
func (x *Cmd) CacheAlias() {
	if len(x.Alias) > 0 {
		x.aliases = strings.Split(x.Alias, `|`)
		for _, alias := range x.aliases {
			if !IsValidName(alias) {
				run.ExitError(InvalidName{alias})
				return
			}
		}
		return
	}
	x.aliases = []string{}
}

// AliasSlice updates the [aliases] internal cache ([CacheAlias]) and
// returns it as a slice.
func (x *Cmd) AliasSlice() []string {
	if x.aliases == nil {
		x.CacheAlias()
	}
	return x.aliases
}

// CacheHide updates the [hidden] cache by splitting [Hide]
// . Remember to call this whenever dynamically altering the value at
// runtime.
func (x *Cmd) CacheHide() {
	if len(x.Hide) > 0 {
		x.hidden = strings.Split(x.Hide, `|`)
		return
	}
	x.hidden = []string{}
}

// HideSlice updates the [hidden] internal cache ([CacheHide]) and
// returns it as a slice.
func (x *Cmd) HideSlice() []string {
	if x.hidden == nil {
		x.CacheHide()
	}
	return x.hidden
}

// Run method resolves [Cmd.Alias] and seeks the leaf [Cmd]. It then
// calls the leaf's first-class [Cmd.Call] function passing itself as
// the first argument along with any remaining command line arguments.
// Run returns nothing because it usually exits the program. Normally,
// Run is called from within main() to convert the Cmd into an actual
// executable program. Use Call instead of Run when delegation is
// needed. However, avoid tight-coupling that comes from delegation with
// Call when possible. Also, Call automatically assumes the proper
// number and type of arguments have already been checked (see
// [Cmd.MinArgs], etc.) which is normally done by Run.
//
// # Completion
//
// Since Run is the main execution entry point for all Bonzai command
// trees it is also responsible for handling bash completion. Only bash
// completion is supported within the binary itself because only bash
// provides self-completion (complete -C foo foo). However, zsh can also be
// made to support it by adding a few functions from the
// oh-my-zsh code base).
//
// Completion mode is triggered by the detection of the COMP_LINE
// environment variable. (complete -C cmd cmd).
//
// When COMP_LINE is set, Run prints a list of possible completions to
// standard output by calling the [Completer.Complete] function of its
// [Comp] field. If [Comp] is nil no completion is attempted. Each
// [Cmd] explicitly manages its own completion and can draw from an
// growing ecosystem of Completers or assign its own. See the
// [core/comp] package for more examples.
//
// # Multicall binary and links
//
// Popularized by BusyBox/Alpine, a multicall binary is a single
// executable that behaves differently based on its name, either through
// copying the binary to another name, or linking (symbolic or hard).
// All Bonzai compiled binaries automatically behave as multicall
// binaries provided the name of the actual binary or link matches the
// name of [Cmd].
//
// Note that this method should never be used to obscure a highly
// sensitive command thinking it won't be discovered. Discovering every
// possible command is very easy to brute force.
//
// # Never panic
//
// All panics are trapped with [run.TrapPanic] which normally exits with 1 and
// outputs the main message. See [run.TrapPanic] for details.
//
// # Valid name check
//
// Throws [InvalidName] error and exist if [Name] does not pass
// [InvalidName] check.
//
// # Subcommand optional arguments
//
// If any argument is detected, delegation through recursive Run calls
// to subcommands is attempted. If more than one argument, each
// argument is assumed to be a [Name] or alias from [Alias] and so on
// (see [Can] for details).
func (x *Cmd) Run(args ...string) {
	defer run.TrapPanic()
	x.recurseIfArgs(args)
	x.exitUnlessValidName()
	x.recurseIfMulti(args)
	x.detectCompletion(args)
	c, args := x.Seek(os.Args[1:])
	if c == nil {
		run.ExitError(IncorrectUsage{c})
		return
	}
	c.init(args)
	c.exitUnlessCallable()
	c.exitIfBadArgs(args)
	c.call(args)
}

// IsHidden returns true if one of [Hide] matches the [Name].
func (x *Cmd) IsHidden() bool {
	for _, hidden := range x.HideSlice() {
		if hidden == x.Name {
			return true
		}
	}
	return false
}

func (x *Cmd) has(c *Cmd) bool {
	for _, this := range x.Cmds {
		if this.Name == c.Name {
			return true
		}
	}
	return false
}

func (x *Cmd) call(args []string) {
	if x.Caller == nil {
		x.Caller = x
	}
	if x.Call != nil {
		if err := x.Call(x, args...); err != nil {
			run.ExitError(err)
			return
		}
		run.Exit()
		return
	}
	if x.Def != nil {
		x.Def.call(args)
		return
	}
	run.Exit()
	return
}

// default to first Cmds if no Call defined
func (x *Cmd) exitUnlessCallable() {
	switch {
	case x.Call == nil && x.Def == nil:
		run.ExitError(Uncallable{x})
		return
	case x.Call != nil && x.Def != nil:
		run.ExitError(CallOrDef{x})
		return
	}
}

// initialize before delegation and Call
func (x *Cmd) init(args []string) {
	if x.Init != nil {
		if err := x.Init(x, args...); err != nil {
			run.ExitError(err)
		}
	}
}

func (x *Cmd) recurseIfArgs(args []string) {
	argslen := len(args)
	if argslen > 0 {
		if argslen == 1 && strings.Contains(args[0], `-`) {
			args = strings.Split(args[0], `-`)
		}
		if c := x.Can(args...); c != nil {
			c.Run()
			return
		}
	}
}

func (x *Cmd) exitIfBadArgs(args []string) {
	switch {
	case len(args) < x.MinArgs:
		run.ExitError(NotEnoughArgs{len(args), x.MinArgs})
		return
	case x.MaxArgs > 0 && len(args) > x.MaxArgs:
		run.ExitError(TooManyArgs{len(args), x.MaxArgs})
		return
	case x.NumArgs > 0 && len(args) != x.NumArgs:
		run.ExitError(WrongNumArgs{len(args), x.NumArgs})
		return
	}
}

func (x *Cmd) exitUnlessValidName() {
	if !IsValidName(x.Name) {
		run.ExitError(InvalidName{x.Name})
		return
	}
}

// called as multicall binary
func (x *Cmd) recurseIfMulti(args []string) {
	name, _ := run.ExeName()
	if name != x.Name {
		if c := x.Can(name); c != nil {
			c.Run()
			return
		}
	}
}

// complete -C foo foo (man bash, Programmable Completion)
func (x *Cmd) detectCompletion(args []string) {
	if line := os.Getenv("COMP_LINE"); len(line) > 0 {

		// find the leaf command
		lineargs := run.ArgsFrom(line)
		cmd, args := x.Seek(lineargs[1:])

		// default completer or package aliases, always exits
		if cmd.Comp == nil {
			run.Exit()
			return
		}

		// not sure we've completed the command name itself yet
		if len(args) == 0 {
			fmt.Println(cmd.Name)
			run.Exit()
			return
		}

		// own completer, delegate
		each.Println(cmd.Comp.Complete(*cmd, args...))
		run.Exit()
		return
	}
}

// String fulfills the [fmt.Stringer] interface for debugging.
func (x Cmd) String() string { return x.Name }

// Root returns the root [Cmd] from the current [Path]. This must always
// be calculated every time since any Cmd can change positions and
// pedigrees at any time at run time. Returns self if no [PathCmds]
// found.
func (x *Cmd) Root() *Cmd {
	cmds := x.PathCmds()
	if len(cmds) > 0 {
		return cmds[0].Caller
	}
	return x.Caller
}

// PrependCmd safely prepends the passed [*Cmd] to the [Cmds] slice.
func (x *Cmd) PrependCmd(cmd *Cmd) {
	old := x.Cmds
	x.Cmds = []*Cmd{cmd}
	if old != nil {
		x.Cmds = append(x.Cmds, old...)
	}
}

// AppendCmd safely appends the passed [*Cmd] to the [Cmds] slice.
func (x *Cmd) AppendCmd(cmd *Cmd) {
	if x.Cmds == nil {
		x.Cmds = []*Cmd{}
	}
	x.Cmds = append(x.Cmds, cmd)
}

// Add creates a new Cmd and sets the [Name] and [Alias] and adds to
// [Cmds] returning a reference to the new Cmd. Name must be
// first.
func (x *Cmd) Add(name string, aliases ...string) *Cmd {
	c := &Cmd{
		Name:  name,
		Alias: strings.Join(aliases, `|`),
	}
	x.aliases = aliases
	x.Cmds = append(x.Cmds, c)
	return c
}

// Resolve looks up a given [Cmd] by name or alias from [Alias]
// (caching a lookup map of aliases in the process).
func (x *Cmd) Resolve(name string) *Cmd {
	if x.Cmds == nil {
		return nil
	}

	for _, c := range x.Cmds {
		if name == c.Name {
			return c
		}
	}

	aliases := x.CmdAliasMap()
	if c, has := aliases[name]; has {
		return c
	}
	return nil
}

// Can returns the [*Cmd] from [Cmds] if the [Cmd.Name] or any
// alias in [Cmd.Alias] for that command matches the name passed. If
// more than one argument is passed calls itself recursively on each
// item in the list.
func (x *Cmd) Can(names ...string) *Cmd {
	var name string
	switch len(names) {
	case 0:
		return nil
	case 1:
		return x.can(names[0])
	}
	name = names[0]
	names = names[1:]
	c := x.can(name)
	if c == nil {
		return nil
	}
	return c.Can(names...)
}

func (x *Cmd) can(name string) *Cmd {
	for _, c := range x.Cmds {
		if c.Name == name {
			return c
		}
	}
	aliases := x.CmdAliasMap() // to trigger cache if needed
	if c, has := aliases[name]; has {
		return c
	}
	return nil
}

// CmdNames returns the names of every [Cmd] from [Cmds]
func (x *Cmd) CmdNames() []string {
	list := []string{}
	for _, c := range x.Cmds {
		if c.Name == "" {
			continue
		}
		list = append(list, c.Name)
	}
	return list
}

// Fill fills out the [text/template] string using the [Cmd] data fields
// and [Cmd.FuncMap] values combined with [bonzai.FuncMap].
func (x *Cmd) Fill(tmpl string) string {
	funcs := to.MergedMaps(FuncMap, x.FuncMap)
	t, err := template.New("t").Funcs(funcs).Parse(tmpl)
	if err != nil {
		log.Println(err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, x); err != nil {
		log.Println(err)
	}
	return buf.String()
}

// Print calls [Fill] on string and prints it with [fmt.Print]. This is
// a rather expensive operation by comparison. Consider the simpler
// alternative or [term.Print].
func (x *Cmd) Print(tmpl string) { fmt.Print(x.Fill(tmpl)) }

// Println calls [Fill] on string and prints it with [fmt.Println]. This is
// a rather expensive operation by comparison. Consider the simpler
// alternative or [term.Print].
func (x *Cmd) Println(tmpl string) { fmt.Println(x.Fill(tmpl)) }

// Opt returns the [Opts] entry matching name if found, empty string if not.
func (x *Cmd) Opt(p string) string {
	if x.opts == nil {
		x.CacheOpts()
	}
	for _, c := range x.opts {
		if p == c {
			return c
		}
	}
	return ""
}

// Seek checks the args for command names returning the deepest along
// with the remaining arguments. Typically the args passed are directly
// from the command line. Seek also sets the Caller on each Cmd found
// during resolution.
func (x *Cmd) Seek(args []string) (*Cmd, []string) {
	if (len(args) == 1 && args[0] == "") || x.Cmds == nil {
		return x, args
	}
	cur := x
	n := 0
	for ; n < len(args); n++ {
		next := cur.Resolve(args[n])
		if next == nil {
			break
		}
		next.Caller = cur
		cur = next
	}
	return cur, args[n:]
}

// PathCmds returns the path of commands used to arrive at this
// command. The path is determined by walking backward from current
// Caller up rather than depending on anything from the command line
// used to invoke the composing binary. Also see [PathNames].
func (x *Cmd) PathCmds() []*Cmd {
	path := qstack.New[*Cmd]()
	path.Unshift(x)
	for p := x.Caller; p != nil; p = p.Caller {
		path.Unshift(p)
	}
	path.Shift()
	return path.Items()
}

// PathNames returns the path of command names used to arrive at this
// command. The path is determined by walking backward from current
// Caller up rather than depending on anything from the command line
// used to invoke the composing binary. Also see Path.
func (x *Cmd) PathNames() []string {
	path := qstack.New[string]()
	path.Unshift(x.Name)
	p := x.Caller
	for p != nil {
		path.Unshift(p.Name)
		if p == p.Caller {
			break
		}
		p = p.Caller
	}
	path.Shift()
	return path.Items()
}
