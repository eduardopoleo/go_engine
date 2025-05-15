package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Shape interface {
	MomentOfInertia() float64
	Draw(body *Body, renderer *renderer.Renderer)
	MarkDebug()
	UnMarkDebug()
	GetHeight() float64
	GetWidth() float64
	UpdateVertices(position vector.Vec2, rotation float64)
}
