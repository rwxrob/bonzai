package bonzai

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"unicode"
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

	// Declaration of shareable variables (set to nil at runtime after
	// caching internal map, see [Cmd.VarsSlice], [Cmd.Get], [Cmd.Set]).
	Vars Vars

	// Work down by this command itself
	Init func(x *Cmd, args ...string) error // initialization with [SeekInit]
	Do   func(x *Cmd, args ...string) error // main (optional if Def or Cmds)

	// Delegated work
	Cmds []*Cmd // composed subcommands (optional if Do or Def)
	Def  *Cmd   // default (optional if Do or Cmds, not required in Cmds)

	// Documentation
	Vers  string           // text (<50 runes) (optional)
	Short string           // text (<50 runes) (optional)
	Long  string           // text/markup (optional)
	Funcs template.FuncMap // own template tags (optional)

	// Faster than "if" conditions in [Cmd.Do] (all optional)
	MinArgs  int    // min
	MaxArgs  int    // max
	NumArgs  int    // exact
	NoArgs   bool   // 0
	RegxArgs string // regx check each arg (document in Long)

	// Self-completion support: complete -C foo foo
	Comp Completer

	caller   *Cmd            // see [Caller],[Seek], delegation
	aliases  []string        // see [cacheAlias]
	opts     []string        // see [cacheOpts]
	hidden   bool            // see [AsHidden] and [IsHidden]
	cmdAlias map[string]*Cmd // see [cacheCmdAlias]
	vars     map[string]*Var // see [Vars]
}

type Vars []Var

func (vs Vars) String() string {
	buf, _ := json.Marshal(vs)
	return string(buf)
}

// Var contains information to be shared between [Cmd] instances and
// contains a [sync.Mutex] allowing safe-for-concurrency modification
// when needed. The Env contains the optional name of an environment
// variable to use instead if set, even if empty (see [Cmd.Get]). If
// Persist is set to true, persistence drivers used by callers can
// determine whether to persist or not. See [Cmd.Set].
type Var struct {
	sync.Mutex
	K       string `json:"k,omitempty"`
	V       string `json:"v,omitempty"`
	Env     string `json:"env,omitempty"`
	Short   string `json:"short,omitempty"`
	Persist bool   `json:"persist,omitempty"`
}

func (v Var) String() string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

// Get returns the value of [os.LookupEnv] if Env was set, otherwise,
// returns the current internal value of Vars[key] (even though Vars is
// always empty after [Cmd.SeekInit] caches the initial values). If key
// was not declared in [Cmd].Vars silently returns empty string.
func (x *Cmd) Get(key string) string {
	v, has := x.vars[key]
	if !has {
		return ""
	}
	if len(v.Env) > 0 {
		if val, has := os.LookupEnv(v.Env); has {
			return val
		}
	}
	return v.V
}

// Set sets the value of the internal Var value for the given key. If
// the Var.Env is found to exist with [os.LookupEnv] then it is also
// set. If key was not declared in [Cmd].Vars silenty returns empty
// string.
func (x *Cmd) Set(key, value string) {
	v, has := x.vars[key]
	if !has {
		return
	}
	v.V = value
	if len(v.Env) > 0 {

		if _, has := os.LookupEnv(v.Env); has {
			os.Setenv(v.Env, v.V)
		}
	}
}

// VarsSlice returns a slice with a copy of the current [Cmd].Vars. Use
// [Cmd.Get] and [Cmd.Set] to access the actual value.
func (x *Cmd) VarsSlice() Vars {
	list := make([]Var, len(x.vars))
	n := 0
	for _, v := range x.vars {
		list[n] = *v
		n++
	}
	return list
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

// Run seeks the leaf command in the arguments passed, validating all
// with [Cmd.SeekInit] and calls its [Cmd].Do method passing itself as
// the first argument along with any remaining arguments. Run is always
// called from [Cmd.Exec] but can be called directly from another
// command's Do method to enable powerful command composition and
// delegation at a high-level. Run returns an error if a command cannot
// be found or the command fails validation in any way.
func (x *Cmd) Run(args ...string) error {
	c, args, err := x.SeekInit(args...)
	if err != nil {
		return err
	}
	return c.call(args)
}

// Nothing is a no-op function that implements the command function signature
// for [Cmd]. It takes a command [*Cmd] and a variadic number of string
// arguments and always returns nil, indicating no operation is performed.
var Nothing = func(*Cmd, ...string) error { return nil }

// Validate checks the integrity of the command. It verifies that
// the command is not nil, validates the length and format of the
// Short and Vers fields, checks the validity of the command
// Name, and ensures that the command is callable by checking the
// associated function Do, Def, and subcommands. It returns an error
// if any validation check fails. Validate is automatically called for
// every command during the [Cmd.SeekInit] descent to the leaf command.
func (c *Cmd) Validate() error {
	switch {

	case c == nil:
		return fmt.Errorf(`developer error: Validate called with nil receiver`)

	case len(c.Short) > 0 && (len(c.Short) > 50 || !unicode.IsLower(rune(c.Short[0]))):
		return ErrInvalidShort{c}

	case len(c.Vers) > 50:
		return ErrInvalidVers{c}

	case IsValidName != nil && !IsValidName(c.Name):
		return ErrInvalidName{c.Name}

	case c.Do == nil && len(c.Cmds) == 0 && c.Def == nil:
		return ErrUncallable{c}
	}
	return nil
}

// ValidateArgs checks the validity of the provided args against the
// constraints defined in the command. It returns an error if the
// arguments do not meet the minimum or maximum requirements, do not
// match the expected number, or fail regex validation. ValidateArgs is
// automatically called from [Cmd.SeekInit] right before the leaf
// command and its arguments are returned.
func (c *Cmd) ValidateArgs(args ...string) error {
	switch {

	case c == nil:
		return fmt.Errorf(`developer error: Validate called with nil receiver`)

	case len(args) < c.MinArgs:
		return ErrNotEnoughArgs{Count: len(args), Min: c.MinArgs}

	case c.MaxArgs > 0 && len(args) > c.MaxArgs:
		return ErrTooManyArgs{Count: len(args), Max: c.MaxArgs}

	case c.NumArgs > 0 && len(args) != c.NumArgs:
		return ErrWrongNumArgs{Count: len(args), Num: c.NumArgs}

	case c.NoArgs && len(args) > 0:
		return ErrTooManyArgs{Count: len(args), Max: 0}

	case len(c.RegxArgs) > 0:
		regx, err := regexp.Compile(c.RegxArgs)
		if err != nil {
			return err
		}
		for n, arg := range args {
			if !regx.MatchString(arg) {
				return ErrInvalidArg{Exp: c.RegxArgs, Index: n}
			}
		}
	}
	return nil
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
		cmd, args := x.Seek(lineargs[1:]...)

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
	cmds := x.Path()
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

// Seek checks the args passed for command names returning the deepest
// along with the remaining arguments. Typically the args passed are
// directly derived from the command line. Seek also sets [Cmd.Caller]
// on each [Cmd] in the path. Seek is indirectly called by [Cmd.Run] and
// [Cmd.Exec]. See [pkg/github.com/rwxrob/bonzai/cmds/help] for
// a practical example of how and why a command might need to call Seek.
// Also see [Cmd.SeekInit] when environment variables and initialization
// functions are wanted as well. Seek returns its receiver and the same
// single empty string argument if there is only one argument and it is
// an empty string (a special case used for completion). Returns self
// and the arguments passed if both Cmds and the default command (Def)
// are nil. If Cmds is nil but Def is not, returns the default command
// and the same list of arguments passed.
func (x *Cmd) Seek(args ...string) (*Cmd, []string) {
	if (len(args) == 1 && args[0] == "") || (x.Cmds == nil && x.Def == nil) {
		return x, args
	}
	if x.Cmds == nil && x.Def != nil {
		x.Def.caller = x
		return x.Def, args
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

// SeekInit is the same as [Cmd.Seek] but Vars are cached and the
// [Cmd].Validate and [Cmd].Init functions are called (if any).
// Returns early with nil values and the error if any Validate or Init
// function produces an error.
func (x *Cmd) SeekInit(args ...string) (*Cmd, []string, error) {
	x.cacheVars()
	if err := x.Validate(); err != nil {
		return x, args, err
	}
	if x.Init != nil {
		if err := x.Init(x, args...); err != nil {
			return x, args, err
		}
	}
	if (len(args) == 1 && args[0] == "") || (x.Cmds == nil && x.Def == nil) {
		return x, args, nil
	}
	if x.Cmds == nil && x.Def != nil {
		x.Def.caller = x
		return x.Def.SeekInit(args...)
	}
	cur := x
	n := 0
	for ; n < len(args); n++ {
		next := cur.resolve(args[n])
		if next == nil {
			break
		}
		next.cacheVars()
		if err := next.Validate(); err != nil {
			return next, args[n:], err
		}
		if next.Init != nil {
			if err := next.Init(next, args[n:]...); err != nil {
				return next, args[n:], err
			}
		}
		next.caller = cur
		cur = next
	}
	err := cur.ValidateArgs(args[n:]...)
	return cur, args[n:], err
}

func (x *Cmd) cacheVars() {
	x.vars = make(map[string]*Var, len(x.Vars))
	for _, v := range x.Vars {
		x.vars[v.K] = &v
	}
	x.Vars = nil
}

// Path returns the path of commands used to arrive at this command.
// The path is determined by walking backward from current caller up
// rather than depending on anything from the command line
// used to invoke the composing binary.
func (x *Cmd) Path() []*Cmd {
	path := []*Cmd{x}
	for p := x.caller; p != nil; p = p.caller {
		path = append(path, p)
	}
	slices.Reverse(path)
	return path[1:]
}

// PathNames returns a slice of strings containing the names of all
// commands in the command path for the command. It retrieves
// the commands using [Cmd.Path] and constructs a slice with their
// respective Name fields before returning it.
func (x *Cmd) PathNames() []string {
	cmds := x.Path()
	names := make([]string, len(cmds))
	for i, c := range cmds {
		names[i] = c.Name
	}
	return names
}

// PathDashed returns a string representation of the command path for
// the command with each command name joined by a dash ('-').
// It utilizes the [Cmd.PathNames] method to obtain the names of the
// commands in the path.
func (x *Cmd) PathDashed() string {
	return strings.Join(x.PathNames(), `-`)
}

func (x *Cmd) walkDeep(level int, fn func(int, *Cmd) error, onError func(error)) {
	if x == nil {
		return
	}
	if err := fn(level, x); err != nil {
		onError(err)
	}
	sublevel := level + 1
	if len(x.Cmds) > 0 {
		for _, cmd := range x.Cmds {
			cmd.walkDeep(sublevel, fn, onError)
		}
	}
}

// WalkDeep recursively traverses the command tree starting from itself,
// applying the function (fn) to each [Cmd] within [Cmd].Cmds with its
// level. If an error occurs while executing fn, the onError function is
// called with the error. Note that [Cmd].Def is only included if it is
// also in the Cmds slice.
func (x *Cmd) WalkDeep(fn func(int, *Cmd) error, onError func(error)) {
	x.walkDeep(0, fn, onError)
}

type leveled struct {
	cmd   *Cmd
	level int
}

func (x *Cmd) walkWide(level int, fn func(int, *Cmd) error, onError func(error)) {
	if x == nil {
		return
	}
	queue := []leveled{leveled{x, level}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if err := fn(cur.level, cur.cmd); err != nil {
			onError(err)
		}
		sublevel := level + 1
		for _, cmd := range cur.cmd.Cmds {
			queue = append(queue, leveled{cmd, sublevel})
		}
	}
}

// WalkWide performs a breadth-first traversal (BFS) of the command tree
// starting from the command itself, applying the function (fn) with the
// current level to each [Cmd] in [Cmd].Cmds recursively. If an error
// occurs during the execution of fn, the onError function is called
// with the error. It uses a queue to process each command and its
// subcommands iteratively. Note that [Cmd].Def is only included if it
// is also in the Cmds slice.
func (x *Cmd) WalkWide(fn func(int, *Cmd) error, onError func(error)) {
	x.walkWide(0, fn, onError)
}

// ErrInvalidName indicates that the provided name for the command is invalid.
// It includes the invalid [Name] that caused the error.
type ErrInvalidName struct {
	Name string
}

func (e ErrInvalidName) Error() string {
	return fmt.Sprintf(`developer error: invalid name: %v`, e.Name)
}

// ErrUncallable indicates that a command requires a Do, one
// [Cmd].Cmds or a [Cmd].Def assigned.
type ErrUncallable struct {
	Cmd *Cmd
}

func (e ErrUncallable) Error() string {
	return fmt.Sprintf(
		`developer error: Cmd (%v) requires Do, Def, or Cmds`, e.Cmd)
}

// ErrDoOrDef suggests that a command cannot have both Do and Def
// functions, providing details on the conflicting [Cmd].
type ErrDoOrDef struct {
	Cmd *Cmd
}

func (e ErrDoOrDef) Error() string {
	return fmt.Sprintf(
		`developer error: Cmd.Do or Cmd.Def (never both) (%v)`, e.Cmd.Path())
}

// ErrNotEnoughArgs indicates that insufficient arguments were provided,
// describing the current [Count] and the minimum [Min] required.
type ErrNotEnoughArgs struct {
	Count int
	Min   int
}

func (e ErrNotEnoughArgs) Error() string {
	return fmt.Sprintf(
		`usage error: %v is not enough arguments, %v required`,
		e.Count, e.Min)
}

// ErrTooManyArgs signifies that too many arguments were provided,
// including the current [Count] and the maximum [Max] allowed.
type ErrTooManyArgs struct {
	Count int
	Max   int
}

func (e ErrTooManyArgs) Error() string {
	return fmt.Sprintf(
		`usage error: %v is too many arguments, %v maximum`,
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
		`usage error: requires exactly %v arguments (not %v)`,
		e.Num, e.Count)
}

// ErrInvalidShort indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Short] description.
type ErrInvalidShort struct {
	Cmd *Cmd
}

func (e ErrInvalidShort) Error() string {
	return fmt.Sprintf(
		`developer error: Cmd.Short (%v) length must be less than 50 runes and must begin with a lowercase letter`, e.Cmd)
}

// ErrInvalidVers indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Vers] description.
type ErrInvalidVers struct {
	Cmd *Cmd
}

func (e ErrInvalidVers) Error() string {
	return fmt.Sprintf(
		`developer error: Cmd.Vers (%v) length must be less than 50 runes`,
		e.Cmd,
	)
}

// ErrInvalidArg indicates that the arguments did not match
// a particular possible regular expression.
type ErrInvalidArg struct {
	Exp   string
	Index int
}

func (e ErrInvalidArg) Error() string {
	return fmt.Sprintf(
		`usage error: arg #%v must match: %v`,
		e.Index+1, e.Exp,
	)
}
