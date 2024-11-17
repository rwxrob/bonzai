package rayfish

import (
	"math"
	"math/rand"

	"gonum.org/v1/gonum/spatial/r3"
)

type Fish struct {
	Center                 r3.Vec
	InitialPos             r3.Vec
	Velocity               r3.Vec
	Radius                 float64
	Rotation               float64
	TailRotation           float64
	Color                  Color
	elapsedTime            float64
	freqX1, freqX2, freqX3 float64
	freqZ1, freqZ2, freqZ3 float64
	ampX1, ampX2, ampX3    float64
	ampZ1, ampZ2, ampZ3    float64
	phaseX2, phaseX3       float64
	phaseZ2, phaseZ3       float64
}

func NewFish(pos r3.Vec, color Color, radius float64) *Fish {
	return &Fish{
		Center:     pos,
		InitialPos: pos,
		Color:      color,
		Radius:     radius,
		freqX1:     0.1 + rand.Float64()*0.7,
		freqX2:     0.1 + rand.Float64()*0.7,
		freqX3:     0.1 + rand.Float64()*0.7,
		freqZ1:     0.1 + rand.Float64()*0.7,
		freqZ2:     0.1 + rand.Float64()*0.7,
		freqZ3:     0.1 + rand.Float64()*0.7,
		ampX1:      0.5 + rand.Float64()*0.3,
		ampX2:      0.2 + rand.Float64()*0.2,
		ampX3:      0.1 + rand.Float64()*0.1,
		ampZ1:      0.5 + rand.Float64()*0.3,
		ampZ2:      0.2 + rand.Float64()*0.2,
		ampZ3:      0.1 + rand.Float64()*0.1,
		phaseX2:    rand.Float64() * 2 * 3.14159,
		phaseX3:    rand.Float64() * 2 * 3.14159,
		phaseZ2:    rand.Float64() * 2 * 3.14159,
		phaseZ3:    rand.Float64() * 2 * 3.14159,
	}
}

func (f *Fish) Update(deltaTime float64) {
	const (
		BOUNDS_SCALE   float64 = 16.0
		TAIL_FREQ      float64 = 4.0
		TAIL_AMPLITUDE float64 = 0.6
		MOVE_SPEED     float64 = 1.0
	)

	f.elapsedTime += deltaTime

	t := f.elapsedTime * MOVE_SPEED

	x := math.Sin(t*f.freqX1)*f.ampX1 +
		math.Sin(t*f.freqX2+f.phaseX2)*f.ampX2 +
		math.Sin(t*f.freqX3+f.phaseX3)*f.ampX3

	z := math.Cos(t*f.freqZ1)*f.ampZ1 +
		math.Cos(t*f.freqZ2+f.phaseZ2)*f.ampZ2 +
		math.Cos(t*f.freqZ3+f.phaseZ3)*f.ampZ3

	f.Center.X = x * BOUNDS_SCALE
	f.Center.Z = z * BOUNDS_SCALE
	f.Center.Y = f.InitialPos.Y + math.Sin(t*0.2)*0.3

	nextT := t + 0.01
	nextX := math.Sin(nextT*f.freqX1)*f.ampX1 +
		math.Sin(nextT*f.freqX2+f.phaseX2)*f.ampX2 +
		math.Sin(nextT*f.freqX3+f.phaseX3)*f.ampX3

	nextZ := math.Cos(nextT*f.freqZ1)*f.ampZ1 +
		math.Cos(nextT*f.freqZ2+f.phaseZ2)*f.ampZ2 +
		math.Cos(nextT*f.freqZ3+f.phaseZ3)*f.ampZ3

	nextX *= BOUNDS_SCALE
	nextZ *= BOUNDS_SCALE

	f.Rotation = -math.Atan2(nextZ-f.Center.Z, nextX-f.Center.X)

	tailSpeed := math.Sqrt(math.Pow(nextX-f.Center.X, 2) + math.Pow(nextZ-f.Center.Z, 2))
	f.TailRotation = TAIL_AMPLITUDE * math.Sin(f.elapsedTime*TAIL_FREQ) * (0.5 + tailSpeed)
}

func (f *Fish) rotatePoint(p r3.Vec) r3.Vec {
	sin, cos := math.Sin(f.Rotation), math.Cos(f.Rotation)
	return r3.Vec{
		p.X*cos - p.Z*sin,
		p.Y,
		p.X*sin + p.Z*cos,
	}
}

func (f *Fish) unrotatePoint(p r3.Vec) r3.Vec {
	sin, cos := math.Sin(f.Rotation), math.Cos(f.Rotation)
	return r3.Vec{
		p.X*cos + p.Z*sin,
		p.Y,
		-p.X*sin + p.Z*cos,
	}
}

func (f *Fish) Intersect(origin, direction r3.Vec) (float64, Color, r3.Vec) {
	localOrigin := r3.Sub(origin, f.Center)
	localOrigin = f.rotatePoint(localOrigin)
	localDir := f.rotatePoint(direction)

	var minT float64 = -1
	var hitNormal r3.Vec
	var hitColor Color

	scaleX := 1.8
	scaleY := 0.8
	scaleZ := 0.9

	scaledOrigin := r3.Vec{
		localOrigin.X / (scaleX * f.Radius),
		localOrigin.Y / (scaleY * f.Radius),
		localOrigin.Z / (scaleZ * f.Radius),
	}
	scaledDir := r3.Vec{
		localDir.X / (scaleX * f.Radius),
		localDir.Y / (scaleY * f.Radius),
		localDir.Z / (scaleZ * f.Radius),
	}

	a := r3.Dot(scaledDir, scaledDir)
	b := 2.0 * r3.Dot(scaledOrigin, scaledDir)
	c := r3.Dot(scaledOrigin, scaledOrigin) - 1.0

	discriminant := b*b - 4*a*c
	if discriminant >= 0 {
		t := (-b - math.Sqrt(discriminant)) / (2 * a)
		if t < 0 {
			t = (-b + math.Sqrt(discriminant)) / (2 * a)
		}
		if t >= 0 {
			minT = t
			hitPoint := r3.Add(localOrigin, r3.Scale(t, localDir))
			hitNormal = r3.Vec{
				hitPoint.X / (scaleX * scaleX * f.Radius * f.Radius),
				hitPoint.Y / (scaleY * scaleY * f.Radius * f.Radius),
				hitPoint.Z / (scaleZ * scaleZ * f.Radius * f.Radius),
			}
			hitNormal = r3.Unit(hitNormal)
			hitColor = f.Color

			eyeX := 1.4 * f.Radius
			eyeY := 0.2 * f.Radius
			eyeZ := 0.4 * f.Radius
			eyeRadius := 0.3 * f.Radius

			eyePos := r3.Vec{eyeX, eyeY, eyeZ}
			relativePoint := r3.Sub(hitPoint, eyePos)
			if r3.Norm(relativePoint) <= eyeRadius {
				hitColor = Color{1, 1, 1}
				hitNormal = r3.Unit(relativePoint)
			}

			eyePos = r3.Vec{eyeX, eyeY, -eyeZ}
			relativePoint = r3.Sub(hitPoint, eyePos)
			if r3.Norm(relativePoint) <= eyeRadius {
				hitColor = Color{1, 1, 1}
				hitNormal = r3.Unit(relativePoint)
			}
		}
	}

	tailBase := r3.Vec{-1.5 * f.Radius, 0, 0}
	tailTop := r3.Vec{-3.0 * f.Radius, 0.9 * f.Radius, 0}
	tailBottom := r3.Vec{-3.0 * f.Radius, -0.9 * f.Radius, 0}

	sin, cos := math.Sin(f.TailRotation), math.Cos(f.TailRotation)

	tailTop = r3.Vec{
		tailBase.X + (tailTop.X-tailBase.X)*cos - (tailTop.Z-tailBase.Z)*sin,
		tailTop.Y,
		tailBase.Z + (tailTop.X-tailBase.X)*sin + (tailTop.Z-tailBase.Z)*cos,
	}

	tailBottom = r3.Vec{
		tailBase.X + (tailBottom.X-tailBase.X)*cos - (tailBottom.Z-tailBase.Z)*sin,
		tailBottom.Y,
		tailBase.Z + (tailBottom.X-tailBase.X)*sin + (tailBottom.Z-tailBase.Z)*cos,
	}

	tailPoints := []r3.Vec{
		tailBase,
		tailTop,
		tailBottom,
	}

	if t, hit := rayTriangleIntersect(localOrigin, localDir, tailPoints[0], tailPoints[1], tailPoints[2]); hit && (t < minT || minT < 0) {
		minT = t
		hitColor = f.Color
		planeNormal := r3.Unit(r3.Cross(
			r3.Sub(tailPoints[1], tailPoints[0]),
			r3.Sub(tailPoints[2], tailPoints[0]),
		))
		if r3.Dot(planeNormal, localDir) > 0 {
			hitNormal = r3.Scale(-1, planeNormal)
		} else {
			hitNormal = planeNormal
		}
	}

	if minT < 0 {
		return -1, Color{0, 0, 0}, r3.Vec{}
	}

	hitNormal = f.unrotatePoint(hitNormal)
	return minT, hitColor, hitNormal
}

func rayTriangleIntersect(origin, dir, v0, v1, v2 r3.Vec) (float64, bool) {
	epsilon := 0.0000001

	edge1 := r3.Sub(v1, v0)
	edge2 := r3.Sub(v2, v0)
	h := r3.Cross(dir, edge2)
	a := r3.Dot(edge1, h)

	if a > -epsilon && a < epsilon {
		return 0, false
	}

	f := 1.0 / a
	s := r3.Sub(origin, v0)
	u := f * r3.Dot(s, h)

	if u < 0.0 || u > 1.0 {
		return 0, false
	}

	q := r3.Cross(s, edge1)
	v := f * r3.Dot(dir, q)

	if v < 0.0 || u+v > 1.0 {
		return 0, false
	}

	t := f * r3.Dot(edge2, q)
	if t > epsilon {
		return t, true
	}

	return 0, false
}
