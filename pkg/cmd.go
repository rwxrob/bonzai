// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/rwxrob/bonzai/pkg/core/ds/qstack"
	"github.com/rwxrob/bonzai/pkg/core/fn/each"
	"github.com/rwxrob/bonzai/pkg/core/is"
	"github.com/rwxrob/bonzai/pkg/core/run"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

type Cmd struct {
	Name    string // ex: delete
	Aliases string // ex: rm|d|del
	Params  string // ex: mon|wed|fri

	// minimal when DocFS overkill
	Usage   string
	Summary string
	Version string

	// Faster than lots of "if" conditions in [Call]
	MinArgs int
	MaxArgs int
	NumArgs int
	NoArgs  bool
	MinParm int
	MaxParm int

	// Descending tree of delegated commands
	Cmds   []*Cmd // delegated, first is always default
	Hidden string // disables completion for Cmds

	// Bash completion support (only)
	Comp Completer

	// Default vars declaration and initial values required by [vars.Cmd]
	// (if used). Does not overwrite existing vars.
	InitVars map[string]string

	// Optional embedded documentation in any format used by help and
	// documentation commands such as [doc.Cmd] from the bonzai/core/cmds
	// package. Embedded content is usually lazy loaded only when the doc
	// command is called. Structure of format of the files can be anything
	// supported by any [Cmd] but Bonzai [mark] is recommended for
	// greatest compatibility. Use of an embedded file system instead of
	// a string allows, for example, support for multiple languages to be
	// embedded into a single binary.
	DocFS embed.FS

	// Functions to be used for the Fill command which is automatically
	// called on most string properties (ex: {{ exename }})
	FuncMap template.FuncMap

	// Where the work happens
	Init   Method // before Cmds/Call, good for validation
	Call   Method // optional if Cmds
	Caller *Cmd

	// Pass bulk input efficiently (when args won't do)
	Input io.Reader

	aliases    []string        // see [CacheAliases]
	params     []string        // see [CacheParams]
	hidden     []string        // see [CacheHidden]
	cmdAliases map[string]*Cmd // see [CacheCmdAliases]
}

// Method defines the main code to execute for a command [Cmd.Call]. By
// convention the parameter list should be named "args" and the caller
// "x". If either is unused an underscore should be used instead.
type Method func(x *Cmd, args ...string) error

// Names returns slice of all [Aliases] and the [Name] as the last item.
func (x *Cmd) Names() []string {
	var names []string
	for _, alias := range x.AliasesSlice() {
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
var IsValidName = is.AllLatinASCIILower

// CacheCmdAliases splits the [Cmd.Aliases] for each [Cmd] in
// [Cmds] with its respective [Cmd.AliasesSlice] and assigns them
// the [Cmd.CmdAliasesMap] cache map. This is primarily used for bash
// tab completion support in [Run] and use as a multicall binary. If
// [Cmds] is nil or [Name] is empty silently returns.
func (x *Cmd) CacheCmdAliases() {
	x.cmdAliases = map[string]*Cmd{}
	if x.Cmds == nil || len(x.Name) == 0 {
		return
	}
	for _, c := range x.Cmds {
		aliases := c.AliasesSlice()
		if len(aliases) == 0 {
			continue
		}
		for _, a := range aliases {
			x.cmdAliases[a] = c
		}
	}
}

// CmdAliasesMap calls [CacheCmdAliases] to update cache if it
// is nil and then returns it. [Hidden] is not applied.
func (x *Cmd) CmdAliasesMap() map[string]*Cmd {
	if x.cmdAliases == nil {
		x.CacheCmdAliases()
	}
	return x.cmdAliases
}

// CacheParams updates the [params] cache by splitting [Params]
// . Remember to call this whenever dynamically altering the value at
// runtime.
func (x *Cmd) CacheParams() {
	if len(x.Params) > 0 {
		x.params = strings.Split(x.Params, `|`)
		return
	}
	x.params = []string{}
}

// ParamsSlice updates the [params] internal cache ([CacheParams]) and
// returns it as a slice.
func (x *Cmd) ParamsSlice() []string {
	if x.params == nil {
		x.CacheParams()
	}
	return x.params
}

// CacheAliases updates the [aliases] cache by splitting [Aliases]
// and adding the [Name] to the end. Remember to call this whenever
// dynamically altering the value at runtime.
func (x *Cmd) CacheAliases() {
	if len(x.Aliases) > 0 {
		x.aliases = strings.Split(x.Aliases, `|`)
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

// AliasesSlice updates the [aliases] internal cache ([CacheAliases]) and
// returns it as a slice.
func (x *Cmd) AliasesSlice() []string {
	if x.aliases == nil {
		x.CacheAliases()
	}
	return x.aliases
}

// CacheHidden updates the [hidden] cache by splitting [Hidden]
// . Remember to call this whenever dynamically altering the value at
// runtime.
func (x *Cmd) CacheHidden() {
	if len(x.Hidden) > 0 {
		x.hidden = strings.Split(x.Hidden, `|`)
		return
	}
	x.hidden = []string{}
}

// HiddenSlice updates the [hidden] internal cache ([CacheHidden]) and
// returns it as a slice.
func (x *Cmd) HiddenSlice() []string {
	if x.hidden == nil {
		x.CacheHidden()
	}
	return x.hidden
}

// Run method resolves [Cmd.Aliases] and seeks the leaf [Cmd]. It then
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
// Completion mode is triggered by the detection of the bash shell and
// the COMP_LINE environment variable. (complete -C cmd cmd) that is
// unique to bash. All other shells can use the "help json" structured
// data create external completion scripts.
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
// argument is assumed to be a [Name] or alias from [Aliases] and so on
// (see [Can] for details). As a convenience, if only one argument is
// passed and that argument contains a dash, it is assumed to be
// a [PathWithDashes] and is split and expanded into a new args list as
// if every field where passed as separate strings instead.
func (x *Cmd) Run(args ...string) {
	defer run.TrapPanic()
	x.recurseIfArgs(args)
	x.exitUnlessValidName()
	x.recurseIfMulti(args)
	x.detectBashCompletion(args)
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

// IsHidden returns true if one of [Hidden] matches the [Name].
func (x *Cmd) IsHidden() bool {
	for _, hidden := range x.HiddenSlice() {
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
	if err := x.Call(x, args...); err != nil {
		run.ExitError(err)
		return
	}
	run.Exit()
}

// default to first Cmds if no Call defined
func (x *Cmd) exitUnlessCallable() {
	if x.Call == nil {
		if len(x.Cmds) > 0 {
			fcmd := x.Cmds[0]
			if fcmd.Call == nil {
				run.ExitError(DefCmdReqCall{x})
				return
			}
			fcmd.Caller = x
			x = fcmd
		} else {
			run.ExitError(NoCallNoCmds{x})
			return
		}
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
	case len(args) > 0 && x.NoArgs:
		run.ExitError(TooManyArgs{len(args), 0})
		return
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

func (x *Cmd) recurseIfMulti(args []string) {
	// called as multicall binary
	name := ExeName
	if name == x.Name {
		name = ExeSymLink
	}
	if name != x.Name {
		// dashed/long (ex: z-bon-multi-symlink)
		if strings.Contains(name, `-`) {
			args = strings.Split(name, `-`)
			first := args[0]
			if first != ExeName {
				run.ExitError(InvalidMultiName{ExeName, name})
			}
			x.Run(args[1:]...)
			return
		}
		// simple (ex: bon)
		if c := x.Can(name); c != nil {
			c.Run()
			return
		}
	}
}

// complete -C foo foo (man bash, Programmable Completion)
func (x *Cmd) detectBashCompletion(args []string) {

	if line := os.Getenv("COMP_LINE"); len(line) > 0 && run.ShellIsBash() {

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
			fmt.Println(" ")
			return
		}

		// own completer, delegate
		each.Println(cmd.Comp.Complete(cmd, args...))
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

// Add creates a new Cmd and sets the [Name] and [Aliases] and adds to
// [Cmds] returning a reference to the new Cmd. Name must be
// first.
func (x *Cmd) Add(name string, aliases ...string) *Cmd {
	c := &Cmd{
		Name:    name,
		Aliases: strings.Join(aliases, `|`),
	}
	x.aliases = aliases
	x.Cmds = append(x.Cmds, c)
	return c
}

// Resolve looks up a given [Cmd] by name or alias from [Aliases]
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

	aliases := x.CmdAliasesMap()
	if c, has := aliases[name]; has {
		return c
	}
	return nil
}

// Can returns the [*Cmd] from [Cmds] if the [Cmd.Name] or any
// alias in [Cmd.Aliases] for that command matches the name passed. If
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
	aliases := x.CmdAliasesMap() // to trigger cache if needed
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

// Param returns Param matching name if found, empty string if not.
func (x *Cmd) Param(p string) string {
	if x.params == nil {
		x.CacheParams()
	}
	for _, c := range x.params {
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
// used to invoke the composing binary. Also see [Path], [PathNames].
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

// Path returns a dotted notation of the [PathNames] including an initial
// dot (for root). This is useful for associating configuration and other
// data specifically with this command. If any arguments are passed then
// will be added with dots between them.
func (x *Cmd) Path(more ...string) string {
	if len(more) > 0 {
		list := x.PathNames()
		list = append(list, more...)
		return "." + strings.Join(list, ".")
	}
	return "." + strings.Join(x.PathNames(), ".")
}

// PathWithDashes is the same as [Path] but with dashes/hyphens instead and
// without the leading dot.
func (x *Cmd) PathWithDashes(more ...string) string {
	path := x.Path(more...)
	return path[1:]
}

// Get is a shorter version of Vars.Get(x.Path()+"."+key) which fetches
// and returns persisted cache values (see [InitVars] and [VarsDriver]).
// If a value has not yet been assigned returns the value from [InitVars]
// and sets it with [Set]. All var keys must be declared and assigned
// initial values with [InitVars] or they cannot be used and throw
// an [UnsupportedVar] [run.ExitError].
func (x *Cmd) Get(key string) string {
	defval, declared := x.InitVars[key]
	if !declared {
		run.ExitError(UnsupportedVar{key})
		return ""
	}
	path := x.Path()
	if path != "." {
		path += "."
	}
	ptr := path + key
	val, code := Vars.Get(ptr)
	if code == NOTFOUND {
		Vars.Set(ptr, defval)
		val = defval
	}
	return val
}

// Set is shorter version of Vars.Set(x.Path()+"."+key.val).
func (x *Cmd) Set(key, val string) {
	path := x.Path()
	if path != "." {
		path += "."
	}
	Vars.Set(path+key, val)
}
