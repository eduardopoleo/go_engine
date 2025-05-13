package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Polygon struct {
	Color         uint32
	LocalVertices []vector.Vec2
	WorldVertices []vector.Vec2
}

// Passing the renderer is termporal for now
func NewBox(color uint32, width float64, height float64) *Polygon {
	return &Polygon{
		Color: color,
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
	v10 := polygon.LocalVertices[1].Subtract(polygon.LocalVertices[0])
	v30 := polygon.LocalVertices[3].Subtract(polygon.LocalVertices[0])

	width := v10.Magnitude()
	height := v30.Magnitude()

	return 0.083333 * ((width * width) + (height * height))
}

func (polygon *Polygon) EdgeAt(idx int) vector.Vec2 {
	nextIdx := (idx + 1) % len(polygon.WorldVertices)
	return polygon.WorldVertices[nextIdx].Subtract(polygon.WorldVertices[idx])
}

func (box *Polygon) Draw(body *Body, rendr *renderer.Renderer) {
	for i := 0; i < len(box.WorldVertices); i++ {
		prev := box.WorldVertices[(i-1+len(box.WorldVertices))%len(box.WorldVertices)]
		curr := box.WorldVertices[i]
		rendr.DrawLine(prev, curr, box.Color)
	}
	// Just a dot in the position of the body.
	rendr.DrawCircle(int32(body.Position.X), int32(body.Position.Y), 1, 0, box.Color)
}

func (polygon *Polygon) MarkDebug() {
	polygon.Color = renderer.DEBUG
}

func (polygon *Polygon) UnMarkDebug() {
	polygon.Color = renderer.WHITE
}
