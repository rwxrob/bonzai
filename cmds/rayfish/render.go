package rayfish

import (
	"gonum.org/v1/gonum/spatial/r3"
)

type Light struct {
	Direction r3.Vec
	Intensity float64
	Ambient   float64
}

type Intersectable interface {
	Intersect(origin, direction r3.Vec) (float64, Color, r3.Vec)
}

func RenderScene(scene *Scene, light *Light, width, height int, aspectRatio float64, pixelCallback func(x, y int, color Color)) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			screenX := (2.0*float64(x)/float64(width) - 1.0) * aspectRatio
			screenY := 1.0 - 2.0*float64(y)/float64(height)

			rayDir := r3.Unit(r3.Vec{screenX, screenY, 1.0})

			color := rayTrace(scene, scene.Camera, rayDir, light)

			pixelCallback(x, y, color)
		}
	}
}

func rayTrace(scene *Scene, origin, direction r3.Vec, light *Light) Color {
	minDist := -1.0
	var hitColor Color
	var hitNormal r3.Vec

	for _, obj := range scene.Objects {
		dist, color, normal := obj.Intersect(origin, direction)
		if dist > 0 && (minDist < 0 || dist < minDist) {
			minDist = dist
			hitColor = color
			hitNormal = normal
		}
	}

	if minDist < 0 {
		return Color{0, 0, 0}
	}

	diffuse := r3.Dot(hitNormal, light.Direction)
	if diffuse < 0 {
		diffuse = 0
	}

	lighting := light.Ambient + diffuse*light.Intensity
	return Color{
		R: hitColor.R * lighting,
		G: hitColor.G * lighting,
		B: hitColor.B * lighting,
	}
}
