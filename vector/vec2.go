package vector

import (
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func (a *Vec2) Add(b Vec2) Vec2 {
	return Vec2{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

func (a *Vec2) Subtract(b Vec2) Vec2 {
	return Vec2{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func (a *Vec2) Dot(b Vec2) float64 {
	return a.X*b.X + a.Y*b.Y
}

// Alias to remember it easier
func (a *Vec2) ProjOnTo(b Vec2) float64 {
	return a.Dot(b)
}

/*
a = [Ax, Ay]
b = [Bx, By]

Ax * By - Ay * Bx
*/
func (a *Vec2) Cross(b Vec2) float64 {
	return a.X*b.Y - a.Y*b.X
}

// Scale could be a better name for this
func (a *Vec2) Multiply(n float64) Vec2 {
	return Vec2{
		X: a.X * n,
		Y: a.Y * n,
	}
}

func (a *Vec2) Magnitude() float64 {
	return float64(math.Sqrt((float64(a.X) * float64(a.X)) + (float64(a.Y) * float64(a.Y))))
}

func (a *Vec2) Rotate(angle float64) Vec2 {
	return Vec2{
		X: a.X*math.Cos(angle) - a.Y*math.Sin(angle),
		Y: a.X*math.Sin(angle) + a.Y*math.Cos(angle),
	}
}

// verify how is this implemented in the reference implementation
func (a *Vec2) Unit() Vec2 {
	magnitude := a.Magnitude()
	return Vec2{
		X: a.X / magnitude,
		Y: a.Y / magnitude,
	}
}

func (vec *Vec2) Normal() Vec2 {
	magnitude := vec.Magnitude()
	return Vec2{
		X: vec.Y / magnitude,
		Y: -vec.X / magnitude,
	}
}
