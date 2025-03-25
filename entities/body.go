package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Body struct {
	Mass         float64
	Radius       int32
	Color        uint32
	Position     vector.Vec2
	Velocity     vector.Vec2
	Acceleration vector.Vec2
	SumForces    vector.Vec2
}

func (body *Body) Render(renderer *renderer.Renderer) {
	renderer.DrawCircle(
		int32(body.Position.X),
		int32(body.Position.Y),
		body.Radius,
		body.Color,
	)
}

func (body *Body) Integrate(dt float64) {
	if body.Mass == 0 {
		return
	}

	body.Acceleration = body.SumForces.Multiply(1.0 / body.Mass)

	// Update velocity first (semi-implicit Euler)
	dampingFactor := 0.99
	body.Velocity = body.Velocity.Add(body.Acceleration.Multiply(dt))
	body.Velocity = body.Velocity.Multiply(dampingFactor)

	// Then update position
	body.Position = body.Position.Add(body.Velocity.Multiply(dt))

	body.SumForces = vector.Vec2{X: 0, Y: 0}
}

/*
	pressing an arrow key enacts force
	force turns into acceleration
	which turns into velocity

	releasing the key releases the force
	then the acceleration is zero
	the velocity right now stays the same cuz I do not have friction yet
*/
