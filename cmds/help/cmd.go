package help

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
)

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Alias: `h|-h|--help|--h|/?`,
	Vers:  `v0.7.7`,
	Short: `display command help`,
	Long: `
		The {{aka .}} command displays the help information for the
		immediate previous command unless it is passed arguments, in which
		case it resolves the arguments as if they were passed to the
		previous command and the help for the leaf command is displayed
		instead.`,

	Do: func(x *bonzai.Cmd, args ...string) error {

		if len(args) > 0 {
			x, args = x.Caller().Seek(args...)
		} else {
			x = x.Caller()
		}

		out, err := mark.Usage(x)
		if err != nil {
			return err
		}

		term.Print(out)

		return nil
	},
}
