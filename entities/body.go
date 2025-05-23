package entities

import (
	"engine/renderer"
	"engine/vector"
)

type Body struct {
	Static  bool
	Name    string
	Texture *renderer.SDLTexture

	// Linear properties
	Mass         float64
	InvMass      float64
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
	E float64 // coefficient of restitution
	F float64 // coefficient of friction
}

func (body *Body) AttachTexture(path string, rend *renderer.Renderer) {
	texture := renderer.LoadTexture(path, rend)
	body.Texture = texture
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

func (body *Body) AngularVelocityProduct(r vector.Vec2) vector.Vec2 {
	/*
		Cross product result in two dimentions
		w = [0,  0,  w]
		r = [rx, ry, 0]
		w X r = [
			X: - w * ry
			Y: w * rx
		]
	*/
	return vector.Vec2{X: -body.AngularVelocity * r.Y, Y: body.AngularVelocity * r.X}
}

func NewCircle(position vector.Vec2, radius int32, color uint32, mass float64) Body {
	circleShape := CircleShape(radius, color)
	circle := Body{
		Position: position,
		Mass:     mass,
		InvMass:  1 / mass,
		Shape:    circleShape,
		Rotation: 0,
		Static:   false,
		E:        1,
		F:        1,
		Name:     "Circle",
	}
	return circle
}

func NewBoxBody(color uint32, width float64, height float64, mass float64, position vector.Vec2, rotation float64, static bool) Body {
	newBoxShape := NewBox(color, width, height)
	box := Body{
		Position: position,
		Mass:     mass,
		InvMass:  1 / mass,
		Shape:    newBoxShape,
		Rotation: rotation,
		Static:   static,
		E:        1,
		F:        1,
	}
	return box
}

func (body *Body) ApplyImpulse(J vector.Vec2, r vector.Vec2) {
	body.applyLinearImpulse(J)
	body.applyAngularImpulse(J, r)
}

func (body *Body) applyLinearImpulse(J vector.Vec2) {
	body.Velocity = body.Velocity.Add(J.Multiply(body.InvMass))
}

func (body *Body) applyAngularImpulse(J vector.Vec2, r vector.Vec2) {
	body.AngularVelocity = body.AngularVelocity + r.Cross(J)/body.Shape.MomentOfInertia()
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
	body.Shape.UpdateVertices(body.Position, body.Rotation)
}

func (body *Body) Destroy() {
	if body.Texture != nil {
		println("Cleaning up texture")
		body.Texture.Destroy()
	}
}
