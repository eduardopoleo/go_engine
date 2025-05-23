package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Circle struct {
	Color  uint32
	Radius int32
}

func CircleShape(radius int32, color uint32) *Circle {
	return &Circle{
		Color:  color,
		Radius: radius,
	}
}

func (circle *Circle) MomentOfInertia() float64 {
	return 0.5 * float64(circle.Radius) * float64(circle.Radius)
}

func (circle *Circle) GetHeight() float64 {
	return 2 * float64(circle.Radius)
}

func (circle *Circle) GetWidth() float64 {
	return 2 * float64(circle.Radius)
}

func (circle *Circle) Draw(body *Body, renderer *renderer.Renderer) {
	// fmt.Printf("drawing x %f, y %f\n", x, y)
	renderer.DrawCircle(
		int32(body.Position.X),
		int32(body.Position.Y),
		circle.Radius,
		body.Rotation,
		circle.Color,
	)
}

func (circle *Circle) UpdateVertices(position vector.Vec2, rotation float64) {
	return // no-op
}

func (circle *Circle) MarkDebug() {
	circle.Color = renderer.DEBUG
}

func (circle *Circle) UnMarkDebug() {
	circle.Color = renderer.WHITE
}
