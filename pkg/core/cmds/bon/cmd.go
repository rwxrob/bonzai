package bon

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
)

var Cmd = &bonzai.Cmd{
	Name:    `bonzai`,
	Aliases: `bon|bonzaicli`,
	Version: `v0.0.1`,
	Vars: bonzai.Vars{
		`some`: thing,
	},
	Commands: []*bonzai.Cmd{
		doc.Cmd,
	},
}
