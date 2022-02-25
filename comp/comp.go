// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

// Completer defines a function to complete the given leaf Command with
// the provided arguments, if any. Completer functions must never be
// passed a nil Command or nil as the args slice. See comp.Standard.

type Completer func(leaf Command, args ...string) []string

// Command interface is only here to break cyclical package imports.
// This enables Completers of any kind to be create and managed
// independently.
type Command interface {
	GetName() string
	GetCommands() []string
	GetHidden() []string
	GetParams() []string
	GetOther() map[string]string
	GetCompleter() Completer
	GetCaller() Command
}
