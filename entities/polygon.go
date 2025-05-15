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

func (polygon *Polygon) GetHeight() float64 {
	vec := polygon.LocalVertices[1].Subtract(polygon.LocalVertices[0])
	return vec.Magnitude()
}

func (polygon *Polygon) GetWidth() float64 {
	vec := polygon.LocalVertices[3].Subtract(polygon.LocalVertices[0])
	return vec.Magnitude()
}

/*
TODO: How do we extend this interface for shapes with different amount of vertices.
How do calculate the momentOfInertia for a shape do not know?
*/
func (polygon *Polygon) MomentOfInertia() float64 {
	if len(polygon.WorldVertices) == 4 {
		v10 := polygon.LocalVertices[1].Subtract(polygon.LocalVertices[0])
		v30 := polygon.LocalVertices[3].Subtract(polygon.LocalVertices[0])

		width := v10.Magnitude()
		height := v30.Magnitude()

		return 0.083333 * ((width * width) + (height * height))
	} else {
		return 5000
	}
}

func (polygon *Polygon) EdgeAt(idx int) vector.Vec2 {
	nextIdx := (idx + 1) % len(polygon.WorldVertices)
	return polygon.WorldVertices[nextIdx].Subtract(polygon.WorldVertices[idx])
}

func (polygon *Polygon) Draw(body *Body, rendr *renderer.Renderer) {
	for i := 0; i < len(polygon.WorldVertices); i++ {
		prev := polygon.WorldVertices[(i-1+len(polygon.WorldVertices))%len(polygon.WorldVertices)]
		curr := polygon.WorldVertices[i]
		rendr.DrawLine(prev, curr, polygon.Color)
	}

	// Just a dot in the position of the body.
	rendr.DrawCircle(int32(body.Position.X), int32(body.Position.Y), 1, 0, polygon.Color)
}

func (polygon *Polygon) UpdateVertices(position vector.Vec2, rotation float64) {
	for i := 0; i < len(polygon.WorldVertices); i++ {
		polygon.WorldVertices[i] = polygon.LocalVertices[i].Rotate(rotation)
		polygon.WorldVertices[i] = polygon.WorldVertices[i].Add(position)
	}
}

func (polygon *Polygon) MarkDebug() {
	polygon.Color = renderer.DEBUG
}

func (polygon *Polygon) UnMarkDebug() {
	polygon.Color = renderer.WHITE
}
