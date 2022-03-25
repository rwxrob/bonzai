package main

import (
	"fmt"
	"strings"

	Z "github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/inc/help"
	"github.com/rwxrob/term"
)

func main() {
	Cmd.Run()
}

var Cmd = &Z.Cmd{
	Name:      `bonzai`,
	Summary:   `bonzai command tree utility`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, sh2slice},
}

var sh2slice = &Z.Cmd{
	Name:     `sh2slice`,
	Summary:  `splits a shell command into arguments`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(_ *Z.Cmd, args ...string) error {
		list := []string{}
		if len(args) == 0 {
			args = append(args, term.Read())
		}
		for _, a := range args {
			// FIXME add awareness or globs and quoted segments
			for _, aa := range strings.Fields(a) {
				list = append(list, fmt.Sprintf("%q", aa))
			}
		}
		fmt.Println(strings.Join(list, ","))
		return nil
	},
}
