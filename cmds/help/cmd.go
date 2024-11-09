package help

import (
	"io"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
)

var Renderer mark.Renderer

var Cmd = &bonzai.Cmd{
	Name: `help`,

	Long: ` 
		The {{pre .Name}} command displays the help information for the
		immediate command before it unless *one*, or ***more***,
		***arguments*** is passed and matches a potential command path for

		# Already a Go thing

		~~~go
		fmt.Println("something")
		~~~

		1. adsfasdf
		2. jkaldsfkj

		- one
		- two

		* ten
		* twenty

		the previous command. In this way this command can be used at the
		top level so users can quickly add it to get essential help
		information about any command or just the previous command. `,

	Call: func(x *bonzai.Cmd, args ...string) error {

		if len(args) > 0 {
			x, args = x.Caller.Seek(args)
		} else {
			x = x.Caller
		}

		if Renderer != nil {
			r, err := Renderer.Render(x, x.Funcs, x.Mark())
			if err != nil {
				return err
			}
			out, err := io.ReadAll(r)
			if err != nil {
				return err
			}
			term.Print(string(out))
			return nil
		}

		term.Print(x.MarkString())

		return nil
	},
}
