package bonzai

// Configurer specifies how to configure a Bonzai Cmd. Configurations
// must always be maintained in YAML (of which JSON is a subset) and be
// fully compatible with gopkg.in/yaml.v3.
//
// The location where the YAML is persisted is not specified by this
// interface. Implementations may use any persistence layer that
// guarantees atomic OverWrites and can be edited with a local system
// editor.
//
// The id may be anything that can be a YAML key, but use of Unicode
// Letter class runes is strongly recommended.  Usually, the id will be
// derived from the name of the executable binary or root Bonzai node
// command.
//
// Configuration data is designed to change less frequently than cache
// data (see Cacher) and is always updated by editing the entire YAML
// file (never through future Set methods). It is important that
// configuration data not be abused and unnecessarily bloated to remain
// performant.
//
// The Init and OverWrite methods are destructive.  The argument passed
// to OverWrite will be marshalled to YAML and completely replace the
// existing configuration in a way that must guarantee atomic,
// system-wide write safety. (Locking a file is insufficient alone.)
//
// Query implementations must trim any initial or trailing white space
// (usually just a single line return from yq, for example) in order to
// ensure that the resulting values for edge matches can be used without
// needing a trim.
type Configurer interface {
	Init() error              // must initialize a new configuration
	Data() string             // must return full YAML
	Print()                   // must print full YAML to os.Stdout
	Edit() error              // must open full YAML in local editor
	OverWrite(with any) error // safely replace all configuration
	Query(q string) string    // yq compatible query string
	QueryPrint(q string)      // prints result to os.Stdout
}

// CacheMap specifies how to persist (cache) simple string key/value
// data. Implementations of CacheMap can persist in different ways,
// files, network storage, or cloud databases, etc. Must log errors
// rather than panic (unavailable source, etc.)
type CacheMap interface {
	Init() error                 // initialize completely new cache
	Data() string                // k=v with \r and \n escaped in v
	Print()                      // (printed)
	Get(key string) string       // accessor
	Set(key, val string) error   // mutator
	Del(key string)              // destroyer
	OverWrite(with string) error // safely replace all cache
}

// Completer defines a function to complete the given leaf Command with
// the provided arguments, if any. Completer functions must never be
// passed a nil Command or nil as the args slice. See comp.Standard.
type Completer func(leaf Command, args ...string) []string

// Section is a section from the Other attribute.
type Section interface {
	GetTitle() string
	GetBody() string
}

// Command interface encapsulates the Z.Cmd implementation under the
// bonzai/z package enabling the use of the interface type when an
// interface is needed, for example, when implementing Completers to
// avoid cyclical import dependencies. For consistency, dynamic
// attributes like Title() have been given a GetTitle() variation as
// well.
type Command interface {
	GetName() string
	GetTitle() string
	GetAliases() []string
	GetSummary() string
	GetUsage() string
	GetVersion() string
	GetCopyright() string
	GetLicense() string
	GetDescription() string
	GetSite() string
	GetSource() string
	GetIssues() string
	GetCommands() []Command
	GetCommandNames() []string
	GetParams() []string
	GetHidden() []string
	GetOther() []Section
	GetOtherTitles() []string
	GetCompleter() Completer
	GetCaller() Command
	GetMinArgs() int
	GetMinParm() int
	GetMaxParm() int
	GetReqConf() bool
	GetReqVars() bool
	GetUsageFunc() UsageFunc
}

// UsageFunc allows dynamic creation of usage strings for interactive
// help and error messages. Every Z.Cmd has one as does the Z package
// itself (which defaults to InferredUsage). The Z package version is
// used with a Cmd has not assigned its own. UsageFunc should take the
// Command interface as the only argument, but it is acceptable for Cmd
// implementations to cast the Command passed to a specific Cmd to give
// that implementation access to rest of the Cmd symbol scope.
type UsageFunc func(x Command) string
