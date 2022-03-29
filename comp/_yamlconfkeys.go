// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"github.com/rwxrob/fn/filt"
)

// YAMLKeys completer will look up the command in the local
// configuration with x.Branch and return a list of all keys in the
// local YAML configuration by calling its Config.
func YAMLKeys(x Command, args ...string) []string {
	// TODO need to get yq working first
	return filt.HasPrefix(list, args[0])
}
