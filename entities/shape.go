package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Shape interface {
	MomentOfInertia() float64
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

func (circle *Circle) MomentOfInertia() float64 {
	return 0.5 * float64(circle.Radius) * float64(circle.Radius)
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
type Box struct {
	Color         uint32
	Type          ShapeType
	Width         float64
	Height        float64
	LocalVertices []vector.Vec2
	WorldVertices []vector.Vec2
}

// Passing the renderer is termporal for now
func NewBox(color uint32, width float64, height float64) *Box {
	return &Box{
		Color:  color,
		Type:   BOX,
		Width:  width,
		Height: height,
		LocalVertices: []vector.Vec2{
			{X: -width / 2.0, Y: -height / 2.0},
			{X: width / 2.0, Y: -height / 2.0},
			{X: width / 2.0, Y: height / 2.0},
			{X: -width / 2.0, Y: height / 2.0},
		},
		WorldVertices: []vector.Vec2{
			{X: -width / 2.0, Y: -height / 2.0},
			{X: width / 2.0, Y: -height / 2.0},
			{X: width / 2.0, Y: height / 2.0},
			{X: -width / 2.0, Y: height / 2.0},
		},
	}
}

func (box *Box) MomentOfInertia() float64 {
	return 0.083333 * ((box.Width * box.Width) + (box.Height * box.Height))
}

func (box *Box) Draw(x float64, y float64, rotation float64, renderer *renderer.Renderer) {

}
