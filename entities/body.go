package entities

import (
	"engine/vector"
	"fmt"
)

type Body struct {
	// Linear properties
	Mass         float64
	Position     vector.Vec2
	Velocity     vector.Vec2
	Acceleration vector.Vec2
	SumForces    vector.Vec2

	// Angular properties
	Shape               Shape
	Rotation            float64
	AngularVelocity     float64
	AngularAcceleration float64
	SumTorque           float64
}

func (body *Body) IntegrateLinear(dt float64) {
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

func (body *Body) IntegrateAngular(dt float64) {
	if body.Mass == 0 {
		return
	}

	body.AngularAcceleration = body.SumTorque * (1 / (body.Shape.MomentOfInertia() * body.Mass))
	body.AngularVelocity += body.AngularAcceleration * dt
	body.Rotation += body.AngularVelocity * dt
	fmt.Printf("rotations %f\n", body.Rotation)
	body.SumTorque = 0
}
