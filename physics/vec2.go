package physics

type Vec2 struct {
	X int32
	Y int32
}

func (vec *Vec2) Add(otherVec Vec2) *Vec2 {
	vec.X += otherVec.X
	vec.Y += otherVec.Y

	return vec
}
