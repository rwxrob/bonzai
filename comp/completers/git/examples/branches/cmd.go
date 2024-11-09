package main

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp/completers/git"
)

var Cmd = &bonzai.Cmd{Name: `branches`, Comp: git.CompBranches}
