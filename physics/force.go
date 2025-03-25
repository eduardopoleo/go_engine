package physics

import (
	"engine/constants"
	"engine/entities"
	"engine/vector"
)

func NewWeightForce(mass float64) vector.Vec2 {
	return vector.Vec2{X: 0, Y: mass * constants.GRAVITY}
}

func NewSpringForce(body *entities.Body, anchor *entities.Body) vector.Vec2 {
	position := body.Position.Subtract(anchor.Position)
	forceDirection := position.Unit()
	displacement := position.Magnitude() - constants.SPRING_REST_LENGTH
	forceMagnitude := -1 * constants.SPRING_CONSTANT * displacement

	return forceDirection.Multiply(forceMagnitude)
}
