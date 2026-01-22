package main

import (
	"github.com/BuddhiLW/bonzai"
	"github.com/BuddhiLW/bonzai/comp/completers/git"
)

var Cmd = &bonzai.Cmd{
	Name: `test`,
	Comp: git.CompTags,
}
