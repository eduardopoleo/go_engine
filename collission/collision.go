package collission

import (
	"engine/entities"
	"engine/vector"
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
	// resolveImpulse(collision, bodyA, bodyB)
}

func calculateCirCleCirCleCollission(bodyA *entities.Body, bodyB *entities.Body, circleA *entities.Circle, circleB *entities.Circle) *Collision {
	d := bodyB.Position.Subtract(bodyA.Position)
	distanceAB := d.Magnitude()
	if distanceAB > float64((circleA.Radius + circleB.Radius)) {
		return nil
	}

	collisionNormal := d.Unit()

	start := bodyB.Position.Subtract(collisionNormal.Multiply(float64(circleB.Radius)))
	end := bodyA.Position.Subtract(collisionNormal.Multiply(float64(circleA.Radius)))
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
	totalMass := bodyA.Mass + bodyB.Mass
	// Calculate the % of penetration
	da := collision.Depth * bodyA.Mass / totalMass
	db := collision.Depth * bodyB.Mass / totalMass

	// Apply to the bodies using the normal to transform the scalar into a vector.
	bodyA.Position = bodyA.Position.Subtract(collision.Normal.Multiply(da))
	bodyB.Position = bodyB.Position.Add(collision.Normal.Multiply(db))
}

func resolveImpulse(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
	E := math.Min(bodyA.E, bodyB.E)
	velRel := bodyA.Velocity.Subtract(bodyB.Velocity)
	jMag := -(1 + E) * velRel.Dot(collision.Normal) / ((1 / bodyA.Mass) + (1 / bodyB.Mass))
	jn := collision.Normal.Multiply(jMag)

	bodyA.Velocity = bodyA.Velocity.Add(jn)
	bodyB.Velocity = bodyB.Velocity.Add(jn.Multiply(-1))
}
