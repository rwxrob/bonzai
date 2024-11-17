package rayfish

import (
	"math/rand"

	"gonum.org/v1/gonum/spatial/r3"
)

type Scene struct {
	Camera  r3.Vec
	Objects []Intersectable
}

func CreateScene(numFish int, hasGround bool) (*Scene, *Light) {
	size := numFish
	if hasGround {
		size++
	}
	objects := make([]Intersectable, size)

	baseHeight := 2.0
	heightStep := 2.5

	for i := 0; i < numFish; i++ {
		pos := r3.Vec{
			X: rand.Float64()*4.0 - 2.0,
			Y: baseHeight + float64(i)*heightStep + rand.Float64()*1.0 - 0.5,
			Z: rand.Float64()*4.0 - 2.0,
		}

		size := 3.0 + rand.Float64()*2.0
		objects[i] = NewFish(pos, randomColor(), size)
	}

	if hasGround {
		objects[numFish] = NewGround(0.0)
	}

	light := &Light{
		Direction: r3.Unit(r3.Vec{-0.5, 2, -0.5}),
		Intensity: 0.9,
		Ambient:   0.3,
	}

	return &Scene{
		Camera:  r3.Vec{0, 7.0, -16.0},
		Objects: objects,
	}, light
}

func UpdateScene(scene *Scene, deltaTime float64) {
	for _, obj := range scene.Objects {
		if fish, ok := obj.(*Fish); ok {
			fish.Update(deltaTime)
		}
	}
}
