package physics

import (
	"engine/constants"
	"engine/entities"
	"engine/vector"
)

func NewWeightForce(mass float32) vector.Vec2 {
	return vector.Vec2{X: 0, Y: mass * constants.GRAVITY}
}

func NewSpringForce(particle *entities.Particle, anchor *entities.Particle) vector.Vec2 {
	position := particle.Position.Subtract(anchor.Position)
	forceDirection := position.Unit()
	displacement := position.Magnitude() - constants.SPRING_REST_LENGTH
	forceMagnitude := -1 * constants.SPRING_CONSTANT * displacement

	return forceDirection.Multiply(forceMagnitude)
}
