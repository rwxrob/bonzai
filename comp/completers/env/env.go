// Package env provides environment variable completion functionality for the Bonzai command tree.
// It supports case-sensitive and case-insensitive completion of environment variable names
// and values, with optional prefixing.
package env

import (
	"os"
	"sort"
	"strings"
)

// CompletionOptions defines the configuration for environment variable completion.
type CompletionOptions struct {
	// Prefix is prepended to the environment variable search.
	// For example, if Prefix is "APP_", only environment variables
	// starting with "APP_" will be included in completions.
	Prefix string

	// Insensitive determines if case-insensitive matching should be used.
	// When true, "path" will match "PATH", "Path", "path", etc.
	Insensitive bool
}

type configuredCompleter struct {
	opts CompletionOptions
}

func (c configuredCompleter) complete(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}

	prefix := c.opts.Prefix
	arg := args[0]
	if c.opts.Insensitive {
		prefix = strings.ToLower(prefix)
		arg = strings.ToLower(arg)
	}

	// If the argument is the start of the prefix, treat as empty
	if strings.HasPrefix(prefix, arg) {
		arg = ""
	}

	// Remove leading $ for common completion
	arg = strings.TrimPrefix(arg, "$")

	// Remove prefix from the argument to preserve continued completion
	arg = strings.TrimPrefix(arg, prefix)

	return c.findMatches(prefix, arg)
}

func (c configuredCompleter) findMatches(prefix, arg string) []string {
	var matches []string

	for _, env := range os.Environ() {
		key, _, found := strings.Cut(env, "=")
		if found && key != "" && c.isMatch(key, prefix, arg) {
			matches = append(matches, key)
		}
	}

	sort.Strings(matches)
	return matches
}

func (c configuredCompleter) isMatch(key, prefix, arg string) bool {
	if c.opts.Insensitive {
		key = strings.ToLower(key)
	}
	return strings.HasPrefix(key, prefix+arg)
}
