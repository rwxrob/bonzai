package comp

import "github.com/rwxrob/bonzai"

var (
	current                *bonzai.Cmd // [bonzai.CmdCompleter]
	CmdsOpts               = Combine{Cmds, Opts}
	CmdsAliases            = Combine{Cmds, Aliases}
	CmdsAliasesOpts        = Combine{Cmds, Aliases, Opts}
	FileDirCmdsOpts        = Combine{FileDir, CmdsOpts}
	FileDirCmdsAliasesOpts = Combine{FileDir, CmdsAliasesOpts}
)
