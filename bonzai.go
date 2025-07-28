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
	if DefaultPersister != nil {
		if err := DefaultPersister.Setup(); err != nil {
			panic(err)
		}
	}
}

// DefaultPersister for any [Cmd] that is not created with its own
// persistence using [Cmd].Persist. The [Cmd.Get] and [Cmd.Set] use this if
// Cmd does not have its own. If assigned, its [Persist] method is
// called during init of the bonzai package.
var DefaultPersister Persister

// Persister specifies anything that implements a persistence layer for
// storage and retrieval of key/value combinations.
//
// # Empty values
//
// Since all implementations must return and set strings, the empty
// string value is considered unset or non-existent. This is consistent
// with working with [pkg/os.Getenv]. Therefore, there is no Delete or
// Has equivalent since a Set("") works to delete a value and
// a len(Get())>0 is the same as Has.
//
// # Setup
//
// Set up an existing persistence store or create and initialize a new
// one if one does not yet exist. Never clears or deletes one that has
// been previously initialized (which is outside the scope of this
// interface). Usually this is called within an init() function after
// the other specific configurations of the driver have been set (much
// like database or other drivers). When called from init should usually
// prompt a panic since something has gone wrong during initialization
// and no attempt to run main should proceed, but this depends on the
// severity of the error, and that is up to the implementations to
// decide.
//
// # Get
//
// Retrieves a value for a specific key in a case-sensitive way or
// returns an empty string if not found.
//
// # Set
//
// Assigns a value for a given key. If the key did not exist, must
// create it. Callers can choose to check for the declaration of a key
// before calling Set, such as with [Cmd.Vars] and [Cmd.Get] and
// [Cmd.Set] (which themselves are not implementations of this interface
// although they use one internally).
type Persister interface {
	Setup() error          // setup existing or create (never clear)
	Get(key string) string // accessor, "" if non-existent
	Set(key, val string)   // mutator, "" to effectively delete
}

type Cmd struct {
	Name  string // ex: delete (required)
	Alias string // ex: rm|d|del (optional)
	Opts  string // ex: mon|wed|fri (optional)

	// Declaration of shareable variables (set to nil at runtime after
	// caching internal map, see [Cmd.VarsSlice], [Cmd.Get], [Cmd.Set]).
	Vars    Vars
	Persist Persister

	// Work down by this command itself
	Init func(x *Cmd, args ...string) error // initialization with [SeekInit]
	Do   func(x *Cmd, args ...string) error // main (optional if Def or Cmds)

	// Delegated work
	Cmds []*Cmd // composed subcommands (optional if Do or Def)
	Def  *Cmd   // default (optional if Do or Cmds, not required in Cmds)

	// Documentation
	Usage string           // text (<70 runes) (optional)
	Short string           // text (<50 runes) (optional)
	Long  string           // text/markup (optional)
	Vers  string           // text (<50 runes) (optional)
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
// contains a [pkg/sync.Mutex] allowing safe-for-concurrency modification
// when needed. This is used instead of [pkg/sync.RWMutex] because it is
// common for the first Get operation to also set the declared initial
// value such as when working with persistence.
//
// Var is rarely used directly being declared in [Cmd].Vars and used by
// [Cmd.Get] and [Cmd.Set] under the hood after being cached internally
// by the [Cmd.Run] method.
//
// # Empty string means undefined
//
// Bonzai variables are assumed to be undefined when empty (much like
// the use of traditional environment variables). This means that
// command creators must use something besides an empty string for
// booleans and such. There are no other specific string-to-type
// marshaling standards defined by Bonzai conventions, only that an
// empty string means "undefined".
//
// # Mandatory explicit declaration
//
// [Cmd.Run], [Cmd.Get], or [Cmd.Set] panics if the key referenced is
// not declared in the [Cmd].Vars slice.
//
// # Environment variables
//
// When an environment variable name (E) is provided and an initial
// value (V) is not set, the value of E is looked up is taken as the
// exact name of an environment variable that, if found to exist, even
// if blank, is used exactly as if the initial V value had been
// explicitly assigned that value returned from the environment
// variable value lookup. One use case is to allow users to decide if
// they prefer to manage initial values through environment variables
// instead of the alternatives. Note that when V is assigned indirectly
// in this way that any persistence is also applied in addition to the
// in-memory assignment of the value of V.
//
// When an environment variable name (E) is provided and an initial
// value (V) is also set, the environment variable, if found to exist,
// even if blank, is said to "shadow" the initial and current value (V)
// as well as any persistence. Get and Set operations are applied to the
// in-memory environment variable only. This is useful for use cases
// that temporarily alter behavior without changing default values or
// persistent values, such as when debugging, or creating "dry runs".
// This also provides a well-supported conventional alternative to
// getopts style options:
//
//	LANG=en greet
//	GOOS=arch go build
//
// # Persistence
//
// When persistence (P) is true then either internal persistence
// ([Cmd].Persist) or default persistence ([DefaultPersister]) is checked
// and used if available (not nil). If neither is available then it is
// as if P is false. In-memory values (V) are always kept in sync
// with those persisted depending on the method of persistence employed
// by the [Persister]. Both are always safe for concurrency.
//
// When persistence (P) is true but the initial value (V) is empty, then
// persistence is assumed to have the value or be okay with an undefined
// persisted value.
//
// When persistence (P) is true and the initial value (V) is set and
// [Cmd.Get] retrieves an empty value (after the delegated call to
// [Persister].Get), then the initial value (V) is passed to
// [Persister].Set persisting it for the next time.
//
// # Explicit inheritance
//
// When the inherits (I) field is set, then the variable is inherited
// from some command above it and if not found must panic. When found,
// the inheriting Var internally assigns a reference ([Var].X) to the
// [Cmd] in which the inherited Var was declared [Var].X. The
// variable-related operations [Cmd.Get] and [Cmd.Set] then directly
// operate on the Var pointed to by the inherited Var ref (X) instead of
// itself. Setting the inheritor field (I) in addition to any
// other field causes a panic.
//
// # Scope and depth of inheritance
//
// There is no concept of variable scope or depth. This means any
// subcommand at any depth level may declare that it wants to inherit
// and therefore retrieve and modify any declared variable above it. It
// is assumed that any developer composing a Bonzai command branch would
// look at such declarations in Vars to decide if it is safe to import
// and compose that subcommand into the compilation of the overall
// command. The requirement to declare all Vars that do inherit makes
// this knowledge explicitly visible for all Bonzai commands.
//
// # Global package variable injection
//
// If the [Var].G is assigned the reference to a string variable, which
// are usually expected to be package globals that can be used by any
// command at all in the package, then Get is called on the command when
// the [Cmd].Vars list is cached during the early SeekInit phase. This
// means that parent commands can trigger the fetching of those
// variables and setting their values early on before anything else.
// This saves considerable boilerplate for otherwise calling Get every
// time it is needed when---and only when---the value is not expected to
// change over the life of the program execution (e.g. API keys, account
// IDs, etc.).
type Var struct {
	sync.Mutex
	K string  `json:"k,omitempty"` // key
	V string  `json:"v,omitempty"` // value
	E string  `json:"e,omitempty"` // environment variable name
	S string  `json:"s,omitempty"` // short description
	P bool    `json:"p,omitempty"` // persistent
	I string  `json:"i,omitempty"` // inherits
	R bool    `json:"r,omitempty"` // required to be non-empty
	X *Cmd    `json:"-"`           // inherited from
	G *string `json:"-"`           // fetch first value into global
}

func (v Var) String() string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

// WithPersister overrides or adds [Cmd].Persist. The [Persister].Setup
// method is called and panics on error.
func (x Cmd) WithPersister(a Persister) *Cmd {
	if err := a.Setup(); err != nil {
		panic(err)
	}
	x.Persist = a
	return &x
}

// Get returns the value of [pkg/os.LookupEnv] if [Var].E was set and
// a corresponding environment variable was found overriding everything
// else in priority and shadowing any initially declared in-memory value
// or persisted value even if the found env variable is set to empty
// string. [Var].V is completely ignored in this case. The
// environment variable never changes the persisted value.
//
// Otherwise, if [Var].P is true
// attempts to look it up from either internal persistence set with
// [Cmd].Persist or the package [DefaultPersister] default if either is
// not nil. If a persister is available but returns an empty string then
// it is assumed the initial value has never been persisted and the
// in-memory cached value is returned and also persisted with [Cmd].Set
// ensuring that initial [Cmd].Vars declarations are persisted if they
// need to be for the first time. Note that none of this works until
// [Cmd.Run] is called which caches the Vars and sets up
// persistence. Panics if key was never declared. Locks the variable so
// safe for concurrency but persisters must implement their own
// file-level locking if shared between multiple processes.
func (x *Cmd) Get(key string) string {
	// declaration is mandatory
	v, has := x.vars[key]
	if !has {
		panic(`developer-error: not declared in Vars: ` + key)
	}

	// inherited, recurse
	if v.X != nil {
		return v.X.Get(key)
	}

	v.Lock()
	defer v.Unlock()

	// env var shadows everything, even if empty
	if len(v.E) > 0 {

		if val, has := os.LookupEnv(v.E); has {
			if len(v.V) == 0 {
				v.V = val
			}
			return val
		}
	}

	// local persister, usually setup in x.Init
	if v.P && x.Persist != nil {
		pv := x.Persist.Get(key)
		if len(pv) > 0 {
			v.V = pv
		} else {
			if len(v.V) > 0 {
				x.Persist.Set(key, v.V)
			}
		}
		return v.V

	}

	// package-wide default persister
	if v.P && DefaultPersister != nil {
		pv := DefaultPersister.Get(key)
		if len(pv) > 0 {
			v.V = pv
		} else {
			if len(v.V) > 0 {
				DefaultPersister.Set(key, v.V)
			}
		}
		return v.V
	}

	return ""
}

// Set assigns the value of the internal vars value for the given key. If
// the [Var].E is found to exist with [os.LookupEnv] then it is only set
// instead. This allows environment variables to shadow in-memory and
// persistent variables. If [Cmd].Persister is set, then attempts to
// persist using the internal persister created with [Cmd].Persister or
// the [DefaultPersister] default if either is not nil. Otherwise,
// assumes no persistence and only changes the in-memory value.
// Panics if key was not declared in [Cmd].Vars. Locks the Var in
// question so safe for concurrency.
func (x *Cmd) Set(key, value string) {

	// declaration is mandatory
	v, has := x.vars[key]
	if !has {
		panic(`developer-error: not declared in Vars: ` + key)
	}

	// inherited, recurse
	if v.X != nil {
		v.X.Set(key, value)
		return
	}

	v.Lock()
	defer v.Unlock()

	// env var shadows everything if found, even if empty
	if len(v.E) > 0 {
		if _, has := os.LookupEnv(v.E); has {
			os.Setenv(v.E, value)
			return
		}
	}
	v.V = value

	// set first persistence found if Persist
	if v.P {
		if x.Persist != nil {
			x.Persist.Set(key, value)
			return
		}
		if DefaultPersister != nil {
			DefaultPersister.Set(key, value)
			return
		}
	}
}

func (x *Cmd) lookvar(key string) *Var { return x.vars[key] }

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
// command. It is not set until [Cmd.Seek] or [Cmd.SeekInit] is called
// or indirectly by [Cmd.Run] or [Cmd.Exec]. Caller is set to itself if
// there is no caller (see [Cmd.IsRoot]).
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
// property set to true so that [Cmd.IsHidden] returns true. Use cases
// include convenient inclusion of leaf commands that are already
// available elsewhere (like help or var) and allowing deprecated
// commands to be supported but hidden in help output. See the
// [pkg/github.com/rwxrob/bonzai/mark/funcs] package for examples.
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
// Exec traps any panics unless the DEBUG environment variable is
// set (truthy).
//
// # Multicall
//
// Exec uses [pkg/os.Args][0] compared to the [Cmd].Name to resolve what to
// run enabling the use of multicall binaries with dashes in the name (a
// common design pattern used by other monolith multicalls such as Git and
// BusyBox/Alpine).
func (x *Cmd) Exec(args ...string) {
	log.SetFlags(0)
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
	c.resolveInheritedVars()
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
		return fmt.Errorf(`developer-error: Validate called with nil receiver`)

	case len(c.Short) > 0 && (len(c.Short) > 50 || !unicode.IsLower(rune(c.Short[0]))):
		return ErrInvalidShort{c}

	case len(c.Usage) > 70:
		return ErrInvalidUsage{c}

	case len(c.Vers) > 50:
		return ErrInvalidVers{c}

	case IsValidName != nil && !IsValidName(c.Name):
		return ErrInvalidName{c.Name}

	case c.Do == nil && len(c.Cmds) == 0 && c.Def == nil:
		return ErrUncallable{c}
	}
	for _, v := range c.Vars {
		if len(v.I) > 0 &&
			(len(v.K) > 0 || len(v.V) > 0 || len(v.E) > 0 || v.X != nil || v.P || v.G != nil) {
			return ErrBadVarInheritance{v}
		}
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
		return fmt.Errorf(`developer-error: Validate called with nil receiver`)

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

// String fulfills the [pkg/fmt.Stringer] interface for [pkg/fmt.Print]
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
// [Cmd.Validate] and [Cmd.Init] functions are called (if any).
// Returns early with nil values and the error if any Validate or Init
// function produces an error.
func (x *Cmd) SeekInit(args ...string) (*Cmd, []string, error) {
	if err := x.Validate(); err != nil {
		return x, args, err
	}
	x.cacheVars()
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

func (x *Cmd) resolveInheritedVars() {
	for _, v := range x.vars {
		if len(v.I) == 0 {
			continue
		}
		for cur := x.Caller(); cur != nil && !cur.IsRoot(); cur = cur.Caller() {
			curvar := cur.lookvar(v.I)
			if curvar == nil {
				continue
			}
			if curvar.K == v.I {
				v.X = cur
				if v.R {
					if len(cur.Get(v.I)) == 0 {
						panic(`required but not set: ` + v.I)
					}
				}
			}
		}
		if v.X == nil {
			panic(`failed to find inherited Var: ` + v.I)
		}
	}
}

func (x *Cmd) cacheVars() {
	x.vars = make(map[string]*Var, len(x.Vars))
	if x.Persist != nil {
		x.Persist.Setup()
	}
	for _, v := range x.Vars {
		if len(v.I) > 0 {
			x.vars[v.I] = &v
			// have to put off X assignment until after callers resolved
			continue
		}
		x.vars[v.K] = &v
		if v.R && len(x.Get(v.K)) == 0 {
			panic(`required variable not set: ` + v.K)
		}
		if v.G != nil {
			*(v.G) = x.Get(v.K)
		}
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
	return fmt.Sprintf(`developer-error: invalid name: %v`, e.Name)
}

// ErrUncallable indicates that a command requires a Do, one
// [Cmd].Cmds or a [Cmd].Def assigned.
type ErrUncallable struct {
	Cmd *Cmd
}

func (e ErrUncallable) Error() string {
	return fmt.Sprintf(
		`developer-error: Cmd (%v) requires Do, Def, or Cmds`, e.Cmd)
}

// ErrDoOrDef suggests that a command cannot have both Do and Def
// functions, providing details on the conflicting [Cmd].
type ErrDoOrDef struct {
	Cmd *Cmd
}

func (e ErrDoOrDef) Error() string {
	return fmt.Sprintf(
		`developer-error: Cmd.Do or Cmd.Def (never both) (%v)`, e.Cmd.Path())
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
		`developer-error: Cmd.Short (%v) length must be less than 50 runes and must begin with a lowercase letter`, e.Cmd)
}

// ErrInvalidUsage indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Usage] description.
type ErrInvalidUsage struct {
	Cmd *Cmd
}

func (e ErrInvalidUsage) Error() string {
	return fmt.Sprintf(
		`developer-error: Cmd.Usage (%v) length must be less than 70 runes`, e.Cmd)
}

// ErrInvalidVers indicates that the short description length exceeds 50
// characters, providing a reference to the [Cmd] and its [Vers] description.
type ErrInvalidVers struct {
	Cmd *Cmd
}

func (e ErrInvalidVers) Error() string {
	return fmt.Sprintf(
		`developer-error: Cmd.Vers (%v) length must be less than 50 runes`,
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

// ErrBadVarInheritance indicates a [Cmd].Vars [Var] that violated
// Bonzai variable rules regarding inheritance.
type ErrBadVarInheritance struct {
	Var Var
}

func (e ErrBadVarInheritance) Error() string {
	return fmt.Sprintf(
		`inherits requires no other fields be set: %v`, e.Var.I,
	)
}
