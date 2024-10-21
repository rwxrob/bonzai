package bonzai

import "fmt"

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

type NotEnoughArgs struct {
	Count int
	Min   int
}

func (e NotEnoughArgs) Error() string {
	return fmt.Sprintf(E_NotEnoughArgs, e.Count, e.Min)
}

type TooManyArgs struct {
	Count int
	Max   int
}

func (e TooManyArgs) Error() string {
	return fmt.Sprintf(E_TooManyArgs, e.Count, e.Max)
}

type WrongNumArgs struct {
	Count int
	Num   int
}

func (e WrongNumArgs) Error() string {
	return fmt.Sprintf(E_WrongNumArgs, e.Count, e.Num)
}

type UsesConf struct {
	Cmd *Cmd
}

func (e UsesConf) Error() string {
	return fmt.Sprintf(E_UsesConf, e.Cmd.Name)
}

type UsesVars struct {
	Cmd *Cmd
}

func (e UsesVars) Error() string {
	return fmt.Sprintf(E_UsesVars, e.Cmd.Name)
}

type NoCallNoCommands struct {
	Cmd *Cmd
}

func (e NoCallNoCommands) Error() string {
	return fmt.Sprintf(E_NoCallNoCommands, e.Cmd.Name)
}

type DefCmdReqCall struct {
	Cmd *Cmd
}

func (e DefCmdReqCall) Error() string {
	return fmt.Sprintf(E_DefCmdReqCall, e.Cmd.Name)
}

type IncorrectUsage struct {
	Cmd *Cmd
}

func (e IncorrectUsage) Error() string {
	return fmt.Sprintf(E_IncorrectUsage, e.Cmd.Name, e.Cmd.FUsage())
}

type MultiCallCmdNotFound struct {
	CmdName string
}

func (e MultiCallCmdNotFound) Error() string {
	return fmt.Sprintf(E_MultiCallCmdNotFound, e.CmdName)
}

type MultiCallCmdNotCmd struct {
	CmdName string
	It      any
}

func (e MultiCallCmdNotCmd) Error() string {
	return fmt.Sprintf(E_MultiCallCmdNotCmd, e.CmdName, e.It)
}

type MultiCallCmdArgNotString struct {
	CmdName string
	It      any
}

func (e MultiCallCmdArgNotString) Error() string {
	return fmt.Sprintf(E_MultiCallCmdArgNotString, e.CmdName, e.It)
}

type MissingVar struct {
	Path string
}

func (e MissingVar) Error() string {
	return fmt.Sprintf(E_MissingVar, e.Path)
}
