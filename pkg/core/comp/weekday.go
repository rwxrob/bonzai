// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package comp

import (
	"github.com/rwxrob/bonzai/pkg/core/fn/filt"
)

type _ThreeLetterEngWeekday struct{}

var ThreeLetterEngWeekday = new(_ThreeLetterEngWeekday)

func (_ThreeLetterEngWeekday) Complete(_ any, args ...string) []string {
	list := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	if len(args) == 0 {
		return list
	}
	return filt.HasPrefix(list, args[0])
}
