package Z

import (
	"fmt"
	"strings"

	"github.com/rwxrob/fn/filt"
)

// UsageGroup uses Bonzai usage notation, a basic form of regular
// expressions, to describe the arguments allowed where each argument is
// a literal string (avoid spaces). The arguments are joined with bars
// (|) and wrapped with parentheses producing a regex group.  The min
// and max are then applied by adding the following regex decorations
// after the final parenthesis:
//
//                - min=1 max=1 (exactly one)
//     ?          - min=0 max=0 (none or many)
//     +          - min=1 max=0 (one or more)
//     {min,}     - min>0 max=0 (min, no max)
//     {min,max}  - min>0 max>0 (min and max)
//     {,max}     - min=0 max>0 (max, no min)
//
// An empty args slice returns an empty string. If only one arg, then
// that arg is simply returned and min and max are ignored. Arguments
// that are empty strings are ignored. No transformation is done to the
// string itself (such as removing white space).
func UsageGroup(args []string, min, max int) string {
	args = filt.NotEmpty(args)
	switch len(args) {
	case 0:
		return ""
	case 1:
		return args[0]
	default:
		var dec string
		switch {
		case min == 1 && max == 1:
		case min == 0 && max == 0:
			dec = "?"
		case min == 1 && max == 0:
			dec = "+"
		case min > 1 && max == 0:
			dec = fmt.Sprintf("{%v,}", min)
		case min > 0 && max > 0:
			dec = fmt.Sprintf("{%v,%v}", min, max)
		case min == 0 && max > 1:
			dec = fmt.Sprintf("{,%v}", max)
		}
		return "(" + strings.Join(args, "|") + ")" + dec
	}
}
