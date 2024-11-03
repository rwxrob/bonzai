package sunrise

import (
	"strconv"
	"time"

	"github.com/rwxrob/bonzai/comp"
	bonzai "github.com/rwxrob/bonzai/pkg"
)

var Cmd = &bonzai.Cmd{
	Name:  `sunrise`,
	Vers:  `v0.1.0`,
	Short: `showcase all colors of terminal`,
	Comp:  comp.Cmds,

	MaxArgs: 1,

	Call: func(x *bonzai.Cmd, args ...string) error {
		var ms int64 = 10
		if len(args) > 0 {
			ms, _ = strconv.ParseInt(args[0], 10, 64)
		}
		Sunrise(time.Duration(ms) * time.Millisecond)
		return nil
	},
}
