/*
Copyright 2022 Robert S. Muhlestein.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package comp

// Completer defines a function to complete the given leaf Command with
// the provided arguments, if any. Completer functions must never be
// passed a nil Command or nil as the args slice. See comp.Standard.

type Completer func(leaf Command, args []string) []string

// Command interface is only here to break cyclical package imports.
// This enables Completers of any kind to be create and managed
// independently.
type Command interface {
	GetName() string
	GetCommands() []string
	GetHidden() []string
	GetParams() []string
	GetCompleter() Completer
}
