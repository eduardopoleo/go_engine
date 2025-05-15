package game

import (
	"engine/collision"
	"engine/entities"
	"engine/physics"
	"engine/vector"
)

type World struct {
	Bodies  []*entities.Body
	Forces  []*vector.Vec2
	Torques []*vector.Vec2
}

func (world *World) AddBody(body *entities.Body) {
	world.Bodies = append(world.Bodies, body)
}

func (world *World) AddForce(force *vector.Vec2) {
	world.Forces = append(world.Forces, force)
}

func (world *World) AddTorque(torque *vector.Vec2) {
	world.Forces = append(world.Forces, torque)
}

func (world *World) Update(dt float64) {
	for _, body := range world.Bodies {
		weight := physics.NewWeightForce(body.Mass)
		body.SumForces = body.SumForces.Add(weight)

		for _, force := range world.Forces {
			body.SumForces = body.SumForces.Add(*force)
		}

		for _, torque := range world.Torques {
			body.SumForces = body.SumForces.Add(*torque)
		}

		body.Update(dt)
	}
}

func (world *World) HandleCollisions() {
	// Check and resolve collisions for all bodies
	for a := 0; a < len(world.Bodies)-1; a++ {
		for b := a + 1; b < len(world.Bodies); b++ {
			bodyA := &world.Bodies[a]
			bodyB := &world.Bodies[b]
			collision.Resolve(*bodyA, *bodyB)
		}
	}
}
