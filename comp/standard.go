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

import (
	"github.com/rwxrob/bonzai/filter"
)

// Standard completion is resolved as follows:
//
//     1. If leaf has Completer function, delegate to it
//
//     2. If leaf has no arguments, return all Commands and Params
//
//     3. If first argument is the name of a Command or one of
//        Aliases return only the Command name
//
//     4. Otherwise, return every Command or Param that is not in the
//        Hidden list and HasPrefix matching the first arg
//
// See comp.Completer.
func Standard(x Command, args []string) []string {

	// if has completer, delegate
	if c := x.GetCompleter(); c != nil {
		return c(x, args)
	}

	// fetch to avoid redundant calls
	commands := x.GetCommands()
	params := x.GetParams()
	hidden := x.GetHidden()

	// check for unique first argument command
	if len(args) == 0 {
		return []string{x.GetName()}
	}

	// build list of visible commands and params
	list := []string{}
	list = append(list, commands...)
	list = append(list, params...)
	list = filter.Minus(list, hidden)

	// catch edge case for first command
	if args[0] == " " {
		return list
	}

	return filter.HasPrefix(list, args[0])
}
