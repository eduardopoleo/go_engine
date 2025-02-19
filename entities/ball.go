package entities

import (
	"engine/physics"
	"engine/renderer"
)

type Ball struct {
	Position physics.Vec2
	Radius   int32
	Color    uint32
}

func (ball *Ball) Render(renderer *renderer.Renderer) {
	renderer.DrawCircle(ball.Position.X, ball.Position.Y, ball.Radius, ball.Color)
}
