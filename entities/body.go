package entities

import (
	"engine/vector"
)

type Body struct {
	Static bool

	// Linear properties
	Mass         float64
	Position     vector.Vec2
	Velocity     vector.Vec2
	Acceleration vector.Vec2
	SumForces    vector.Vec2

	// Angular properties
	Shape               Shape
	Rotation            float64
	AngularVelocity     float64
	AngularAcceleration float64
	SumTorque           float64

	// Impulse
	E float64 // coefficient of restituation
}

func (body *Body) IntegrateLinear(dt float64) {
	if body.Mass == 0 || body.Static {
		return
	}

	body.Acceleration = body.SumForces.Multiply(1.0 / body.Mass)

	// Update velocity first (semi-implicit Euler)
	dampingFactor := 0.99
	body.Velocity = body.Velocity.Add(body.Acceleration.Multiply(dt))
	body.Velocity = body.Velocity.Multiply(dampingFactor)

	// Then update position
	body.Position = body.Position.Add(body.Velocity.Multiply(dt))

	body.SumForces = vector.Vec2{X: 0, Y: 0}
}

func NewBoxBody(color uint32, width float64, height float64, mass float64, position vector.Vec2, rotation float64) Body {
	newBoxShape := NewBox(color, width, height)
	box := Body{
		Position: position,
		Mass:     mass,
		Shape:    newBoxShape,
		Rotation: rotation,
	}
	box.Shape = newBoxShape
	return box
}

func (body *Body) IntegrateAngular(dt float64) {
	if body.Mass == 0 || body.Static {
		return
	}

	body.AngularAcceleration = body.SumTorque * (1 / (body.Shape.MomentOfInertia() * body.Mass))
	body.AngularVelocity += body.AngularAcceleration * dt
	body.Rotation += body.AngularVelocity * dt
	body.SumTorque = 0
}

func (body *Body) Update(dt float64) {
	body.IntegrateLinear(dt)
	body.IntegrateAngular(dt)

	if box, ok := body.Shape.(*Polygon); ok {
		for i := 0; i < len(box.WorldVertices); i++ {
			box.WorldVertices[i] = box.LocalVertices[i].Rotate(body.Rotation)
			box.WorldVertices[i] = box.WorldVertices[i].Add(body.Position)
		}
	}
}
