// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

// give own completer for days of the week
type _ThreeLetterEngWeekday struct{}

var ThreeLetterEngWeekday = new(_ThreeLetterEngWeekday)

func (ThreeLetterEngWeekday) Complete(_ *bonzai.Cmd, args ...string) []string {
	list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	if len(args) == 0 {
		return list
	}
	return filt.HasPrefix(list, args[0])
}
