package physics

type Vec2 struct {
	X float32
	Y float32
}

func (vec *Vec2) Add(otherVec Vec2) Vec2 {
	return Vec2{
		X: vec.X + otherVec.X,
		Y: vec.Y + otherVec.Y,
	}
}

func (vec *Vec2) Multiply(n float32) Vec2 {
	return Vec2{
		X: vec.X * n,
		Y: vec.Y * n,
	}
}
