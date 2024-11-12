package help

import (
	"io"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
)

var Renderer mark.Renderer

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Short: `display command help`,
	Alias: `-h|--help|--h|/?`,
	Long: ` 
		The {{.Name}} command displays the help information for the
		immediate previous command unless it is passed arguments, in which
		case it resolves the arguments is if they were passed to the
		previous command itself and the leaf command help is displayed
		instead. Output is always Markdown.`,

	Do: func(x *bonzai.Cmd, args ...string) error {

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
