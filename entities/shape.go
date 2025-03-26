package entities

import (
	"engine/renderer"
	"engine/vector"
	"math"
)

type Shape interface {
	MomentOfInertia(mass float64) float64
	Draw(x float64, y float64, rotation float64, renderer *renderer.Renderer)
}

type ShapeType int

const (
	CIRCLE ShapeType = iota
	POLYGON
	BOX
)

/*
Circle
*/
type Circle struct {
	Color  uint32
	Type   ShapeType
	Radius int32
}

func NewCircle(radius int32, color uint32) *Circle {
	return &Circle{
		Color:  color,
		Radius: radius,
		Type:   CIRCLE,
	}
}

func (circle *Circle) MomentOfInertia(mass float64) float64 {
	// TODO
	return float64(circle.Radius)
}

func (circle *Circle) Draw(x float64, y float64, rotation float64, renderer *renderer.Renderer) {
	renderer.DrawCircle(
		int32(x),
		int32(y),
		circle.Radius,
		rotation,
		circle.Color,
	)
}

/*
Polygon
*/
type Polygon struct {
	Color    uint32
	Type     ShapeType
	Vertices []vector.Vec2
}

func NewPolygon(color uint32, vertices []vector.Vec2) Polygon {
	return Polygon{
		Color:    color,
		Type:     POLYGON,
		Vertices: vertices,
	}
}

func (polygon *Polygon) MomentOfInertia(mass float64) float64 {
	// TODO
	return math.Pi
}

func (polygon *Polygon) Draw(x float64, y float64, renderer *renderer.Renderer) {
	// TODO
}
