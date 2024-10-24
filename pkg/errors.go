package bonzai

import (
	"fmt"
)

const (
	E_InvalidName      = `invalid name/alias detected: %v`
	E_InvalidMultiName = `%q must begin with %q: %q`
	E_DefCmdReqCall    = `default (first) %q Commands requires Call`
	E_IncorrectUsage   = `usage: %v %v`
	E_MissingVar       = `missing var for %v`
	E_NoCallNoCommands = `%v requires either Call or Commands`
	E_NotEnoughArgs    = `%v is not enough arguments, %v required`
	E_TooManyArgs      = `%v is too many arguments, %v maximum`
	E_WrongNumArgs     = `%v arguments, %v required`
	E_UnsupportedVar   = `unsupported var: %v`
	E_VarsInitFailed   = `var initialization failed: %v`
)

type VarsInitFailed struct {
	Err error
}

func (e VarsInitFailed) Error() string {
	return fmt.Sprintf(E_VarsInitFailed, e.Err)
}

type UnsupportedVar struct {
	Name string
}

func (e UnsupportedVar) Error() string {
	return fmt.Sprintf(E_UnsupportedVar, e.Name)
}

type InvalidMultiName struct {
	Got  string
	Want string
}

func (e InvalidMultiName) Error() string {
	return fmt.Sprintf(E_InvalidMultiName,
		e.Want, e.Got, e.Got+"-"+e.Want)
}

type InvalidName struct {
	Name string
}

func (e InvalidName) Error() string {
	return fmt.Sprintf(E_InvalidName, e.Name)
}

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
	return fmt.Sprintf(E_IncorrectUsage,
		e.Cmd.Name,
		e.Cmd.Fill(e.Cmd.Usage),
	)
}

type MissingVar struct {
	Path string
}

func (e MissingVar) Error() string {
	return fmt.Sprintf(E_MissingVar, e.Path)
}
