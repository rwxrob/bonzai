package locale

const (
	E_DefCmdReqCall            = "default (first) %q Commands requires Call"
	E_IncorrectUsage           = "usage: %v %v"
	E_MissingConf              = "missing conf value for %v"
	E_MissingVar               = "missing var for %v"
	E_MultiCallCmdArgNotString = "multicall match for %v, but arg not string: %T"
	E_MultiCallCmdNotCmd       = "multicall match for %v, but first in slice not *Z.Cmd: %T"
	E_MultiCallCmdNotFound     = "multicall command not found: %v"
	E_NoCallNoCommands         = "%v requires either Call or Commands"
	E_NotEnoughArgs            = `error: %v is not enough arguments, %v required`
	E_TooManyArgs              = `error: %v is too many arguments, %v maximum`
	E_UsesConf                 = "%v requires Z.Conf"
	E_UsesVars                 = "%v requires Z.Vars"
	E_WrongNumArgs             = `error: %v is wrong number of arguments, %v required`
)
