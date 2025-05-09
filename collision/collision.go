package collision

import (
	"engine/entities"
	"engine/renderer"
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
	if bodyA.Static && bodyB.Static {
		return
	}

	circleA, isCircleA := bodyA.Shape.(*entities.Circle)
	circleB, isCircleB := bodyB.Shape.(*entities.Circle)

	polygonA, isPolygonA := bodyA.Shape.(*entities.Polygon)
	polygonB, isPolygonB := bodyB.Shape.(*entities.Polygon)

	var collision *Collision
	if isCircleA && isCircleB {
		collision = calculateCirCleCirCleCollission(bodyA, bodyB, circleA, circleB)
	}

	if isPolygonA && isPolygonB {
		collision = CalculatePolygonPolygonCollision(bodyA, bodyB, polygonA, polygonB)
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

func CalculatePolygonPolygonCollision(bodyA *entities.Body, bodyB *entities.Body, polygonA *entities.Polygon, polygonB *entities.Polygon) *Collision {
	penetrationAB, edgeA, vertexB := calculatePenetration(polygonA, polygonB)
	penetrationBA, edgeB, vertexA := calculatePenetration(polygonB, polygonA)

	/*
		If the max penetration for both is positive it means that there was no collision
	*/
	if penetrationBA >= 0 || penetrationAB >= 0 {
		return nil
	}

	var depth float64
	var normal vector.Vec2
	var start vector.Vec2
	var end vector.Vec2
	/*
		We want to resolve the smallest penetration in this case the highest value.
		Highest value (less negative) means less penetration.

		depth needs to be positive because both penetrationAB penetrationBA are negative by definition are negative we need to multiply by -
	*/
	if penetrationAB >= penetrationBA {
		depth = -penetrationAB
		normal = edgeA.Normal()
		start = vertexB
		end = vertexB.Add(normal.Multiply(depth))
	} else {
		depth = -penetrationBA
		edgeBNormal := edgeB.Normal()
		normal = edgeBNormal.Multiply(-1) // We need to go from A->B
		start = vertexA.Add(edgeBNormal.Multiply(depth))
		end = vertexA
	}

	return &Collision{
		BodyA:  *bodyA,
		BodyB:  *bodyB,
		Depth:  depth,
		Normal: normal,
		Start:  start,
		End:    end,
	}
}

func calculatePenetration(polygonA *entities.Polygon, polygonB *entities.Polygon) (float64, vector.Vec2, vector.Vec2) {
	penetration := float64(math.MinInt)
	var collidingVertex vector.Vec2
	var collidingEdge vector.Vec2

	for idx, vertexA := range polygonA.WorldVertices {
		edge := polygonA.EdgeAt(idx)
		normal := edge.Normal()

		minPenetration := float64(math.MaxInt)
		var minVertex vector.Vec2
		for _, vertexB := range polygonB.WorldVertices {
			vba := vertexB.Subtract(vertexA)
			currPenetration := vba.Dot(normal)
			/*
				The smallest (or most negative) value is the closest point or the point that
				has penetrated the most looking through this edge
			*/
			if currPenetration < minPenetration {
				minPenetration = currPenetration
				minVertex = vertexB
			}
		}

		if minPenetration > penetration {
			penetration = minPenetration
			collidingEdge = edge
			collidingVertex = minVertex
		}
	}

	return penetration, collidingEdge, collidingVertex
}

func resolvePenetration(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
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

// func resolveImpulse(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
// 	e := math.Min(bodyA.E, bodyB.E)
// 	relativeVelocity := bodyA.Velocity.Subtract(bodyB.Velocity)
// 	velDotNormal := relativeVelocity.Dot(collision.Normal)

// 	invMassA := 1.0 / bodyA.Mass
// 	invMassB := 1.0 / bodyB.Mass

// 	invMassSum := invMassA + invMassB

// 	jMag := -(1 + e) * velDotNormal / invMassSum

// 	impulse := collision.Normal.Multiply(jMag)

// 	if !bodyA.Static {
// 		bodyA.Velocity = bodyA.Velocity.Add(impulse.Multiply(invMassA))
// 	}

// 	if !bodyB.Static {
// 		bodyB.Velocity = bodyB.Velocity.Subtract(impulse.Multiply(invMassB))
// 	}
// }

func resolveImpulse(collision *Collision, bodyA *entities.Body, bodyB *entities.Body) {
	e := math.Min(bodyA.E, bodyB.E)

	// r is the distance from the center of mass to the point of collision aprox.
	ra := collision.End.Subtract(bodyA.Position)
	rb := collision.Start.Subtract(bodyB.Position)

	// V = v + w X r at the point of contact determined by r
	Va := bodyA.Velocity.Add(bodyA.AngularVelocityProduct(ra))
	Vb := bodyB.Velocity.Add(bodyA.AngularVelocityProduct(rb))

	// vrel = Va - Vb
	normal := collision.Normal
	vRel := Va.Subtract(Vb)
	vRelNormal := vRel.Dot(normal)

	/*
	                       -(1 + e)(Vrel n)
	   Jn =  ---------------------------------------------
	            1        1     (ra X n)^2     (rb x n)^2
	         ------- + ----- + ----------- + ------------
	            Ma       Mb        Ia             Ib
	*/
	num := -(1 + e) * vRelNormal
	linearDen := bodyA.InvMass + bodyB.InvMass
	AngularDenA := ra.Cross(normal) * ra.Cross(normal) / bodyA.Shape.MomentOfInertia()
	AngularDenB := rb.Cross(normal) * rb.Cross(normal) / bodyB.Shape.MomentOfInertia()

	J := num / (linearDen + AngularDenA + AngularDenB)
	Jn := normal.Multiply(J)

	bodyA.ApplyImpulse(Jn, ra)
	bodyB.ApplyImpulse(Jn, rb)

	// f := math.Min(bodyA.F, bodyB.F)
}

func PolygonPolygonCollisionDebugger(collision *Collision, rend renderer.Renderer) {
	if collision != nil {
		rend.DrawFilledCircle(int32(collision.Start.X), int32(collision.Start.Y), 2, renderer.RED)
		rend.DrawFilledCircle(int32(collision.End.X), int32(collision.End.Y), 2, renderer.RED)

		drawEnd := collision.Start.Add(collision.Normal.Multiply(15))
		rend.DrawLine(collision.Start, drawEnd, renderer.RED)
	}
}
