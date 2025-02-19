package entities

import (
	"engine/physics"
	"engine/renderer"
)

type Particle struct {
	Mass         float32
	Radius       int32
	Color        uint32
	Position     physics.Vec2
	Velocity     physics.Vec2
	Acceleration physics.Vec2
	SumForces    physics.Vec2
}

func (particle *Particle) Render(renderer *renderer.Renderer) {
	renderer.DrawCircle(
		int32(particle.Position.X),
		int32(particle.Position.Y),
		particle.Radius,
		particle.Color,
	)
}

func (particle *Particle) Integrate(dt float32) {
	if particle.Mass == 0 {
		return
	}
	/*
		F = M * A
		Forces are applied on every frame.
		Consequently, the acceleration is set anew on every frame
	*/
	particle.Acceleration = particle.SumForces.Multiply(1.0 / particle.Mass)
	// V = A * dt

	/*
		The velocity and position stay the same and are only affected when there are forces
		acting on the particle
	*/
	particle.Velocity = particle.Velocity.Add(particle.Acceleration.Multiply(dt))
	// P = V * dt
	particle.Position = particle.Position.Add(particle.Velocity.Multiply(dt))
	// Clear the sum of forces so that the next
	particle.SumForces = physics.Vec2{X: 0, Y: 0}
}
