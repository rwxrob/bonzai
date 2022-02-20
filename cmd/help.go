// Copyright 2022 Robert S. Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"log"

	"github.com/rwxrob/bonzai"
)

var Help = &bonzai.Cmd{
	Name:    `help`,
	Aliases: []string{"h"},
	Call: func(none ...string) error {
		log.Println("would print help")
		return nil
	},
}
