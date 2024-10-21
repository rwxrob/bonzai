// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package bonzai

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/rwxrob/bonzai/pkg/core/ds/qstack"
	"github.com/rwxrob/bonzai/pkg/core/fn/each"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
	"github.com/rwxrob/bonzai/pkg/core/fn/redu"
	"github.com/rwxrob/bonzai/pkg/core/run"
	"github.com/rwxrob/bonzai/pkg/core/to"
)

type Cmd struct {
	Name    string   `json:"name,omitempty"`    // plain
	Aliases []string `json:"aliases,omitempty"` // plain
	DocFrom string   `json:"docfrom,omitempty"` // template, embedded

	// Directly assign (instead of parse from DocFrom)
	Usage     string `json:"usage,omitempty"`       // template
	Short     string `json:"summary,omitempty"`     // template
	Long      string `json:"description,omitempty"` // template
	Version   string `json:"version,omitempty"`     // template
	Copyright string `json:"copyright,omitempty"`   // template
	License   string `json:"license,omitempty"`     // template
	Site      string `json:"site,omitempty"`        // template (url)
	Source    string `json:"source,omitempty"`      // template (url)
	Issues    string `json:"issues,omitempty"`      // template (url)

	// Descending tree of delegated commands
	Commands []*Cmd   `json:"commands,omitempty"` // delegated
	Hidden   []string `json:"hidden,omitempty"`   // no completion

	// Leaf command parameter possibilities (simple completion)
	Params []string `json:"params,omitempty"` // just completed

	// Initial variable values, keys will complete with buitin var
	// command. If nil no var command is added.
	Vars map[string]string `json:"vars,omitempty"`

	// Set to true to indicate can be compiled independently or used as
	// a base command in a multicall binary. When true and Vars is non-nil
	// cache for Vars is initialized if not present using the [Cmd.Name]
	IsBase bool `json:isbase, omitempty`

	// Disables the automatic help command. Useful for leaf commands where
	// help could be an actual value.
	NoHelp bool `json:nohelp,omitempty`

	// Functions to be used for the Fill command which is automatically
	// called on most string properties (ex: {{ exename }}).
	FuncMap template.FuncMap `json:"-"`

	// If unassigned will be assigned [DefComp] if either [Commands] or
	// [Params] is not empty.
	Comp Completer `json:"-"`

	// Where the work happens.
	Init   Method `json:"-"` // before Commands/Call, good for validation
	Call   Method `json:"-"` // optional if Commands
	Caller *Cmd   `json:"-"`

	// Pass bulk input efficiently (when args won't do)
	Input io.Reader

	// Faster than lots of "if" conditions in Call
	MinArgs int  `json:"minargs,omitempty"`
	MaxArgs int  `json:"maxargs,omitempty"`
	NumArgs int  `json:"numargs,omitempty"`
	NoArgs  bool `json:"noargs,omitempty"`
	MinParm int  `json:"minparm,omitempty"`
	MaxParm int  `json:"minparm,omitempty"`

	_aliases map[string]*Cmd // see cacheAliases called from Run->Seek->Resolve
}

func (x *Cmd) FUsage() string     { return x.Fill(x.Usage) }
func (x *Cmd) FShort() string     { return x.Fill(x.Short) }
func (x *Cmd) FLong() string      { return x.Fill(x.Long) }
func (x *Cmd) FVersion() string   { return x.Fill(x.Version) }
func (x *Cmd) FCopyright() string { return x.Fill(x.Copyright) }
func (x *Cmd) FLicense() string   { return x.Fill(x.License) }
func (x *Cmd) FSite() string      { return x.Fill(x.Site) }
func (x *Cmd) FSource() string    { return x.Fill(x.Source) }
func (x *Cmd) FIssues() string    { return x.Fill(x.Issues) }

// Completer specifies a struct with a Complete function that will
// complete the given bonzai.Cmd with the given arguments.
// The Complete function must never panic and always return at least an
// empty slice of strings. Not all completers require the passed
// command. By convention such the unused [Cmd] should use a single
// underscore as their name. In the less likely event that arguments are
// not to do the completion (only the Cmd itself) then an underscore
// should also be used in place of "args".
type Completer interface {
	Complete(x *Cmd, args ...string) []string
}

// Method defines the main code to execute for a command (Cmd). By
// convention the parameter list should be named "args" if there are
// args expected and underscore (_) if not. Methods must never write
// error output to anything but standard error and should almost always
// use the slog package to do so. By convention "caller" is usually
// named "x". If the caller is unused the argument name is underscore (_)
// by convention.
type Method func(caller *Cmd, args ...string) error

// Names returns the Name and any Aliases grouped such that the Name is
// always last. Any alias beginning with anything but a letter (L) is
// omitted.
func (x *Cmd) Names() []string {
	var names []string
	for _, a := range x.Aliases {
		if len(a) == 0 {
			continue
		}
		if unicode.IsLetter([]rune(a)[0]) {
			names = append(names, a)
		}
	}
	names = append(names, x.Name)
	return names
}

// UsageGroup uses Bonzai usage notation, a basic form of regular
// expressions, to describe the arguments allowed where each argument is
// a literal string (avoid spaces). The arguments are joined with bars
// (|) and wrapped with parentheses producing a regex group.  The min
// and max are then applied by adding the following regex decorations
// after the final parenthesis:
//
//   - min=1 max=1 (exactly one)
//     ?          - min=0 max=0 (none or many)
//   - - min=1 max=0 (one or more)
//     {min,}     - min>0 max=0 (min, no max)
//     {min,max}  - min>0 max>0 (min and max)
//     {,max}     - min=0 max>0 (max, no min)
//
// An empty args slice returns an empty string. If only one arg, then
// that arg is simply returned and min and max are ignored. Arguments
// that are empty strings are ignored. No transformation is done to the
// string itself (such as removing white space).
func UsageGroup(args []string, min, max int) string {
	args = filt.NotEmpty(args)
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	default:
		var dec string
		switch {
		case min == 1 && max == 1:
		case min == 0 && max == 0:
			dec = "?"
		case min == 1 && max == 0:
			dec = "+"
		case min > 1 && max == 0:
			dec = fmt.Sprintf("{%v,}", min)
		case min > 0 && max > 0:
			dec = fmt.Sprintf("{%v,%v}", min, max)
		case min == 0 && max > 1:
			dec = fmt.Sprintf("{,%v}", max)
		}
		return "(" + strings.Join(args, "|") + ")" + dec
	}
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

// Title returns [Name] and [FShort] combined (if exists). If the Name
// field of the commands is not defined will return a "{ERROR}".
func (x *Cmd) Title() string {
	if x.Name == "" {
		return "{ERROR: Name is empty}"
	}
	summary := x.FShort()
	switch {
	case len(summary) > 0:
		return x.Name + " - " + summary
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

	copyright := x.FCopyright()
	license := x.FLicense()
	version := x.FVersion()

	switch {

	case len(copyright) > 0 && len(license) == 0 && len(version) == 0:
		return x.Name + " " + copyright

	case len(copyright) > 0 && len(license) > 0 && len(version) > 0:
		return x.Name + " (" + version + ") " +
			copyright + "\nLicense " + license

	case len(copyright) > 0 && len(license) > 0:
		return x.Name + " " + copyright + "\nLicense " + license

	case len(copyright) > 0 && len(version) > 0:
		return x.Name + " (" + version + ") " + copyright

	case len(copyright) > 0:
		return x.Name + "\n" + copyright

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
// # Delegation and completion
//
// Since Run is the main execution entry point for all Bonzai command
// trees it is also responsible for handling completion (tab or
// otherwise). Therefore, all Run methods have two modes: delegation and
// completion (both are executions of the Bonzai binary command tree).
// Delegation is the default mode.
//
// Completion mode is triggered by the detection of the bash shell and
// the COMP_LINE environment variable triggering the self-completion
// (complete -C cmd cmd) that is unique to bash. All other shells can
// use the "help json" structured data create external completion
// scripts.
//
// When COMP_LINE is set, Run prints a list of possible completions to
// standard output by calling the [Completer.Complete] function of its
// [Comp] field or [DefComp] if [Commands] or [Params] are set. Each
// Cmd therefore manages its own completion and can draw from an
// ecosystem of Completers or assign its own. See the
// [core/comp] package for more.
//
// # Never panic
//
// All panics are trapped with [run.TrapPanic] which normally exits with 1 and
// outputs the main message. See [run.TrapPanic] for details.
func (x *Cmd) Run() {
	defer run.TrapPanic()

	// complete -C cmd cmd
	if line := os.Getenv("COMP_LINE"); len(line) > 0 && run.ShellIsBash() {
		var list []string

		// find the leaf command
		lineargs := run.ArgsFrom(line)
		cmd, args := x.Seek(lineargs[1:])

		// default completer or package aliases, always exits
		if cmd.Comp == nil && (len(cmd.Commands) > 0 || len(cmd.Params) > 0) {
			list = append(list, DefComp.Complete(cmd, args...)...)
			each.Println(list)
			run.Exit()
			return
		}

		// own completer, delegate
		each.Println(cmd.Comp.Complete(cmd, args...))
		run.Exit()
		return
	}

	// DELEGATION

	// seek should never fail to return something, but ...
	cmd, args := x.Seek(os.Args[1:])

	if cmd == nil {
		run.ExitError(IncorrectUsage{cmd})
		return
	}

	// initialize before delegation and Call
	if cmd.Init != nil {
		if err := cmd.Init(cmd, args...); err != nil {
			run.ExitError(err)
		}
	}

	// default to first Command if no Call defined
	if cmd.Call == nil {
		if len(cmd.Commands) > 0 {
			fcmd := cmd.Commands[0]
			if fcmd.Call == nil {
				run.ExitError(DefCmdReqCall{cmd})
				return
			}
			fcmd.Caller = cmd
			cmd = fcmd
		} else {
			run.ExitError(NoCallNoCommands{cmd})
			return
		}
	}

	switch {
	case len(args) > 0 && cmd.NoArgs:
		run.ExitError(TooManyArgs{len(args), 0})
		return
	case len(args) < cmd.MinArgs:
		run.ExitError(NotEnoughArgs{len(args), cmd.MinArgs})
		return
	case cmd.MaxArgs > 0 && len(args) > cmd.MaxArgs:
		run.ExitError(TooManyArgs{len(args), cmd.MaxArgs})
		return
	case cmd.NumArgs > 0 && len(args) != cmd.NumArgs:
		run.ExitError(WrongNumArgs{len(args), cmd.NumArgs})
		return
	}

	// delegate
	if cmd.Caller == nil {
		cmd.Caller = x
	}

	if err := cmd.Call(cmd, args...); err != nil {
		run.ExitError(err)
		return
	}
	run.Exit()
}

// Root returns the root [Cmd] from the current Path. This must always be
// calculated every time since any Cmd can change positions and
// pedigrees at any time at run time. Returns self if no PathCmds found.
func (x *Cmd) Root() *Cmd {
	cmds := x.PathCmds()
	if len(cmds) > 0 {
		return cmds[0].Caller
	}
	return x.Caller
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

// Resolve looks up a given Command by name or alias from Aliases
// (caching a lookup map of aliases in the process).
func (x *Cmd) Resolve(name string) *Cmd {

	if x.Commands == nil {
		return nil
	}

	for _, c := range x.Commands {
		if name == c.Name {
			return c
		}
	}

	if x._aliases == nil {
		x.cacheAliases()
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
// justification.  Hidden commands are not included. Aliases that begin
// with anything but a letter (L) are not included. Note that the order
// of the Commands is preserved (not necessarily alphabetic).
func (x *Cmd) UsageCmdTitles() string {
	var set []string
	var summaries []string
	for _, c := range x.Commands {
		if x.IsHidden(c.Name) {
			continue
		}
		set = append(set, strings.Join(c.Names(), "|"))
		summaries = append(summaries, c.FShort())
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
// Hidden commands. Hidden command typically begin with underscore.
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

// Seek checks the args for command names returning the deepest along
// with the remaining arguments. Typically the args passed are directly
// from the command line. Seek also sets the Caller on each Cmd found
// during resolution.
func (x *Cmd) Seek(args []string) (*Cmd, []string) {
	if args == nil || x.Commands == nil {
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
// used to invoke the composing binary. Also see Path, PathNames.
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

// Path returns a dotted notation of the PathNames including an initial
// dot (for root). This useful for associating configuration and other
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

// Log is currently short for log.Printf() but may be supplemented in
// the future to have more fine-grained control of logging.
func (x *Cmd) Log(format string, a ...any) {
	log.Printf(format, a...)
}

/*

// FIXME make use of builtin Vars

// Get is a shorter version of Z.Vars.Get(x.Path()+"."+key) for
// convenience.
func (x *Cmd) Get(key string) (string, error) {

	path := x.Path()
	if path != "." {
		path += "."
	}

	ptr := path + key

	// FIXME used to be global, now needs to be scoped per Exe
	v := vars.Get(ptr)
	if v != "" {
		return v, nil
	}

	if x.Vars != nil {
		v, _ = x.Vars[key]
		if v != "" {
			x.Set(key, v)
		}
		return v, nil
	}

	return "", nil
}

// Set is a shorter version of Z.Vars.Set(x.Path()+"."+key.val) for
// convenience. Logs the error Z.Vars is not defined (see UseVars).
func (x *Cmd) Set(key, val string) error {
	if Vars == nil {
		return UsesVars{x}
	}
	path := x.Path()
	if path != "." {
		path += "."
	}
	return Vars.Set(path+key, val)
}

// Del is a shorter version of Z.Vars.Del(x.Path()+"."+key.val) for
// convenience. Logs the error Z.Vars is not defined (see UseVars).
func (x *Cmd) Del(key string) error {
	if Vars == nil {
		return UsesVars{x}
	}
	path := x.Path()
	if path != "." {
		path += "."
	}
	Vars.Del(path + key)
	return nil
}
*/

// Fill fills out the passed [text/template] string using the [Cmd] instance
// as the data object source for the template. It is called by the Get*
// family of field accessors but can be called directly as well. Also
// see [bonzai.FuncMap] for list of predefined template functions.
// Filled versions of most fields are available as dynamic methods
// beginning with F (ex: [FShort]).
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

// UsageError returns [IncorrectUsage] for self.
func (x *Cmd) UsageError() error {
	return IncorrectUsage{x}
}

// ----------------------- builtin command: var -----------------------

/*

// TODO add automatic var command

var vars Map

func init() {
	dir, _ := os.UserCacheDir()
	vars = New()
	vars.Id = run.ExeName
	vars.Dir = dir
	vars.File = `vars`
}

//go:embed text/vars_help.md
var _vars_help string

var varsCmd = &Cmd{
	Name:    `var`,
	Short: help.S(_vars_help),
	Desc:    help.D(_vars_help),
	Commands: []*Cmd{
		getCmd, // default
		help.Cmd, initCmd, setCmd, fileCmd, dataCmd, editCmd, deleteCmd,
	},
}

//go:embed get.md
var getDoc string

var getCmd = &Z.Cmd{
	Name:        `get`,
	Short:     `print a cached variable with a new line`,
	Commands:    []*Z.Cmd{help.Cmd},
	Long: getDoc,
	NumArgs:     1,

	Call: func(x *Z.Cmd, args ...string) error {
		val, err := x.Caller.Caller.Get(args[0])
		if err != nil {
			return err
		}
		term.Print(val)
		return nil
	},
}

var setCmd = &Z.Cmd{
	Name:        `set`,
	Short:     `safely sets (persists) a cached variable`,
	Usage:       `(help|<name>) [<args>...]`,
	Long: setDoc,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,

	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) > 1 {
			val := strings.Join(args[1:], " ")
			if err := x.Caller.Caller.Set(args[0], val); err != nil {
				return err
			}
		}
		return getCmd.Call(x, args[0])
	},
}

//go:embed text/vars_set.md
var setDoc string

var fileCmd = &Z.Cmd{
	Name:     `file`,
	Aliases:  []string{"f"},
	Short:  `outputs full path to the cached vars file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		term.Print(vars.Path())
		return nil
	},
}

//go:embed text/vars_init.md
var initDoc string

var initCmd = &Z.Cmd{
	Name:        `init`,
	Aliases:     []string{"i"},
	Short:     `(re)initializes current variable cache`,
	Commands:    []*Z.Cmd{help.Cmd},
	UseVars:     true, // but fulfills at init() above
	Long: initDoc,
	Call: func(x *Z.Cmd, _ ...string) error {
		if term.IsInteractive() {
			r := term.Prompt(`Really initialize %v? (y/N) `, vars.DirPath())
			if r != "y" {
				return nil
			}
		}
		return Z.Vars.Init()
	},
}

//go:embed text/vars_data.md
var dataDoc string

var dataCmd = &Z.Cmd{
	Name:        `data`,
	Aliases:     []string{"d"},
	Short:     `outputs contents of the cached variables file`,
	Long: dataDoc,
	Commands:    []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, _ ...string) error {
		fmt.Print(vars.Data())
		return nil
	},
}

//go:embed text/vars_edit.md
var editDoc string

var editCmd = &Z.Cmd{
	Name:        `edit`,
	Short:     `edit variables file ({{execachedir "vars"}}) `,
	Long: editDoc,
	Aliases:     []string{"e"},
	Commands:    []*Z.Cmd{help.Cmd},
	Call:        func(x *Z.Cmd, _ ...string) error { return vars.Edit() },
}

var deleteCmd = &Z.Cmd{
	Name:        `delete`,
	Aliases:     []string{`d`, `del`, `unset`},
	Short:     `delete variable(s) from cache`,
	Usage:       `(help|<name>...)`,
	Commands:    []*Z.Cmd{help.Cmd},
	MinArgs:     1,
	Long: ` The {{aka}} command deletes the specified variable from cache.`,

	Call: func(x *Z.Cmd, args ...string) error {
		for _, i := range args {
			vars.Del(x.Caller.Caller.Path(i))
		}
		return nil
	},
}

// CallDummy is useful for testing when non-nil Call function needed.
var CallDummy = func(_ *Cmd, args ...string) error { return nil }
*/
