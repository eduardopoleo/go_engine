package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Particle struct {
	Mass         float64
	Radius       int32
	Color        uint32
	Position     vector.Vec2
	Velocity     vector.Vec2
	Acceleration vector.Vec2
	SumForces    vector.Vec2
}

func (particle *Particle) Render(renderer *renderer.Renderer) {
	renderer.DrawCircle(
		int32(particle.Position.X),
		int32(particle.Position.Y),
		particle.Radius,
		particle.Color,
	)
}

func (particle *Particle) Integrate(dt float64) {
	if particle.Mass == 0 {
		return
	}

	particle.Acceleration = particle.SumForces.Multiply(1.0 / particle.Mass)

	// Update velocity first (semi-implicit Euler)
	dampingFactor := 0.99
	particle.Velocity = particle.Velocity.Add(particle.Acceleration.Multiply(dt))
	particle.Velocity = particle.Velocity.Multiply(dampingFactor)

	// Then update position
	particle.Position = particle.Position.Add(particle.Velocity.Multiply(dt))

	particle.SumForces = vector.Vec2{X: 0, Y: 0}
}

/*
	pressing an arrow key enacts force
	force turns into acceleration
	which turns into velocity

	releasing the key releases the force
	then the acceleration is zero
	the velocity right now stays the same cuz I do not have friction yet
*/
