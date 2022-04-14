package Z

import "fmt"

type NotEnoughArgs struct {
	Cmd *Cmd
}

func (e NotEnoughArgs) Error() string {
	return fmt.Sprintf("not enough args, %v required", e.Cmd.MinArgs)
}

type TooManyArgs struct {
	Cmd *Cmd
}

func (e TooManyArgs) Error() string {
	return fmt.Sprintf("too many args, %v maximum", e.Cmd.MaxArgs)
}

type WrongNumArgs struct {
	Cmd *Cmd
}

func (e WrongNumArgs) Error() string {
	return fmt.Sprintf("wrong number of args, %v required", e.Cmd.NumArgs)
}

type ConfRequired struct {
	Cmd *Cmd
}

func (e ConfRequired) Error() string {
	return fmt.Sprintf("%v requires Z.Conf", e.Cmd.Name)
}

type VarsRequired struct {
	Cmd *Cmd
}

func (e VarsRequired) Error() string {
	return fmt.Sprintf("%v requires Z.Vars", e.Cmd.Name)
}

type NoCallNoCommands struct {
	Cmd *Cmd
}

func (e NoCallNoCommands) Error() string {
	return fmt.Sprintf("%v requires either Call or Commands", e.Cmd.Name)
}

type DefCmdReqCall struct {
	Cmd *Cmd
}

func (e DefCmdReqCall) Error() string {
	return fmt.Sprintf("default (first) %q Commands requires Call", e.Cmd.Name)
}

type IncorrectUsage struct {
	Cmd *Cmd
}

func (e IncorrectUsage) Error() string {
	return fmt.Sprintf("usage: %v %v", e.Cmd.Name, UsageFunc(e.Cmd))
}

type MultiCallCmdNotFound struct {
	CmdName string
}

func (e MultiCallCmdNotFound) Error() string {
	return fmt.Sprintf("multicall command not found: %v", e.CmdName)
}

type MultiCallCmdNotCmd struct {
	CmdName string
	It      any
}

func (e MultiCallCmdNotCmd) Error() string {
	return fmt.Sprintf(
		"multicall match for %v, but first in slice not *Z.Cmd: %T", e.CmdName, e.It)
}

type MultiCallCmdArgNotString struct {
	CmdName string
	It      any
}

func (e MultiCallCmdArgNotString) Error() string {
	return fmt.Sprintf(
		"multicall match for %v, but arg not string: %T", e.CmdName, e.It)
}
