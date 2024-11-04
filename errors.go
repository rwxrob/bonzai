package bonzai

import (
	"fmt"
)

type VarsInitFailed struct {
	Err error
}

func (e VarsInitFailed) Error() string {
	return fmt.Sprintf(`var initialization failed: %v`, e.Err)
}

type UnsupportedVar struct {
	Name string
}

func (e UnsupportedVar) Error() string {
	return fmt.Sprintf(
		`unsupported var: %v`, e.Name)
}

type InvalidMultiName struct {
	Got  string
	Want string
}

func (e InvalidMultiName) Error() string {
	return fmt.Sprintf(`%q must begin with %q: %q`,
		e.Want, e.Got, e.Got+"-"+e.Want)
}

type InvalidName struct {
	Name string
}

func (e InvalidName) Error() string {
	return fmt.Sprintf(`invalid name: %v`, e.Name)
}

type NotEnoughArgs struct {
	Count int
	Min   int
}

func (e NotEnoughArgs) Error() string {
	return fmt.Sprintf(`%v is not enough arguments, %v required`,
		e.Count, e.Min)
}

type TooManyArgs struct {
	Count int
	Max   int
}

func (e TooManyArgs) Error() string {
	return fmt.Sprintf(`%v is too many arguments, %v maximum`,
		e.Count, e.Max)
}

type WrongNumArgs struct {
	Count int
	Num   int
}

func (e WrongNumArgs) Error() string {
	return fmt.Sprintf(
		`%v arguments, %v required`,
		e.Count, e.Num)
}

type Uncallable struct {
	Cmd *Cmd
}

func (e Uncallable) Error() string {
	return fmt.Sprintf(`Cmd requires Call or Def: %v`, e.Cmd.Name)
}

type CallOrDef struct {
	Cmd *Cmd
}

func (e CallOrDef) Error() string {
	return fmt.Sprintf(`Call or Def (not both): %v`, e.Cmd.Name)
}

type NoCallNoDef struct {
	Cmd *Cmd
}

func (e NoCallNoDef) Error() string {
	return fmt.Sprintf(`either Call or Def required if no Cmds: %v`, e.Cmd.Name)
}

type IncorrectUsage struct {
	Cmd *Cmd
}

func (e IncorrectUsage) Error() string {
	return fmt.Sprintf(`usage: %v %v`,
		e.Cmd.Name,
		e.Cmd.Fill(e.Cmd.Usage),
	)
}

type MissingVar struct {
	Path string
}

func (e MissingVar) Error() string {
	return fmt.Sprintf(`missing var for %v`, e.Path)
}
