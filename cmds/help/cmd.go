package help

import (
	"io"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/mark/funcs"
	"github.com/rwxrob/bonzai/term"
	"github.com/rwxrob/bonzai/to"
)

var Renderer mark.Renderer

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Vers:  `v0.4.1`,
	Short: `display command help`,
	Alias: `-h|--help|--h|/?`,
	Long: `
		The {{aka .}} command displays the help information for the
		immediate previous command unless it is passed arguments, in which
		case it resolves the arguments as if they were passed to the
		previous command and the help for the leaf command is displayed
		instead.`,

	Do: func(x *bonzai.Cmd, args ...string) error {
		funcs := funcs.Map

		if len(args) > 0 {
			x, args = x.Caller().Seek(args)
		} else {
			x = x.Caller()
		}

		funcs = to.MergedMaps(funcs, x.Funcs)

		if Renderer != nil {
			r, err := Renderer.Render(x, funcs, x.Mark())
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

		x.Funcs = funcs
		term.Print(x.MarkString())

		return nil
	},
}
