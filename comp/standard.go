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

import "github.com/rwxrob/bonzai/filter"

// Standard completion is resolved as follows:
//
//     1. If leaf has Completer function, delegate to it
//
//     2. If leaf has no arguments, return all Commands and Params
//
//     3. Otherwise, return every Command or Param that is not in the
//        Hidden list and HasPrefix matching the first arg
//
// See comp.Completer.
func Standard(x Command, args []string) []string {

	if c := x.GetCompleter(); c != nil {
		return c(x, args)
	}

	list := []string{}
	list = append(list, x.GetCommands()...)
	list = append(list, x.GetParams()...)
	list = filter.Minus(list, x.GetHidden())

	if len(args) > 0 {
		list = filter.HasPrefix(list, args[0])
	}

	return list
}
