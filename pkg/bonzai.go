package bonzai

// Configurer specifies the package configuration driver interface. One
// (and only one) implementation will be assigned to Z.Conf. Every
// implementation *must* assign itself to Z.Conf on init so that
// Bonzai tree developers need only import the implementation package.
//
// Furthermore, implementations must be maintained in YAML (of which
// JSON is a subset) and be fully compatible with gopkg.in/yaml.v3.
//
// The location to which YAML is to be persisted is not specified by
// this interface. Most implementations will use a file in
// os.UserConfigDir. Implementations may use any persistence layer that
// guarantees atomic OverWrites and can be edited with a local system
// editor.
//
// The Id may be anything that can be a YAML key, but use of Unicode
// Letter class runes is strongly recommended. Usually, the id will be
// derived from the name of the executable binary or Cmd.Root Bonzai node
// command.
//
// Configuration data is designed to change less frequently than cache
// data (see Vars) and is always updated by editing the entire YAML file
// (never through future Set methods). It is important that
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
	Init() error                    // must initialize a new configuration
	SoftInit() error                // must init if not yet initialized
	Data() (string, error)          // must return full YAML
	Print() error                   // must print full YAML to os.Stdout
	Edit() error                    // must open full YAML in local editor
	OverWrite(with any) error       // safely replace all configuration
	Query(q string) (string, error) // yq compatible query string
	QueryPrint(q string) error      // prints result to os.Stdout
}

// Vars specifies the package persistent variables driver interface. All
// implementations must assign themselves to Z.Vars during init. One
// (and only one) persistent variable driver is allowed per executable.
//
// Implementations must persist (cache) simple string key/value
// variables Implementations of Vars can persist in different ways, but
// most will write to os.UserCacheDir.  Files, network storage, or cloud
// databases, etc. are all allowed and expected.  However, each must
// always present the data in a .key=val format with \r and \n escaped
// and the key never must contain an equal (=). (Equal signs in the
// value are ignored.) This is the fastest format to read and parse.
type Vars interface {
	Init() error                 // initialize completely new cache
	SoftInit() error             // initialize if not already initialized
	Data() string                // k=v with \r and \n escaped in v
	Print()                      // (printed)
	Get(key string) string       // accessor
	Set(key, val string) error   // mutator
	Del(key string) error        // destroyer
	OverWrite(with string) error // safely replace all cache
}

// Completer specifies a struct with a Complete function that will
// complete the given bonzai.Command with the given arguments.
// The Complete function must never panic and always return at least an
// empty slice of strings.
type Completer interface {
	Complete(x Command, args ...string) []string
}

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
	GetShortcutsMap() map[string][]string
	GetShortcuts() []string
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
	GetComp() Completer
	GetCaller() Command
	GetMinArgs() int
	GetMinParm() int
	GetMaxParm() int
	GetUseConf() bool
	GetUseVars() bool
}
