package entities

import (
	"engine/renderer"
)

type Shape interface {
	MomentOfInertia() float64
	Draw(body *Body, renderer *renderer.Renderer)
	MarkDebug()
	UnMarkDebug()
}
