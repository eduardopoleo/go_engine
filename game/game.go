package game

import (
	"engine/collision"
	"engine/constants"
	"engine/entities"
	"engine/renderer"
	"engine/vector"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	DebugMode           bool
	Renderer            renderer.Renderer
	Bodies              []entities.Body
	PushForce           vector.Vec2
	TimeToPreviousFrame uint64
	Collisions          []*collision.Collision
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	rendr := renderer.NewRenderer(name, width, height)
	game.Renderer = rendr
	game.PushForce = vector.Vec2{}
	game.TimeToPreviousFrame = sdl.GetTicks64()

	// bottom := entities.NewBoxBody(
	// 	renderer.WHITE, float64(width-20), 50, 2, vector.Vec2{X: float64(width / 2), Y: float64(height - 20)}, 0, true,
	// )
	// left := entities.NewBoxBody(
	// 	renderer.WHITE, 50, float64(height-20), 2, vector.Vec2{X: 20, Y: float64(height / 2)}, 0, true,
	// )

	// right := entities.NewBoxBody(
	// 	renderer.WHITE, 50, float64(height-20), 2, vector.Vec2{X: float64(width - 20), Y: float64(height / 2)}, 0, true,
	// )

	bigBox := entities.NewBoxBody(
		renderer.WHITE, 150, 150, 2, vector.Vec2{X: float64(width / 2), Y: float64(height / 2)}, 0, true,
	)
	bigBox.Name = "bigBox"
	bigBox.Rotation = 1.4

	circle := entities.NewCircle(vector.Vec2{X: 400, Y: 600}, 40, renderer.WHITE, 2.0)

	// game.Bodies = append(game.Bodies, bottom)
	// game.Bodies = append(game.Bodies, left)
	// game.Bodies = append(game.Bodies, right)
	game.Bodies = append(game.Bodies, bigBox)
	game.Bodies = append(game.Bodies, circle)

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
			} else if event.Key() == renderer.D {
				game.DebugMode = !game.DebugMode
			}
		case renderer.MOUSEMOTION:
			x, y, _ := sdl.GetMouseState()
			circle := &game.Bodies[1]
			circle.Position.X = float64(x)
			circle.Position.Y = float64(y)
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

	// Update other bodies

	for i := range game.Bodies {
		body := &game.Bodies[i]
		// weight := physics.NewWeightForce(body.Mass)
		// body.SumForces = body.SumForces.Add(weight)
		body.SumForces = body.SumForces.Add(game.PushForce)
	}

	// Check and resolve collisions for all bodies
	for a := 0; a < len(game.Bodies)-1; a++ {
		for b := a + 1; b < len(game.Bodies); b++ {
			bodyA := &game.Bodies[a]
			bodyB := &game.Bodies[b]
			col := collision.Resolve(bodyA, bodyB)

			if game.DebugMode && col != nil {
				game.Collisions = append(game.Collisions, col)
				bodyA.Shape.MarkDebug()
				bodyB.Shape.MarkDebug()
			}
		}
	}

	// IntegrateLinear last all bodies not in between
	for i := range game.Bodies {
		body := &game.Bodies[i]
		body.Update(deltaTime)
	}
}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	if game.DebugMode {
		for _, col := range game.Collisions {
			collision.PolygonPolygonCollisionDebugger(col, game.Renderer)
		}

		game.Collisions = nil
	}

	for i := range game.Bodies {
		body := &game.Bodies[i]
		body.Shape.Draw(body, &game.Renderer)

		if game.DebugMode {
			body.Shape.UnMarkDebug()
		}
	}

	game.Renderer.Render()
}
