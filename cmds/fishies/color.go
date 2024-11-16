package fishies

import (
	"math/rand"
)

type Color struct {
	R, G, B float64
}

func randomColor() Color {
	return Color{
		R: 0.3 + rand.Float64()*0.7,
		G: 0.3 + rand.Float64()*0.7,
		B: 0.3 + rand.Float64()*0.7,
	}
}
