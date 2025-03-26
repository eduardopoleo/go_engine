package entities

import (
	"engine/vector"
)

type Body struct {
	// General fields
	Mass     float64
	Position vector.Vec2
	Shape    Shape
	// Linear properties
	Velocity     vector.Vec2
	Acceleration vector.Vec2
	SumForces    vector.Vec2
	// Angular properties
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
