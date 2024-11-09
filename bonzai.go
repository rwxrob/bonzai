package bonzai

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/rwxrob/bonzai/run"
)

// Completer specifies a struct with a [Completer.Complete] function
// that will complete the first argument (usually a command of some kind)
// based on the remaining arguments. The [Complete] method must never
// panic and always return at least an empty slice of strings. By
// convention Completers that do not make use of or other arguments
// should use an underscore identifier since they are ignored.
type Completer interface {
	Complete(x Cmd, args ...string) []string
}

type Cmd struct {
	Name  string // ex: delete
	Alias string // ex: rm|d|del
	Opts  string // ex: mon|wed|fri

	// Work done by this command
	Call Method // if nil, Def must be set

	// Delegation to subcommands
	Def  *Cmd   // default [Cmd] if no Call and no matching Cmds
	Cmds []*Cmd // compiled/composed commands in no particular order
	Hide string // disable completion: ex: old|defunct

	// Minimal when [Cmd.Docs] is overkill
	Usage string // text
	Vers  string // text (<50 runes)
	Short string // text (<50 runes)
	Long  string // text/markup

	// Faster than lots of "if" conditions in [Cmd.Call]. Consider
	// [Cmd.Init] when more complex argument validation is needed.
	MinArgs   int    // min
	MaxArgs   int    // max
	NumArgs   int    // exact, also used for NoArg (0)
	MatchArgs string // PCRE/Go regular expression (requires Usage)

	// Self-completion support: complete -C foo foo
	Comp Completer

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

// WithName sets the [Name] of the command [x] to the specified [name]
// and returns a pointer to a copy of the updated command.
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

// IsValidName is assigned a function that returns a boolean for the
// given name. Note that if this is changed certain characters may break
// the creation of multicall binary links and bash completion. See the
// [pkg/github.com/rwxrob/bonzai/is] package for alternatives.
var IsValidName = allLatinASCIILowerWithDashes

func allLatinASCIILowerWithDashes(in string) bool {
	if len(in) == 0 || in[0] == '-' || in[len(in)-1] == '-' {
		return false
	}
	for _, r := range in {
		if ('a' <= r && r <= 'z') || r == '-' {
			continue
		}
		return false
	}
	return true
}

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
				run.ExitError(ErrInvalidName{alias})
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
// Throws [ErrInvalidName] error and exist if [Name] does not pass
// [ErrInvalidName] check.
//
// # Subcommand optional arguments
//
// If any argument is detected, delegation through recursive Run calls
// to subcommands is attempted. If more than one argument, each
// argument is assumed to be a [Name] or alias from [Alias] and so on
// (see [Can] for details).
func (x *Cmd) Run(args ...string) {
	defer run.TrapPanic()
	x.exitUnlessValidName()
	x.exitUnlessValidShort()
	x.recurseIfMulti(args)
	x.detectCompletion(args)
	if len(args) == 0 {
		args = os.Args[1:]
	}
	c, args := x.Seek(args)
	if c == nil {
		run.ExitError(ErrIncorrectUsage{c})
		return
	}
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
		run.ExitError(ErrUncallable{x})
		return
	case x.Call != nil && x.Def != nil:
		run.ExitError(ErrCallOrDef{x})
		return
	}
}

func (x *Cmd) exitIfBadArgs(args []string) {
	switch {
	case len(args) < x.MinArgs:
		run.ExitError(ErrNotEnoughArgs{len(args), x.MinArgs})
		return
	case x.MaxArgs > 0 && len(args) > x.MaxArgs:
		run.ExitError(ErrTooManyArgs{len(args), x.MaxArgs})
		return
	case x.NumArgs > 0 && len(args) != x.NumArgs:
		run.ExitError(ErrWrongNumArgs{len(args), x.NumArgs})
		return
	}
}

func (x *Cmd) exitUnlessValidName() {
	if !IsValidName(x.Name) {
		run.ExitError(ErrInvalidName{x.Name})
		return
	}
}

func (x *Cmd) exitUnlessValidShort() {
	if len(x.Short) > 50 {
		run.ExitError(ErrInvalidShort{x})
		return
	}
}

// called as multicall binary
func (x *Cmd) recurseIfMulti(args []string) {
	name := run.ExeName()
	if name == x.Name {
		return
	}
	if c := x.Can(name); c != nil {
		c.Run(args...)
		return
	}
	if strings.Contains(name, `-`) {
		fields := strings.Split(name, `-`)
		if fields[0] == x.Name {
			fields = fields[1:]
		}
		if c := x.Can(fields[0]); c != nil {
			if len(fields) > 1 {
				args = append(fields[1:], args...)
			}
			c.Run(args...)
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
		for _, completion := range cmd.Comp.Complete(*cmd, args...) {
			fmt.Println(completion)
		}
		run.Exit()
		return
	}
}

// String fulfills the [fmt.Stringer] interface for debugging by simply
// printing the Name of the [Cmd].
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

// IsRoot determines if the command [x] is the root command by checking
// if its [Caller] is the same as itself. It returns true if [x] is
// the root command; otherwise, it returns false.
func (x *Cmd) IsRoot() bool { return x.Caller == x }

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
// from the command line. Seek also sets the [Caller] on each [Cmd] found
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
	path := []*Cmd{x}
	for p := x.Caller; p != nil; p = p.Caller {
		path = append(path, p)
	}
	slices.Reverse(path)
	return path[1:]
}

// PathNames returns the path of command names used to arrive at this
// command. The path is determined by walking backward from current
// Caller up rather than depending on anything from the command line
// used to invoke the composing binary. Also see Path.
func (x *Cmd) PathNames() []string {
	names := []string{x.Name}
	p := x.Caller
	for p != nil {
		names = append(names, p.Name)
		if p == p.Caller {
			break
		}
		p = p.Caller
	}
	slices.Reverse(names)
	return names[1:]
}

// Title generates a formatted string representing the command title,
// including its Name, Alias, and Short description.
// If Name is empty, it defaults to "NONAME".
func (m Cmd) Title() string {
	out := new(strings.Builder)
	if len(m.Name) == 0 {
		m.Name = `NONAME`
	}
	out.WriteString(m.Name)
	if len(m.Alias) > 0 {
		out.WriteString(" (" + m.Alias + ")")
	}
	if len(m.Short) > 0 {
		out.WriteString(" - " + m.Short)
	}
	return out.String()
}

// Mark outputs raw, unfilled, BonzaiMark for rendering with
// a [pkg/github.com/rwxrob/bonzai/mark.Renderer] as the third argument.
func (m *Cmd) Mark() io.Reader {
	out := new(strings.Builder)
	out.WriteString("# Name\n\n")
	out.WriteString(m.Title() + "\n\n")
	out.WriteString("# Synopsis\n\n")
	out.WriteString(m.CmdTree() + "\n")
	if len(m.Long) > 0 {
		out.WriteString(dedent(m.Long))
	}
	return strings.NewReader(out.String())
}

var isblank = regexp.MustCompile(`^\s*$`)

func indentation(in string) int {
	var n int
	var v rune
	for n, v = range []rune(in) {
		if !unicode.IsSpace(v) {
			break
		}
	}
	return n
}

func dedent(in string) string {
	lines := strings.Split(in, "\n")
	for len(lines) == 1 && isblank.MatchString(lines[0]) {
		return ""
	}
	var n int
	for len(lines[n]) == 0 || isblank.MatchString(lines[n]) {
		n++
	}
	starts := n
	indent := indentation(lines[n])
	for ; n < len(lines); n++ {
		if len(lines[n]) >= indent {
			lines[n] = lines[n][indent:]
		}
	}
	return strings.Join(lines[starts:], "\n")
}

func (m *Cmd) cmdTree(depth int) string {
	out := new(strings.Builder)
	for range depth {
		out.WriteString("  ")
	}
	out.WriteString(m.Title() + "\n")
	depth++
	for _, c := range m.Cmds {
		out.WriteString(c.cmdTree(depth))
	}
	return out.String()
}

// CmdTree generates and returns a formatted string representation of
// the command tree for the [Cmd] instance. It aligns dashes in the
// output for better readability, adjusting spaces based on the position
// of the dashes.
func (m *Cmd) CmdTree() string {
	lines := strings.Split(m.cmdTree(1), "\n")
	dashindex := make([]int, len(lines))
	var dashcol int
	for i, line := range lines {
		n := strings.Index(line, `-`)
		dashindex[i] = n
		if n > dashcol {
			dashcol = n
		}
	}
	for i, line := range lines {
		n := dashindex[i]
		numspace := dashcol - n
		spaces := new(strings.Builder)
		for range numspace {
			spaces.WriteString(` `)
		}
		if n > 0 {
			lines[i] = line[:n] + spaces.String() + line[n:]
		}
	}
	return strings.Join(lines[1:], "\n")
}

// -------------------------- ErrInvalidName --------------------------

type ErrInvalidName struct {
	Name string
}

func (e ErrInvalidName) Error() string {
	return fmt.Sprintf(`invalid name: %v`, e.Name)
}

// -------------------------- ErrIncorrectUsage --------------------------

type ErrIncorrectUsage struct {
	Cmd *Cmd
}

func (e ErrIncorrectUsage) Error() string {
	if len(e.Cmd.Usage) == 0 {
		return fmt.Sprintf(`incorrect usage for "%v" command`, e.Cmd.Name)
	}
	return fmt.Sprintf(`usage: %v %v`,
		e.Cmd.Name,
		e.Cmd.Usage,
	)
}

// ---------------------------- ErrUncallable ----------------------------

type ErrUncallable struct {
	Cmd *Cmd
}

func (e ErrUncallable) Error() string {
	return fmt.Sprintf(`Cmd requires Call or Def: %v`, e.Cmd.Name)
}

// ----------------------------- ErrCallOrDef ----------------------------

type ErrCallOrDef struct {
	Cmd *Cmd
}

func (e ErrCallOrDef) Error() string {
	return fmt.Sprintf(`Call or Def (not both): %v`, e.Cmd.Name)
}

// --------------------------- ErrNotEnoughArgs --------------------------

type ErrNotEnoughArgs struct {
	Count int
	Min   int
}

func (e ErrNotEnoughArgs) Error() string {
	return fmt.Sprintf(`%v is not enough arguments, %v required`,
		e.Count, e.Min)
}

// ---------------------------- ErrTooManyArgs ---------------------------

type ErrTooManyArgs struct {
	Count int
	Max   int
}

func (e ErrTooManyArgs) Error() string {
	return fmt.Sprintf(`%v is too many arguments, %v maximum`,
		e.Count, e.Max)
}

// --------------------------- ErrWrongNumArgs ---------------------------

type ErrWrongNumArgs struct {
	Count int
	Num   int
}

func (e ErrWrongNumArgs) Error() string {
	return fmt.Sprintf(
		`%v arguments, %v required`,
		e.Count, e.Num)
}

// -------------------------- ErrInvalidShort -------------------------

type ErrInvalidShort struct {
	Cmd *Cmd
}

func (e ErrInvalidShort) Error() string {
	return fmt.Sprintf(`short length >50 %v: %v`, e.Cmd, e.Cmd.Short)
}
