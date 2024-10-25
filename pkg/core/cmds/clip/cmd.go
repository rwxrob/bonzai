package clip

import (
	"os"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/vars"
)

var Cmd = &bonzai.Cmd{
	Name:    `clip`,
	Summary: `manage and play YouTube clips with mpv`,
	Version: `v0.0.1`,
	Cmds:    []*bonzai.Cmd{PlayCmd},
}

var PlayCmd = &bonzai.Cmd{
	Name: `play`,
	Call: func(_ *bonzai.Cmd, _ ...string) error {
		data := vars.Get(
			`.clip.data`,
			`CLIP_DATA`,
			os.UserCacheDir(),
		)
		_ = data
		//dir = os.Getenv(`CLIP_DIR`)
		//screen = os.Getenv(`CLIP_SCREEN`)

		/*
			: "${CLIP_DATA:="$HOME/.config/clip/data"}"
			  : "${CLIP_DIR:="$HOME/Movies/clips"}"
			  : "${CLIP_SCREEN:=2}"
			  : "${CLIP_VOLUME:=-50}"
			  : "${PAGER:=more}"
			  : "${EDITOR:=vi}"
			  : "${HELP_BROWSER:=}"
		*/
		return nil
	},
}
