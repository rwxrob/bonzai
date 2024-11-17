package rayfish

import (
	"math"

	"gonum.org/v1/gonum/spatial/r3"
)

type Ground struct {
	Y float64
}

func NewGround(height float64) *Ground {
	return &Ground{Y: height}
}

func (g *Ground) Intersect(origin, direction r3.Vec) (float64, Color, r3.Vec) {
	if direction.Y == 0 {
		return -1, Color{0, 0, 0}, r3.Vec{}
	}

	t := (g.Y - origin.Y) / direction.Y
	if t < 0 {
		return -1, Color{0, 0, 0}, r3.Vec{}
	}

	hitPoint := r3.Add(origin, r3.Scale(t, direction))

	gridSize := 4.0
	x := math.Floor(hitPoint.X / gridSize)
	z := math.Floor(hitPoint.Z / gridSize)
	isEven := math.Mod(math.Abs(x+z), 2) < 1

	darkSquare := Color{0.0, 0.0, 0.0}
	lightSquare := Color{0.5, 0.5, 0.5}

	var color Color
	if isEven {
		color = darkSquare
	} else {
		color = lightSquare
	}

	return t, color, r3.Vec{0, 1, 0}
}
