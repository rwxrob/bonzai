package sunrise

import (
	"strconv"
	"time"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `sunrise`,
	Vers:  `v1.0.3`,
	Short: `showcase all colors of terminal`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{help.Cmd},
	Long: `
    Examples

        sunrise
        sunrise 5
        sunrise 20
	`,

	MaxArgs: 1,

	Do: func(x *bonzai.Cmd, args ...string) error {
		var ms int64 = 10
		if len(args) > 0 {
			ms, _ = strconv.ParseInt(args[0], 10, 64)
		}
		Sunrise(time.Duration(ms) * time.Millisecond)
		return nil
	},
}
