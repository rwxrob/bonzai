package fishies

import (
	"strconv"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `fishies`,
	Alias: `fish`,
	Short: `terminal-based 3D fish animation`,
	Vers:  `v1.0.0`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		animateCmd,
		help.Cmd,
	},
	Def: animateCmd,
}

var animateCmd = &bonzai.Cmd{
	Name:  `animate`,
	Alias: `a|run|start`,
	Short: `run the fish animation with optional number of fish and ground`,
	Long: `
        Starts a terminal-based 3D animation of swimming fish.
        
        Arguments:
            number_of_fish: Number of fish to show (default: 3)
            ground: Show checkerboard ground (specify "ground" to enable)
        
        Controls:
            Ctrl+C: Exit the animation
        
        Examples:
            fishies         # Run with 3 fish, no ground
            fishies 5       # Run with 5 fish
            fishies 5 ground  # Run with 5 fish and ground`,
	Comp:    comp.Cmds,
	MaxArgs: 2,
	Do: func(x *bonzai.Cmd, args ...string) error {
		numFish := 3
		ground := false

		if len(args) > 0 {
			if n, err := strconv.Atoi(args[0]); err == nil && n > 0 {
				numFish = n
			}
		}
		if len(args) > 1 {
			ground = args[1] == "ground"
		}

		RunAnimation(numFish, ground)
		return nil
	},
}
