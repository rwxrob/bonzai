// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package Z

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/fn/each"
	"github.com/rwxrob/fn/maps"
	"github.com/rwxrob/fn/redu"
	"github.com/rwxrob/structs/qstack"
	"github.com/rwxrob/to"
)

type Cmd struct {

	// main documentation, use Get* for filled template
	Name        string    `json:"name,omitempty"`        // plain
	Aliases     []string  `json:"aliases,omitempty"`     // plain
	Shortcuts   ArgMap    `json:"shortcuts,omitempty"`   // plain
	Summary     string    `json:"summary,omitempty"`     // template
	Usage       string    `json:"usage,omitempty"`       // template
	Version     string    `json:"version,omitempty"`     // template
	Copyright   string    `json:"copyright,omitempty"`   // template
	License     string    `json:"license,omitempty"`     // template
	Description string    `json:"description,omitempty"` // template
	Other       []Section `json:"other,omitempty"`       // template

	// for use with text/template for run-time additions to main
	// documentation and embedded file output (ex: {{ exename }})
	Dynamic template.FuncMap `json:"-"`

	// administrative URLs
	Site   string `json:"site,omitempty"`   // template, https:// assumed
	Source string `json:"source,omitempty"` // template, usually git url
	Issues string `json:"issues,omitempty"` // template, https:// assumed

	// descending tree, completable
	Commands []*Cmd   `json:"commands,omitempty"`
	Params   []string `json:"params,omitempty"`
	Hidden   []string `json:"hidden,omitempty"`

	// default values for Var variables, also for initialization
	VarDefs VarVals `json:"vardefs,omitempty"`

	// standard or custom completer, usually of form compfoo.New()
	Comp bonzai.Completer `json:"-"`

	// where the work happens
	Caller *Cmd   `json:"-"`
	Call   Method `json:"-"`
	Init   Method `json:"-"` // before Commands/Call, good for validation

	// pass bulk input efficiently (when args won't do)
	Input io.Reader

	// faster than lots of "if" in Call
	MinArgs  int  `json:"-"` // minimum number of args required (including parms)
	MaxArgs  int  `json:"-"` // maximum number of args required (including parms)
	NumArgs  int  `json:"-"` // exact number of args required (including parms)
	NoArgs   bool `json:"-"` // must not have any args
	MinParm  int  `json:"-"` // minimum number of params required
	MaxParm  int  `json:"-"` // maximum number of params required
	UseConf  bool `json:"-"` // requires Z.Conf be assigned
	UseVars  bool `json:"-"` // requires Z.Var be assigned
	ConfVars bool `json:"-"` // vars default to conf values before VarDefs

	_aliases  map[string]*Cmd   // see cacheAliases called from Run->Seek->Resolve
	_sections map[string]string // see cacheSections called from Run
}

// Section contains the Other sections of a command. Composition
// notation (without Title and Body) is not only supported but
// encouraged for clarity when reading the source for documentation.
type Section struct {
	Title string
	Body  string
}

func (s Section) GetTitle() string { return s.Title }
func (s Section) GetBody() string  { return s.Body }

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
// a "{ERROR}". Fills template for Summary.
func (x *Cmd) Title() string {
	if x.Name == "" {
		return "{ERROR: Name is empty}"
	}
	summary := x.GetSummary()
	switch {
	case len(summary) > 0:
		return x.Name + " - " + summary
	default:
		return x.Name
	}
}

// GetLegal returns a single line with the combined values of the
// Name, Version, Copyright, and License. If Version is empty or nil an
// empty string is returned instead. GetLegal() is used by the
// version builtin command to aggregate all the version information into
// a single output.
func (x *Cmd) GetLegal() string {

	copyright := x.GetCopyright()
	license := x.GetLicense()
	version := x.GetVersion()

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

// OtherTitles returns just the ordered titles from Other.
func (x *Cmd) OtherTitles() []string { return maps.Keys(x._sections) }

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

func (x *Cmd) cacheSections() {
	x._sections = map[string]string{}
	if len(x.Other) == 0 {
		return
	}
	for _, s := range x.Other {
		x._sections[s.Title] = s.Body
	}
}

// Run method resolves Shortcuts and Aliases and seeks the leaf Cmd. It
// then calls the leaf's first-class Call function passing itself as the
// first argument along with any remaining command line arguments.  Run
// returns nothing because it usually exits the program. Normally, Run
// is called from within main() to convert the Cmd into an actual
// executable program.  Exiting can be controlled, however, by calling
// ExitOn or ExitOff (primarily for testing). Use Call instead of Run
// when delegation is needed. However, avoid tight-coupling that comes
// from delegation with Call when possible. Also, Call automatically
// assumes the proper number and type of arguments have already been
// checked (see MinArgs, MaxArgs, NumArgs, etc.) which is normally done
// by Run. Use a high-level branch pkg instead (which is idiomatic for
// good Bonzai branch development).
//
// Handling Completion
//
// Since Run is the main execution entry point for all Bonzai command
// trees it is also responsible for handling completion (tab or
// otherwise). Therefore, all Run methods have two modes: delegation and
// completion (both are executions of the Bonzai binary command tree).
// Delegation is the default mode. Completion mode is triggered by the
// detection of the COMP_LINE environment variable.
//
// COMP_LINE
//
// When COMP_LINE is set, Run prints a list of possible completions to
// standard output by calling its Comp.Complete function
// (default Z.Comp). Each Cmd therefore manages its own completion and
// can draw from a rich ecosystem of Completers or assign its own custom
// one. This enables very powerful completions including dynamic
// completions that query the network or the local execution
// environment. Since Go can run on pretty much every device
// architecture right now, that's a lot of possibilities.  Even
// a line-based calculator can be implemented as a Completer. AI
// completers are also fully supported by this approach.  Intelligent
// completion eliminates the need for overly complex and error-prone
// (getopt) argument signatures for all Bonzai commands.
//
// Why COMP_LINE?
//
// Setting COMP_LINE has been a bash shell standard for more than a few
// decades. (Unfortunately, zsh dubiously chose to not support it for no
// good reason.) COMP_LINE completion, therefore, is the only planned
// method of detecting completion context. Enabling it in bash for any
// command becomes a simple matter of "complete -C foo foo" (rather than
// forcing users to evaluate thousands of lines of shell code to enable
// completion for even minimally complex command trees as other
// "commanders" require). Any code will work that sets COMP_LINE before
// calling Cmd.Run and receives a list of lines to standard
// output with completion candidates.
func (x *Cmd) Run() {
	defer TrapPanic()

	// returns throughout are because Exit* can be disabled
	// see ExitOff/On

	x.cacheSections()

	// COMPLETION

	line := os.Getenv("COMP_LINE")
	if line != "" {
		var list []string

		// find the leaf command
		lineargs := ArgsFrom(line)
		cmd, args := x.Seek(lineargs[1:])

		// default completer or package aliases, always exits
		if cmd.Comp == nil {
			if Comp != nil {
				list = append(list, Comp.Complete(cmd, args...)...)
			}
			each.Println(list)
			Exit()
			return
		}

		// own completer, delegate
		each.Println(cmd.Comp.Complete(cmd, args...))
		Exit()
		return
	}

	// DELEGATION

	// seek should never fail to return something, but ...
	cmd, args := x.Seek(os.Args[1:])

	if cmd == nil {
		ExitError(IncorrectUsage{cmd})
		return
	}

	// initialize before delegation and Call
	if cmd.Init != nil {
		if err := cmd.Init(cmd, args...); err != nil {
			ExitError(err)
		}
	}

	// default to first Command if no Call defined
	if cmd.Call == nil {
		if len(cmd.Commands) > 0 {
			fcmd := cmd.Commands[0]
			if fcmd.Call == nil {
				ExitError(DefCmdReqCall{cmd})
				return
			}
			fcmd.Caller = cmd
			cmd = fcmd
		} else {
			ExitError(NoCallNoCommands{cmd})
			return
		}
	}

	switch {
	case len(args) > 0 && cmd.NoArgs:
		ExitError(TooManyArgs{len(args), 0})
		return
	case len(args) < cmd.MinArgs:
		ExitError(NotEnoughArgs{len(args), cmd.MinArgs})
		return
	case cmd.MaxArgs > 0 && len(args) > cmd.MaxArgs:
		ExitError(TooManyArgs{len(args), cmd.MaxArgs})
		return
	case cmd.NumArgs > 0 && len(args) != cmd.NumArgs:
		ExitError(WrongNumArgs{len(args), cmd.NumArgs})
		return
	case cmd.UseConf && Conf == nil:
		ExitError(UsesConf{x})
		return
	case cmd.UseVars && Vars == nil:
		ExitError(UsesVars{x})
		return
	}

	// delegate
	if cmd.Caller == nil {
		cmd.Caller = x
	}

	if err := cmd.Call(cmd, args...); err != nil {
		ExitError(err)
		return
	}
	Exit()
}

// Root returns the root Cmd from the current Path. This must always be
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
		summaries = append(summaries, c.GetSummary())
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

// UsageCmdShortcuts returns a single string with the Shortcuts
// indented and with a maximum title signature length for
// justification.
func (x *Cmd) UsageCmdShortcuts() string {
	var set []string
	var summaries []string
	for k, v := range x.Shortcuts {
		set = append(set, k)
		summaries = append(summaries, strings.Join(v, " "))
	}
	longest := redu.Longest(set)
	var buf string
	for n := 0; n < len(set); n++ {
		if len(summaries[n]) > 0 {
			buf += fmt.Sprintf(
				`%-`+strconv.Itoa(longest)+"v - %v\n", set[n], summaries[n])
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

// Seek checks the args for command names returning the deepest along
// with the remaining arguments. Typically the args passed are directly
// from the command line. Seek also sets the Caller on each Cmd found
// during resolution. Cmd.Shortcuts are expanded.
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
	if len(cur.Shortcuts) > 0 && len(args[n:]) > 0 {
		if v, has := cur.Shortcuts[args[n]]; has {
			nargs := v
			nargs = append(nargs, args[n+1:]...)
			return cur.Seek(nargs)
		}
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
// dot (for root). This is compatible yq query expressions and useful
// for associating configuration and other data specifically with this
// command. If any arguments are passed then will be added with dots
// between them.
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

// C is a shorter version of
// strings.TrimSpace(Z.Conf.Query(x.Path()+"."+q)) for convenience. The
// yq YAML/JSON queries unreliably add line returns sometimes and other
// times not. If a line return is wanted, the caller will need to use
// fmt.Println. Also see UseConf.
func (x *Cmd) C(q string) (string, error) {
	if Conf == nil {
		return "", UsesConf{x}
	}
	path := x.Path()
	if path != "." {
		path += "."
	}
	res, err := Conf.Query(path + q)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(res), nil
}

// Get is a shorter version of Z.Vars.Get(x.Path()+"."+key) for
// convenience but also observes VarDefsFromConf if set. Also see
// UseVars.
func (x *Cmd) Get(key string) (string, error) {

	if Vars == nil {
		return "", UsesVars{x}
	}

	path := x.Path()
	if path != "." {
		path += "."
	}

	ptr := path + key

	v := Vars.Get(ptr)
	if v != "" {
		return v, nil
	}

	if x.ConfVars {
		v, _ = Conf.Query(ptr)
		if v != "" && v != "null" {
			x.Set(key, v)
			return v, nil
		}
	}

	if x.VarDefs != nil {
		v, _ = x.VarDefs[key]
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

// Fill fills out the passed text/template string using the Cmd instance
// as the data object source for the template. It is called by the Get*
// family of field accessors but can be called directly as well. Also
// see markfunc.go for list of predefined template functions.
func (x *Cmd) Fill(tmpl string) string {
	funcs := to.MergedMaps(Dynamic, x.Dynamic)
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

// UsageError returns IncorrectUsage for self.
func (x *Cmd) UsageError() error {
	return IncorrectUsage{x}
}

// --------------------- bonzai.Command interface ---------------------

// GetName fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetName() string { return x.Name }

// GetTitle fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetTitle() string { return x.Title() }

// GetAliases fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetAliases() []string { return x.Aliases }

// GetShortcutsMap fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetShortcutsMap() map[string][]string {
	return map[string][]string(x.Shortcuts)
}

// GetShortcuts fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetShortcuts() []string { return x.Shortcuts.Keys() }

// GetSummary fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetSummary() string { return x.Fill(x.Summary) }

// GetUsage fulfills the bonzai.Command interface. Uses x.Usage if not
// empty. Otherwise, calls UsageFunc to get usage, then uses Fill.
func (x *Cmd) GetUsage() string {
	usage := x.Usage
	if usage == "" {
		usage = UsageFunc(x)
	}
	return x.Fill(usage)
}

// GetVersion fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetVersion() string { return x.Fill(x.Version) }

// GetCopyright fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetCopyright() string { return x.Fill(x.Copyright) }

// GetLicense fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetLicense() string { return x.Fill(x.License) }

// GetDescription fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetDescription() string { return x.Fill(x.Description) }

// GetSite fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetSite() string { return x.Fill(x.Site) }

// GetSource fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetSource() string { return x.Fill(x.Source) }

// GetIssues fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetIssues() string { return x.Fill(x.Issues) }

// GetMinArgs fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetMinArgs() int { return x.MinArgs }

// GetMinParm fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetMinParm() int { return x.MinParm }

// GetMaxParm fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetMaxParm() int { return x.MaxParm }

// GetUseConf fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetUseConf() bool { return x.UseConf }

// GetUseVars fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetUseVars() bool { return x.UseVars }

// GetCommands fulfills the bonzai.Command interface.
func (x *Cmd) GetCommands() []bonzai.Command {
	var commands []bonzai.Command
	for _, s := range x.Commands {
		commands = append(commands, bonzai.Command(s))
	}
	return commands
}

// GetCommandNames fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetCommandNames() []string { return x.CmdNames() }

// GetHidden fulfills the bonzai.Command interface.
func (x *Cmd) GetHidden() []string { return x.Hidden }

// GetParams fulfills the bonzai.Command interface.
func (x *Cmd) GetParams() []string { return x.Params }

// GetOther fulfills the bonzai.Command interface. Uses Fill.
func (x *Cmd) GetOther() []bonzai.Section {
	var sections []bonzai.Section
	for _, s := range x.Other {
		s.Body = x.Fill(s.Body)
		sections = append(sections, bonzai.Section(s))
	}
	return sections
}

// GetOtherTitles fulfills the bonzai.Command interface. No Fill.
func (x *Cmd) GetOtherTitles() []string {
	var titles []string
	for _, title := range x.OtherTitles() {
		titles = append(titles, title)
	}
	return titles
}

// GetComp fulfills the Command interface.
func (x *Cmd) GetComp() bonzai.Completer { return x.Comp }

// GetCaller fulfills the bonzai.Command interface.
func (x *Cmd) GetCaller() bonzai.Command { return x.Caller }
