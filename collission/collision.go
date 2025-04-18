package collission

import (
	"engine/entities"
	"engine/vector"
	"fmt"
	"math"
)

type Collision struct {
	BodyA  entities.Body
	BodyB  entities.Body
	Normal vector.Vec2
	Start  vector.Vec2
	End    vector.Vec2
	Depth  float64
}

func ResolveCollision(bodyA *entities.Body, bodyB *entities.Body) {
	if bodyA.Static && bodyB.Static {
		return
	}

	circleA, isCircleA := bodyA.Shape.(*entities.Circle)
	circleB, isCircleB := bodyB.Shape.(*entities.Circle)

	var collision *Collision
	if isCircleA && isCircleB {
		collision = calculateCirCleCirCleCollission(bodyA, bodyB, circleA, circleB)
	}

	if collision == nil {
		return
	}

	resolvePenetration(collision, bodyA, bodyB)
	resolveImpulse(collision, bodyA, bodyB)
}

func calculateCirCleCirCleCollission(bodyA *entities.Body, bodyB *entities.Body, circleA *entities.Circle, circleB *entities.Circle) *Collision {
	d := bodyB.Position.Subtract(bodyA.Position)
	distanceAB := d.Magnitude()
	if distanceAB > float64((circleA.Radius + circleB.Radius)) {
		return nil
	}

	collisionNormal := d.Unit()

	start := bodyB.Position.Subtract(collisionNormal.Multiply(float64(circleB.Radius)))
	end := bodyA.Position.Add(collisionNormal.Multiply(float64(circleA.Radius)))
	dep := end.Subtract(start)
	depth := dep.Magnitude()

	return &Collision{
		BodyA:  *bodyA,
		BodyB:  *bodyB,
		Normal: collisionNormal,
		Start:  start,
		End:    end,
		Depth:  depth,
	}
}

func resolvePenetration(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
	fmt.Printf("penetration %f\n", collision.Depth)
	const minPenetration = 0.5
	if collision.Depth < minPenetration {
		return
	}

	invMassA := 1 / bodyA.Mass
	invMassB := 1 / bodyB.Mass
	invSum := invMassA + invMassB
	// Calculate the % of penetration
	da := collision.Depth / invSum * invMassA
	db := collision.Depth / invSum * invMassB

	// Apply to the bodies using the normal to transform the scalar into a vector.
	bodyA.Position = bodyA.Position.Subtract(collision.Normal.Multiply(da))
	bodyB.Position = bodyB.Position.Add(collision.Normal.Multiply(db))
}

func resolveImpulse(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
	e := math.Min(bodyA.E, bodyB.E)
	relativeVelocity := bodyA.Velocity.Subtract(bodyB.Velocity)
	velDotNormal := relativeVelocity.Dot(collision.Normal)

	invMassA := 1.0 / bodyA.Mass
	invMassB := 1.0 / bodyB.Mass

	invMassSum := invMassA + invMassB

	jMag := -(1 + e) * velDotNormal / invMassSum

	impulse := collision.Normal.Multiply(jMag)

	if !bodyA.Static {
		bodyA.Velocity = bodyA.Velocity.Add(impulse.Multiply(invMassA))
	}

	if !bodyB.Static {
		bodyB.Velocity = bodyB.Velocity.Subtract(impulse.Multiply(invMassB))
	}
}
