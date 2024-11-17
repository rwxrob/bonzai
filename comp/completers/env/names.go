package env

import "github.com/rwxrob/bonzai"

// names implements completion for environment variable names.
type names struct {
	configuredCompleter
}

// CompNames is a [bonzai.Completer] that completes environment variable names with case sensitivity.
var CompNames bonzai.Completer = NewCompNames("", false)

// NewCompNames creates a new [bonzai.Completer] that completes environment variable names.
//
// Parameters:
//   - prefix: String to prepend to environment variable search
//   - insensitive: If true, matching will be case-insensitive
//
// Returns:
//   - A [bonzai.Completer] configured with the specified options
func NewCompNames(prefix string, insensitive bool) bonzai.Completer {
	return &names{
		configuredCompleter: configuredCompleter{
			opts: CompletionOptions{
				Prefix:      prefix,
				Insensitive: insensitive,
			},
		},
	}
}

func (n *names) Complete(args ...string) []string {
	return n.complete(args)
}
