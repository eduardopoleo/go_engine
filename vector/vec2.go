package vector

import (
	"math"
)

type Vec2 struct {
	X float64
	Y float64
}

func (vec *Vec2) Add(otherVec Vec2) Vec2 {
	return Vec2{
		X: vec.X + otherVec.X,
		Y: vec.Y + otherVec.Y,
	}
}

func (vec *Vec2) Subtract(otherVec Vec2) Vec2 {
	return Vec2{
		X: vec.X - otherVec.X,
		Y: vec.Y - otherVec.Y,
	}
}

func (vec *Vec2) Multiply(n float64) Vec2 {
	return Vec2{
		X: vec.X * n,
		Y: vec.Y * n,
	}
}

func (vec *Vec2) Magnitude() float64 {
	return float64(math.Sqrt((float64(vec.X) * float64(vec.X)) + (float64(vec.Y) * float64(vec.Y))))
}

func (vec *Vec2) Unit() Vec2 {
	magnitude := vec.Magnitude()
	return Vec2{
		X: vec.X / magnitude,
		Y: vec.Y / magnitude,
	}
}
