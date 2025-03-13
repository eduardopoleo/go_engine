package physics

import (
	"engine/constants"
)

func NewGravityForce(mass float32) Vec2 {
	return Vec2{X: 0, Y: mass * constants.GRAVITY}
}
