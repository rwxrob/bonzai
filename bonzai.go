package bonzai

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

// AllowsPanic if set will turn any panic into an exit with an error
// message.
var AllowPanic bool

func init() {
	val, exists := os.LookupEnv(`DEBUG`)
	if !exists {
		return
	}
	val = strings.ToLower(strings.TrimSpace(val))
	isTruthy := slices.Contains([]string{"t", "true", "on"}, val)
	if num, err := strconv.Atoi(val); (err == nil && num > 0) ||
		isTruthy {
		AllowPanic = true
	}
}

type Cmd struct {
	Name  string // ex: delete (required)
	Alias string // ex: rm|d|del (optional)
	Opts  string // ex: mon|wed|fri (optional)

	// Shareable variables also used for persistent initial values
	Vars map[string]Var

	// Own work (optional if Cmds or Def)
	Do func(x *Cmd, args ...string) error

	// Delegated work
	Def  *Cmd   // default Cmd (optional)
	Cmds []*Cmd // composed commands (optional if Do or Cmds or Def)

	// Documentation
	Vers  string           // text (<50 runes) (optional)
	Short string           // text (<50 runes) (optional)
	Long  string           // text/markup (optional)
	Funcs template.FuncMap // own template tags (optional)

	// Faster than "if" conditions in [Cmd.Do] (all optional)
	MinArgs   int    // min
	MaxArgs   int    // max
	NumArgs   int    // exact, doubles as NoArg (0)
	MatchArgs string // regular expression filter (document in Long)

	// Self-completion support: complete -C foo foo
	Comp Completer

	caller   *Cmd            // delegation
	aliases  []string        // see [cacheAlias]
	opts     []string        // see [cacheOpts]
	hidden   bool            // see [AsHidden] and [IsHidden]
	cmdAlias map[string]*Cmd // see [cacheCmdAlias]
}

// Var contains information to be shared between [Cmd] instances and
// contains a [sync.Mutex] allowing safe-for-concurrency modification
// when needed. The Str
type Var struct {
	sync.Mutex
	Key   string // same as that used in Vars (map[string]Var)
	Short string // short description of variable
	Str   string
	Int   int
	Bool  bool
	Any   any
}

// Completer specifies anything with Complete function based
// on the remaining arguments. The Complete method must never panic
// and always return at least an empty slice of strings. For completing
// data from a [Cmd] use [CmdCompleter] instead. Implementations must
// not examine anything from the command line itself depending entirely
// on the passed arguments instead (which are usually the remaining
// arguments from the command line). Implementations may and will often
// depend on external data sources to determine the possible completion
// values, for example, current host names, users, or data from a web
// API endpoint.
type Completer interface {
	Complete(args ...string) []string
}

// CmdCompleter is a specialized [Completer] that requires a [Cmd]. This
// is used for the following core completions:
//
//   - [pkg/github.com/rwxrob/bonzai/comp.Cmds]
//   - [pkg/github.com/rwxrob/bonzai/comp.Aliases]
//   - [pkg/github.com/rwxrob/bonzai/comp.CmdsAliases]
//   - [pkg/github.com/rwxrob/bonzai/comp.Opts]
//   - [pkg/github.com/rwxrob/bonzai/comp.CmdsOpts]
//   - [pkg/github.com/rwxrob/bonzai/comp.CmdsOptsAliases]
type CmdCompleter interface {
	Completer
	Cmd() *Cmd
	SetCmd(x *Cmd)
}

// Caller returns the internal reference to the parent/caller of this
// command. It is not set until [Cmd.Seek] is called or indirectly by
// [Cmd.Run] or [Cmd.Exec]. Caller is set to itself if there is no
// caller (see [Cmd.IsRoot]).
func (x Cmd) Caller() *Cmd { return x.caller }

// WithName sets the [Cmd].Name to name and returns a pointer to a copy
// of the updated command with the new name. This covers a very specific
// but important use case when a naming conflict exists between two
// different commands at the same level within the command tree, for example,
// a help command that displays help in the local web browser and
// another help command with a new name that sends it to the terminal.
func (x Cmd) WithName(name string) *Cmd {
	x.Name = name
	return &x
}

// AsHidden returns a copy of the [Cmd] with its internal hidden
// property set to true preventing it from appearing in [Cmd.CmdTreeString]
// and some [CmdCompleter] values. Use cases include convenient
// inclusion of leaf commands that are already available elsewhere (like
// help or var) and allowing deprecated commands to be supported but
// hidden.
func (x Cmd) AsHidden() *Cmd {
	x.hidden = true
	return &x
}

// Aliases returns the [Cmd].Alias value split into a slice with the
// [Cmd].Name added as the last item.
func (x *Cmd) Aliases() []string {
	var names []string
	for _, alias := range x.aliasSlice() {
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

// cacheCmdAlias splits the [Cmd].Alias for each [Cmd] in [Cmds] with
// its respective [Cmd.aliasSlice] and assigns them the private cache
// map. If [Cmds] is nil or Name is empty silently returns. This is
// primarily used for bash tab completion support in Run and use as
// a multicall binary. It is exported (instead of private like the cache
// map itself) so that a additional aliases can be created at run time
// and updated.
func (x *Cmd) cacheCmdAlias() {
	x.cmdAlias = map[string]*Cmd{}
	if x.Cmds == nil || len(x.Name) == 0 {
		return
	}
	for _, c := range x.Cmds {
		aliases := c.aliasSlice()
		if len(aliases) == 0 {
			continue
		}
		for _, a := range aliases {
			x.cmdAlias[a] = c
		}
	}
}

// cmdAliasMap returns a cached map of all aliases pointing to this
// command.
func (x *Cmd) cmdAliasMap() map[string]*Cmd {
	if x.cmdAlias == nil {
		x.cacheCmdAlias()
	}
	return x.cmdAlias
}

// cacheOpts updates the [opts] cache by splitting [Opts]. Remember
// to call this whenever dynamically altering the value at
// runtime.
func (x *Cmd) cacheOpts() {
	if len(x.Opts) > 0 {
		x.opts = strings.Split(x.Opts, `|`)
		return
	}
	x.opts = []string{}
}

// OptsSlice returns the [Cmd].Opts as a cached slice (derived from the
// delimited string).
func (x *Cmd) OptsSlice() []string {
	if x.opts == nil {
		x.cacheOpts()
	}
	return x.opts
}

// cacheAlias updates the aliases cache by splitting Alias
// and adding the Name to the end. Remember to call this whenever
// dynamically altering the value at runtime.
func (x *Cmd) cacheAlias() {
	if len(x.Alias) > 0 {
		x.aliases = strings.Split(x.Alias, `|`)
		return
	}
	x.aliases = []string{}
}

// aliasSlice updates the aliases internal cache created from the
// delimited [Cmd].Alias value and returns it as a slice.
func (x *Cmd) aliasSlice() []string {
	if x.aliases == nil {
		x.cacheAlias()
	}
	return x.aliases
}

// Exec is called from main function in the main package and never
// returns, always exiting with either 0 or 1 usually printing any error
// encountered. Exec calls [Cmd.Run] (which returns errors instead of
// exiting) after the following runtime considerations.
//
// # Self-completion
//
// Exec checks the COMP_LINE environment variable and if found assumes
// the program is being called in self-completion context by a user
// tapping the tab key one or more times. See the bash man page section
// "Programmable Completion" (although zsh can also be setup to enable
// bash self-completion). All Bonzai programs can, therefore, be set to
// complete themselves:
//
//	complete -C foo foo
//
// See [Completer] and [CmdCompleter] for more information about
// completion and the [pkg/github.com/rwxrob/bonzai/comp/completers]
// package for a growing collection of community maintained common
// completer implementations. Contributions always welcome.
//
// # Trapped panics
//
// Exec traps any panics with unless the DEBUG environment variable is
// set (truthy).
//
// # Multicall
//
// Exec uses [os.Args][0] compared to the [Cmd].Name to resolve what to
// run enabling the use of multicall binaries with dashes in the name (a
// common design pattern used by other monolith multicalls such as git and
// busybox).
func (x *Cmd) Exec(args ...string) {
	defer trapPanic()
	x.recurseIfMulti(args)
	x.detectCompletion()
	if len(args) == 0 {
		args = os.Args[1:]
	}
	if err := x.Run(args...); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

// trapPanic recovers from any panic and more gracefully displays the
// panic by logging it before exiting with a return value of 1.
var trapPanic = func() {
	if !AllowPanic {
		if r := recover(); r != nil {
			log.Println(r)
			os.Exit(1)
		}
	}
}

// Run seeks the leaf command in the arguments passed, validates it, and
// calls its [Cmd].Do method passing itself as the first argument
// along with any remaining arguments. Run is always called from
// [Cmd.Exec] but can be called directly from another command's Do
// method to enable powerful command composition and delegation at
// a high-level. Run returns an error if a command cannot be found or the
// command fails validation in any way.
func (x *Cmd) Run(args ...string) error {
	c, args := x.Seek(args)

	switch {
	case c == nil:
		return ErrIncorrectUsage{c}

	case len(x.Short) > 50:
		return ErrInvalidShort{x}

	case len(x.Vers) > 50:
		return ErrInvalidVers{x}

	case IsValidName != nil && !IsValidName(x.Name):
		return ErrInvalidName{x.Name}

	case x.Do == nil && x.Def == nil && len(x.Cmds) == 0:
		return ErrUncallable{x}

	case x.Do != nil && x.Def != nil:
		return ErrDoOrDef{x}

	case len(args) < x.MinArgs:
		return ErrNotEnoughArgs{Count: len(args), Min: x.MinArgs}

	case x.MaxArgs > 0 && len(args) > x.MaxArgs:
		return ErrTooManyArgs{Count: len(args), Max: x.MaxArgs}

	case x.NumArgs > 0 && len(args) != x.NumArgs:
		return ErrWrongNumArgs{Count: len(args), Num: x.NumArgs}

	}
	return c.call(args)
}

func (x *Cmd) call(args []string) error {
	if x.caller == nil {
		x.caller = x
	}
	if x.Do != nil {
		if err := x.Do(x, args...); err != nil {
			return err
		}
		return nil
	}
	if x.Def != nil {
		return x.Def.call(args)
	}
	return nil
}

// IsHidden returns true if [Cmd.AsHidden] was used to create.
func (x *Cmd) IsHidden() bool { return x.hidden }

func (x *Cmd) has(c *Cmd) bool {
	for _, this := range x.Cmds {
		if this.Name == c.Name {
			return true
		}
	}
	return false
}

// called as multicall binary
func (x *Cmd) recurseIfMulti(args []string) {
	name := filepath.Base(os.Args[0])
	if name == x.Name {
		return
	}
	if c := x.Can(name); c != nil {
		c.Exec(args...)
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
			c.Exec(args...)
			return
		}
	}
}

func argsFrom(line string) []string {
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

// complete -C foo foo (man bash, Programmable Completion)
func (x *Cmd) detectCompletion() {
	if line := os.Getenv("COMP_LINE"); len(line) > 0 {

		// find the leaf command
		lineargs := argsFrom(line)
		cmd, args := x.Seek(lineargs[1:])

		if cmd.Comp == nil {
			os.Exit(0)
		}

		// not sure we've completed the command name itself yet
		if len(args) == 0 {
			fmt.Println(cmd.Name)
			os.Exit(0)
		}

		if v, is := cmd.Comp.(CmdCompleter); is {
			v.SetCmd(cmd)
		}

		// own completer, delegate
		for _, completion := range cmd.Comp.Complete(args...) {
			fmt.Println(completion)
		}
		os.Exit(0)
	}
}

// String fulfills the [fmt.Stringer] interface for [fmt.Print]
// debugging and template inclusion by simply printing the [Cmd].Name.
func (x Cmd) String() string { return x.Name }

// Root returns the root [Cmd] traversing up its [Cmd.Caller] tree
// returning self if already root. The value returned will always return true
// of its [Cmd.IsRoot] is called.
func (x *Cmd) Root() *Cmd {
	cmds := x.CmdPath()
	if len(cmds) > 0 {
		return cmds[0].caller
	}
	return x.caller
}

// IsRoot determines if the command has no parent commands above it by
// checking if its [Cmd.Caller] is the same as itself returning true if so.
func (x *Cmd) IsRoot() bool { return x.caller == x }

// resolve looks up a given [Cmd] by name or alias from [Alias]
// (caching a lookup map of aliases in the process).
func (x *Cmd) resolve(name string) *Cmd {
	if x.Cmds == nil {
		return nil
	}
	for _, c := range x.Cmds {
		if name == c.Name {
			return c
		}
	}
	aliases := x.cmdAliasMap()
	if c, has := aliases[name]; has {
		return c
	}
	return nil
}

// Can returns the first pointer to a command from the [Cmd].Cmds list
// that has a matching Name or Alias. If more than one argument is passed
// calls itself recursively on each item in the list.
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
	aliases := x.cmdAliasMap() // to trigger cache if needed
	if c, has := aliases[name]; has {
		return c
	}
	return nil
}

// CmdNames returns the names of all [Cmd].Cmds.
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

// Seek checks the args passed for command names returning the deepest along
// with the remaining arguments. Typically the args passed are directly
// derived from the command line. Seek also sets the [Cmd.Caller] for each
// command found as it descends. Seek is indirectly called by [Cmd.Run]
// and [Cmd.Exec]. See [pkg/github.com/rwxrob/bonzai/cmds/help] for
// a practical example of how and why a command might need to call Seek.
func (x *Cmd) Seek(args []string) (*Cmd, []string) {
	if (len(args) == 1 && args[0] == "") || x.Cmds == nil {
		return x, args
	}
	cur := x
	n := 0
	for ; n < len(args); n++ {
		next := cur.resolve(args[n])
		if next == nil {
			break
		}
		next.caller = cur
		cur = next
	}
	return cur, args[n:]
}

// CmdPath returns the path of commands used to arrive at this command.
// The path is determined by walking backward from current caller up
// rather than depending on anything from the command line
// used to invoke the composing binary.
func (x *Cmd) CmdPath() []*Cmd {
	path := []*Cmd{x}
	for p := x.caller; p != nil; p = p.caller {
		path = append(path, p)
	}
	slices.Reverse(path)
	return path[1:]
}

// ErrInvalidName indicates that the provided name for the command is invalid.
// It includes the invalid [Name] that caused the error.
type ErrInvalidName struct {
	Name string
}

func (e ErrInvalidName) Error() string {
	return fmt.Sprintf(`invalid name: %v`, e.Name)
}

// ErrIncorrectUsage signifies that the command was used incorrectly,
// providing a reference to the [Cmd] that encountered the issue.
type ErrIncorrectUsage struct {
	Cmd *Cmd
}

func (e ErrIncorrectUsage) Error() string {
	return fmt.Sprintf(`incorrect usage for "%v" command`, e.Cmd.Name)
}

// ErrUncallable indicates that a command requires a Do or Def function
// for execution, providing a reference to the [Cmd] in question.
type ErrUncallable struct {
	Cmd *Cmd
}

func (e ErrUncallable) Error() string {
	return fmt.Sprintf(`Cmd requires Do or Def of Cmds: %v`, e.Cmd.Name)
}

// ErrDoOrDef suggests that a command cannot have both Do and Def
// functions, providing details on the conflicting [Cmd].
type ErrDoOrDef struct {
	Cmd *Cmd
}

func (e ErrDoOrDef) Error() string {
	return fmt.Sprintf(`Do or Def (not both): %v`, e.Cmd.Name)
}

// ErrNotEnoughArgs indicates that insufficient arguments were provided,
// describing the current [Count] and the minimum [Min] required.
type ErrNotEnoughArgs struct {
	Count int
	Min   int
}

func (e ErrNotEnoughArgs) Error() string {
	return fmt.Sprintf(`%v is not enough arguments, %v required`,
		e.Count, e.Min)
}

// ErrTooManyArgs signifies that too many arguments were provided,
// including the current [Count] and the maximum [Max] allowed.
type ErrTooManyArgs struct {
	Count int
	Max   int
}

func (e ErrTooManyArgs) Error() string {
	return fmt.Sprintf(`%v is too many arguments, %v maximum`,
		e.Count, e.Max)
}

// ErrWrongNumArgs indicates that the number of arguments does not match
// the expected count, showing the current [Count] and the required [Num].
type ErrWrongNumArgs struct {
	Count int
	Num   int
}

func (e ErrWrongNumArgs) Error() string {
	return fmt.Sprintf(
		`%v arguments, %v required`,
		e.Count, e.Num)
}

// ErrInvalidShort indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Short] description.
type ErrInvalidShort struct {
	Cmd *Cmd
}

func (e ErrInvalidShort) Error() string {
	return fmt.Sprintf(`Cmd.Short length >50 for %q: %q`, e.Cmd, e.Cmd.Short)
}

// ErrInvalidVers indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Vers] description.
type ErrInvalidVers struct {
	Cmd *Cmd
}

func (e ErrInvalidVers) Error() string {
	return fmt.Sprintf(`Cmd.Vers length >50 for %q: %q`, e.Cmd, e.Cmd.Vers)
}
