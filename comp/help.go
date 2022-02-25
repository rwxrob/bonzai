// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"github.com/rwxrob/bonzai/check"
	"github.com/rwxrob/bonzai/filter"
)

func Help(x Command, args ...string) []string {

	// check for unique first argument command
	if check.Blank(args) {
		return []string{x.GetName()}
	}

	// build list of params and other keys
	list := []string{}
	list = append(list, x.GetParams()...)
	// FIXME GetCaller nil check
	list = append(list, filter.Keys(x.GetCaller().GetOther())...)

	return filter.HasPrefix(list, args[0])
}
