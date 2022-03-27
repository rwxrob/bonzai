package conf

// Configurer specifies how to configure a Bonzai Cmd. Configurations
// must always be maintained in YAML (of which JSON is a subset) and be
// fully compatible with gopkg.in/yaml.v3. Where the YAML is persisted
// is not specified by this interface (see LocalConfigurer for interface
// with local requirement). The id may be anything that can be a YAML
// key, but use of Unicode Letter class runes is strongly recommended.
// Configuration data is designed to change less frequently than cache
// data (see cache package) and is always updated by editing the entire
// YAML file.
type Configurer interface {
	Init(id string) error          // must initialize a new configuration
	Data(id string) string         // must return full YAML
	Print(id string)               // must print full YAML to stdout
	Edit(id string) error          // must open full YAML in editor
	Write(id string, it any) error // must safely write (system semphore)
	Query(id, q string) string     // yq/jq compatible query string
	QueryPrint(id, q string)       // prints result to stdout
}

// LocalConfigurer implementations must use the local os.UserConfigDir
// to store all configuration information. All other requirements for
// Configurer apply. The id must be a valid directory name for all host
// operating systems on which the application may be installed.
// Lowercase Unicode Letters with no whitespace are strongly
// recommended.
type LocalConfigurer interface {
	Configurer
	Dir(id string) string
	File(id string) string
}
