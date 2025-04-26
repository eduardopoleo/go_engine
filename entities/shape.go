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
	BOX
	POLYGON
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
type Polygon struct {
	Color         uint32
	Type          ShapeType
	LocalVertices []vector.Vec2
	WorldVertices []vector.Vec2
}

// Passing the renderer is termporal for now
func NewBox(color uint32, width float64, height float64) *Polygon {
	return &Polygon{
		Color: color,
		Type:  BOX,
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

/*
TODO: How do we extend this interface for shapes with different amount of vertices.
How do calculate the momentOfInertia for a shape do not know?
*/
func (polygon *Polygon) MomentOfInertia() float64 {
	switch polygonType := polygon.Type; polygonType {
	case BOX:
		v10 := polygon.LocalVertices[1].Subtract(polygon.LocalVertices[0])
		v30 := polygon.LocalVertices[3].Subtract(polygon.LocalVertices[0])

		width := v10.Magnitude()
		height := v30.Magnitude()

		return 0.083333 * ((width * width) + (height * height))
	}
	/*
		How do you return error?
	*/
	return 0
}

func (polygon *Polygon) EdgeAt(idx int) vector.Vec2 {
	nextIdx := (idx + 1) % len(polygon.WorldVertices)
	return polygon.WorldVertices[nextIdx].Subtract(polygon.WorldVertices[idx])
}

func (box *Polygon) Draw(x float64, y float64, rotation float64, rendr *renderer.Renderer) {
	rendr.DrawLine(box.WorldVertices[3], box.WorldVertices[0], renderer.WHITE)
	for i := 1; i < len(box.WorldVertices); i++ {
		prev := box.WorldVertices[i-1]
		curr := box.WorldVertices[i]
		rendr.DrawLine(prev, curr, renderer.WHITE)
	}
	// Just a dot in the position of the body.
	rendr.DrawCircle(int32(x), int32(y), 1, 0, renderer.WHITE)
}
