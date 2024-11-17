package env

import "github.com/rwxrob/bonzai"

// vars implements completion for environment variables with $ prefix.
type vars struct {
	configuredCompleter
}

// CompVars is a [bonzai.Completer] that completes environment variables with case sensitivity.
var CompVars bonzai.Completer = NewCompVars("", false)

// NewCompVars creates a new [bonzai.Completer] that completes environment variables with $ prefix.
//
// Parameters:
//   - prefix: String to prepend to environment variable search
//   - insensitive: If true, matching will be case-insensitive
//
// Returns:
//   - A [bonzai.Completer] configured with the specified options
func NewCompVars(prefix string, insensitive bool) bonzai.Completer {
	return &vars{
		configuredCompleter: configuredCompleter{
			opts: CompletionOptions{
				Prefix:      prefix,
				Insensitive: insensitive,
			},
		},
	}
}

func (v *vars) Complete(args ...string) []string {
	matches := v.complete(args)
	for i := range matches {
		matches[i] = "$" + matches[i]
	}
	return matches
}
