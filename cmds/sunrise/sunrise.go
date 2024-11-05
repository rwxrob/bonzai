package sunrise

import (
	"fmt"
	"math"
	"time"

	"github.com/rwxrob/bonzai/anim"
)

var (
	p         = 3.14
	i float64 = 0
)

// Sunrise prints a wall of all colors in the terminal. The speed of the
// animation can be adjusted with the speed parameter.
func Sunrise(speed time.Duration) {
	anim.SimpleAnimationScreen()
	for {
		i += 0.04
		r := int(128 + 127*math.Sin(i))
		g := int(128 + 127*math.Sin(i+p*(1.0/3)))
		b := int(128 + 127*math.Sin(i+p*(2.0/3)))
		fmt.Printf("\033[48;2;%d;%d;%dm\n", r, g, b)
		time.Sleep(speed) // Sleep for 10ms
	}
}
