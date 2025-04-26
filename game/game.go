package game

import (
	"engine/collision"
	"engine/constants"
	"engine/entities"
	"engine/physics"
	"engine/renderer"
	"engine/vector"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	Renderer            renderer.Renderer
	Bodies              []entities.Body
	PushForce           vector.Vec2
	TimeToPreviousFrame uint64
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	rendr := renderer.NewRenderer(name, width, height)
	game.Renderer = rendr
	game.PushForce = vector.Vec2{}
	game.TimeToPreviousFrame = sdl.GetTicks64()

	game.Bodies = append(game.Bodies, entities.NewBoxBody(
		renderer.WHITE, 100, 100, 2, vector.Vec2{X: 300, Y: 300}, 0,
	))

	game.Bodies = append(game.Bodies, entities.NewBoxBody(
		renderer.WHITE, 100, 100, 2, vector.Vec2{X: 300, Y: 500}, 2.3,
	))

	return game
}

func (game *Game) Input() {
	for event := game.Renderer.PollEvent(); event != nil; event = game.Renderer.PollEvent() {
		switch event.Type {
		case renderer.QUIT:
			println("Quit")
			game.Running = false
			game.Renderer.Destroy()
		case renderer.KEYDOWN:
			if event.Key() == renderer.ESCAPE {
				game.Running = false
				game.Renderer.Destroy()
			} else if event.Key() == renderer.LEFT_ARROW {
				game.PushForce.X = float64(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = float64(50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = float64(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = float64(50 * constants.PIXEL_PER_METER)
			}
		case renderer.KEYUP:
			if event.Key() == renderer.LEFT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = 0
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = 0
			}
		case renderer.MOUSEMOTION:
			x, y, _ := sdl.GetMouseState()
			box2 := &game.Bodies[1]
			box2.Position = vector.Vec2{X: float64(x), Y: float64(y)}

			if event.Key() == renderer.LEFT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = 0
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = 0
			}
		}
	}
}

func (game *Game) Update() {
	/*
		It forces the execution to be about 60 frames / second
	*/
	timeElapsed := sdl.GetTicks64() - game.TimeToPreviousFrame

	if constants.MILLISECONDS_PER_FRAME > timeElapsed {
		sdl.Delay(uint32(constants.MILLISECONDS_PER_FRAME - timeElapsed))
	}

	currentTime := sdl.GetTicks64()
	deltaTime := float64(currentTime-game.TimeToPreviousFrame) / 1000.0

	// Cap the maximum delta time to prevent instability
	if deltaTime > 0.016 {
		deltaTime = 0.016
	}

	game.TimeToPreviousFrame = sdl.GetTicks64()

	windowWidth, windowHeight := game.Renderer.GetWindowSize()

	// Update other bodies

	for i := range game.Bodies {
		body := &game.Bodies[i]
		if _, ok := body.Shape.(*entities.Circle); ok {
			weight := physics.NewWeightForce(body.Mass)
			body.SumForces = body.SumForces.Add(weight)
			body.SumForces = body.SumForces.Add(game.PushForce)
		}
	}

	// Check and resolve collisions for all bodies
	for a := 0; a < len(game.Bodies)-1; a++ {
		for b := a + 1; b < len(game.Bodies); b++ {
			bodyA := &game.Bodies[a]
			bodyB := &game.Bodies[b]
			polygonA, _ := bodyA.Shape.(*entities.Polygon)
			polygonB, _ := bodyB.Shape.(*entities.Polygon)

			collision := collision.CalculatePolygonPolygonCollision(bodyA, bodyB, polygonA, polygonB)

			// return &Collision{
			// 	BodyA:  *bodyA,
			// 	BodyB:  *bodyB,
			// 	Depth:  depth,
			// 	Normal: normal,
			// 	Start:  start,
			// 	End:    end,
			// }

			game.Renderer.ClearScreen()
			if collision != nil {
				fmt.Printf("got here!\n", collision)
				fmt.Printf("\n")
				game.Renderer.DrawFilledCircle(int32(collision.Start.X), int32(collision.Start.Y), 1, renderer.RED)
				game.Renderer.DrawFilledCircle(int32(collision.End.X), int32(collision.End.Y), 1, renderer.RED)

				drawEnd := collision.Start.Add(collision.Normal.Multiply(15))
				game.Renderer.DrawLine(collision.Start, drawEnd, renderer.RED)
			}
		}
	}

	// IntegrateLinear last all bodies not in between
	// use float64
	// use damping factor to increase stability
	for i := range game.Bodies {
		body := &game.Bodies[i]
		body.Update(deltaTime)
		bounce(body, game, windowWidth, windowHeight, deltaTime)
	}
}

func (game *Game) Draw() {

	for i := range game.Bodies {
		body := &game.Bodies[i]
		// fmt.Printf("coordinates for body i, %d X %f, Y %f\n", i, body.Position.X, body.Position.Y)
		body.Shape.Draw(body.Position.X, body.Position.Y, body.Rotation, &game.Renderer)
	}

	game.Renderer.Render()
}

func bounce(body *entities.Body, game *Game, windowWidth float64, windowHeight float64, deltaTime float64) {
	padding := 0.1
	if circle, ok := body.Shape.(*entities.Circle); ok {
		if (body.Position.X - float64(circle.Radius)) <= 0 {
			body.Velocity.X = -constants.RESTITUTION_COEFFICIENT * body.Velocity.X
			body.Position.X = float64(circle.Radius) + padding
		} else if (body.Position.X + float64(circle.Radius)) >= windowWidth {
			body.Velocity.X = -constants.RESTITUTION_COEFFICIENT * body.Velocity.X
			body.Position.X = windowWidth - float64(circle.Radius) - padding
		} else if (body.Position.Y - float64(circle.Radius)) <= 0 {
			body.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * body.Velocity.Y
			body.Position.Y = float64(circle.Radius) + padding
		} else if (body.Position.Y + float64(circle.Radius)) >= windowHeight {
			body.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * body.Velocity.Y
			body.Position.Y = windowHeight - float64(circle.Radius) - padding
		}
	}
}
