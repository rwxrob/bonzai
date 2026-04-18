package help

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
)

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Alias: `h|-h|--help|--h|/?`,
	Vers:  `v0.9.0`,
	Short: `display command help`,
	Long: `
		The {{code .Name}} command displays the help information for the
		immediate previous command unless it is passed arguments, in which
		case it resolves the arguments as if they were passed to the
		previous command and the help for the leaf command is displayed
		instead.`,

	Do: func(x *bonzai.Cmd, args ...string) (err error) {
		if len(args) > 0 {
			x, args, err = x.Caller().SeekInit(args...)
			if err != nil {
				return err
			}
		} else {
			x = x.Caller()
		}
		return mark.Print(x)
	},
}

