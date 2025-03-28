package game

import (
	"engine/constants"
	"engine/entities"
	"engine/physics"
	"engine/renderer"
	"engine/vector"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	Renderer            renderer.Renderer
	Bodies              []entities.Body
	SpringAnchor        entities.Body
	SpringBodies        []entities.Body
	PushForce           vector.Vec2
	TimeToPreviousFrame uint64
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer
	game.PushForce = vector.Vec2{}
	game.TimeToPreviousFrame = sdl.GetTicks64()

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
		case renderer.MOUSE_UP_EVENT:
			if event.Key() == renderer.BUTTON_LEFT {
				mouseX, mouseY := game.Renderer.GetMouseCoordinates()
				body := entities.Body{
					Position: vector.Vec2{X: mouseX, Y: mouseY},
					Mass:     2.0,
				}
				body.Shape = entities.NewCircle(20, renderer.WHITE)
				game.Bodies = append(game.Bodies, body)
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
		weight := physics.NewWeightForce(body.Mass)
		body.SumForces = body.SumForces.Add(weight)
		body.SumForces = body.SumForces.Add(game.PushForce)
	}

	for i := 0; i < len(game.SpringBodies); i++ {
		body := &game.SpringBodies[i]
		weight := physics.NewWeightForce(body.Mass)
		body.SumForces = body.SumForces.Add(weight)
		body.SumForces = body.SumForces.Add(game.PushForce)
	}

	if len(game.SpringBodies) > 0 {
		springForce := physics.NewSpringForce(&game.SpringBodies[0], &game.SpringAnchor)
		game.SpringBodies[0].SumForces = game.SpringBodies[0].SumForces.Add(springForce)
	}

	// Apply spring forces between bodies
	for i := 1; i < len(game.SpringBodies); i++ {
		previousBody := &game.SpringBodies[i-1]
		currentBody := &game.SpringBodies[i]

		springForce := physics.NewSpringForce(currentBody, previousBody)
		currentBody.SumForces = currentBody.SumForces.Add(springForce)
		previousBody.SumForces = previousBody.SumForces.Add(springForce.Multiply(-1))
	}

	// Integrate last all bodies not in between
	// use float64
	// use damping factor to increase stability
	for i := range game.Bodies {
		body := &game.Bodies[i]
		body.Integrate(deltaTime)
		bounce(body, game, windowWidth, windowHeight, deltaTime)
	}

	for i := 0; i < len(game.SpringBodies); i++ {
		body := &game.SpringBodies[i]
		body.Integrate(deltaTime)
		bounce(body, game, windowWidth, windowHeight, deltaTime)
	}

}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	for i := range game.Bodies {
		body := &game.Bodies[i]
		body.Shape.Draw(body.Position.X, body.Position.Y, &game.Renderer)
	}

	game.Renderer.Render()
}

func bounce(body *entities.Body, game *Game, windowWidth float64, windowHeight float64, deltaTime float64) {
	if circle, ok := body.Shape.(*entities.Circle); ok {
		if (body.Position.X - float64(circle.Radius)) <= 0 {
			body.Velocity.X = -constants.RESTITUTION_COEFFICIENT * body.Velocity.X
			body.Position.X = float64(circle.Radius)
		} else if (body.Position.X + float64(circle.Radius)) >= windowWidth {
			body.Velocity.X = -constants.RESTITUTION_COEFFICIENT * body.Velocity.X
			body.Position.X = windowWidth - float64(circle.Radius)
		} else if (body.Position.Y - float64(circle.Radius)) <= 0 {
			body.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * body.Velocity.Y
			body.Position.Y = float64(circle.Radius)
		} else if (body.Position.Y + float64(circle.Radius)) >= windowHeight {
			body.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * body.Velocity.Y
			body.Position.Y = windowHeight - float64(circle.Radius)
		}
	}
}
