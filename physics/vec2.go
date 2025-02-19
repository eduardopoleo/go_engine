package physics

type Vec2 struct {
	X float32
	Y float32
}

func (vec *Vec2) Add(otherVec Vec2) *Vec2 {
	vec.X += otherVec.X
	vec.Y += otherVec.Y

	return vec
}

func (vec *Vec2) Multiply(n float32) *Vec2 {
	vec.X *= n
	vec.Y *= n

	return vec
}
