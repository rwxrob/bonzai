package Z

import "fmt"

type NotEnoughArgs struct {
	Count int
	Min   int
}

func (e NotEnoughArgs) Error() string {
	return fmt.Sprintf(
		"error: %v is not enough args, %v required", e.Count, e.Min)
}

// -------------------------------- -- --------------------------------

type TooManyArgs struct {
	Count int
	Max   int
}

func (e TooManyArgs) Error() string {
	return fmt.Sprintf(
		"error: %v is too many args, %v maximum", e.Count, e.Max)
}

// -------------------------------- -- --------------------------------

type WrongNumArgs struct {
	Count int
	Num   int
}

func (e WrongNumArgs) Error() string {
	return fmt.Sprintf(
		"error: %v is wrong number of args, %v required", e.Count, e.Num)
}

// -------------------------------- -- --------------------------------

type MissingConf struct {
	Path string
}

func (e MissingConf) Error() string {
	return fmt.Sprintf("missing conf value for %v", e.Path)
}

// -------------------------------- -- --------------------------------

type MissingVar struct {
	Path string
}

func (e MissingVar) Error() string {
	return fmt.Sprintf("missing var for %v", e.Path)
}

// -------------------------------- -- --------------------------------

type UsesConf struct {
	Cmd *Cmd
}

func (e UsesConf) Error() string {
	return fmt.Sprintf("%v requires Z.Conf", e.Cmd.Name)
}

// -------------------------------- -- --------------------------------

type UsesVars struct {
	Cmd *Cmd
}

func (e UsesVars) Error() string {
	return fmt.Sprintf("%v requires Z.Vars", e.Cmd.Name)
}

// -------------------------------- -- --------------------------------

type NoCallNoCommands struct {
	Cmd *Cmd
}

func (e NoCallNoCommands) Error() string {
	return fmt.Sprintf("%v requires either Call or Commands", e.Cmd.Name)
}

// -------------------------------- -- --------------------------------

type DefCmdReqCall struct {
	Cmd *Cmd
}

func (e DefCmdReqCall) Error() string {
	return fmt.Sprintf("default (first) %q Commands requires Call", e.Cmd.Name)
}

// -------------------------------- -- --------------------------------

type IncorrectUsage struct {
	Cmd *Cmd
}

func (e IncorrectUsage) Error() string {
	return fmt.Sprintf("usage: %v %v", e.Cmd.Name, e.Cmd.GetUsage())
}

// -------------------------------- -- --------------------------------

type MultiCallCmdNotFound struct {
	CmdName string
}

func (e MultiCallCmdNotFound) Error() string {
	return fmt.Sprintf("multicall command not found: %v", e.CmdName)
}

// -------------------------------- -- --------------------------------

type MultiCallCmdNotCmd struct {
	CmdName string
	It      any
}

func (e MultiCallCmdNotCmd) Error() string {
	return fmt.Sprintf(
		"multicall match for %v, but first in slice not *Z.Cmd: %T", e.CmdName, e.It)
}

// -------------------------------- -- --------------------------------

type MultiCallCmdArgNotString struct {
	CmdName string
	It      any
}

func (e MultiCallCmdArgNotString) Error() string {
	return fmt.Sprintf(
		"multicall match for %v, but arg not string: %T", e.CmdName, e.It)
}
