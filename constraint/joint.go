package constraint

import (
	"engine/entities"
	"engine/vector"
)

type JointConstraint struct {
	A *entities.Body
	B *entities.Body

	AnchorA vector.Vec2 // anchor point in local space
	AnchorB vector.Vec2
}

/*

[
		2 (anchorA - anchorB),
		2 (ra X (ral - rb)),
	  2 * (rb - ra),
		2 * (rbl X (rb - ra))
	]


*/

func (constraint *JointConstraint) Solve() {

}
