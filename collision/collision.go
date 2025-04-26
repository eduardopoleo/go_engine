package collision

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

	// resolvePenetration(collision, bodyA, bodyB)
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
	penetrationAB, edgeA, vertexA := calculatePenetration(polygonA, polygonB)
	penetrationBA, edgeB, vertexB := calculatePenetration(polygonB, polygonA)

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
		Highest value (less negative) means less penetration
	*/
	if penetrationAB >= penetrationBA {
		depth = -penetrationAB
		normal = edgeA.Normal()
		start = vertexA
		end = vertexA.Add(normal.Multiply(depth))
	} else {
		depth = -penetrationBA
		normal = edgeB.Normal()
		normal = normal.Multiply(-1)
		start = vertexB.Subtract(normal.Multiply(depth))
		end = vertexB
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
