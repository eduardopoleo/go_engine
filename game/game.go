package game

import (
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
	TimeToPreviousFrame uint64
	World               World
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	rendr := renderer.NewRenderer(name, width, height)
	game.Renderer = rendr
	game.TimeToPreviousFrame = sdl.GetTicks64()

	bottom := entities.NewBoxBody(
		renderer.WHITE, float64(width-20), 50, 2, vector.Vec2{X: float64(width / 2), Y: float64(height - 20)}, 0, true,
	)
	left := entities.NewBoxBody(
		renderer.WHITE, 50, float64(height-20), 2, vector.Vec2{X: 20, Y: float64(height / 2)}, 0, true,
	)

	right := entities.NewBoxBody(
		renderer.WHITE, 50, float64(height-20), 2, vector.Vec2{X: float64(width - 20), Y: float64(height / 2)}, 0, true,
	)

	bigBox := entities.NewBoxBody(
		renderer.WHITE, 150, 150, 2, vector.Vec2{X: float64(width / 2), Y: float64(height / 2)}, 0, true,
	)
	bigBox.Name = "bigBox"
	bigBox.Rotation = 1.4
	bigBox.AttachTexture("./assets/crate.png", &game.Renderer)

	game.World.AddBody(&bottom)
	game.World.AddBody(&left)
	game.World.AddBody(&right)
	game.World.AddBody(&bigBox)

	return game
}

func (game *Game) Input() {
	for event := game.Renderer.PollEvent(); event != nil; event = game.Renderer.PollEvent() {
		switch event.Type {
		case renderer.QUIT:
			game.Running = false
			game.Cleanup()
		case renderer.KEYDOWN:
			if event.Key() == renderer.ESCAPE {
				game.Running = false
				game.Cleanup()
			}
		case renderer.MOUSE_BUTTON_LEFT_UP:
			x, y, _ := sdl.GetMouseState()
			circle := entities.NewCircle(
				vector.Vec2{X: float64(x), Y: float64(y)},
				25,
				renderer.WHITE,
				2,
			)
			circle.AttachTexture("./assets/bowlingball.png", &game.Renderer)
			game.World.AddBody(&circle)
		case renderer.MOUSE_BUTTON_RIGHT_UP:
			x, y, _ := sdl.GetMouseState()
			polygonShape := entities.NewBox(renderer.WHITE, 50, 50)
			polygon := entities.Body{
				Position: vector.Vec2{X: float64(x), Y: float64(y)},
				Mass:     2.0,
				InvMass:  float64(1 / 2.0),
				Shape:    polygonShape,
				Rotation: 0.7,
				Static:   false,
				E:        0.1,
				F:        0.7,
				Name:     "Polygon",
			}
			polygon.AttachTexture("./assets/crate.png", &game.Renderer)
			game.World.AddBody(&polygon)
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

	game.World.Update(deltaTime)
	game.World.HandleCollisions()
}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	for i := range game.World.Bodies {
		body := game.World.Bodies[i]
		if body.Texture == nil {
			body.Shape.Draw(body, &game.Renderer)

			if game.DebugMode {
				body.Shape.UnMarkDebug()
			}
		} else {
			body.Texture.Draw(
				body.Position.X,
				body.Position.Y,
				body.Rotation,
				body.Shape.GetWidth(),
				body.Shape.GetHeight(),
				&game.Renderer,
			)
		}
	}

	game.Renderer.Render()
}

func (game *Game) Cleanup() {
	for _, body := range game.World.Bodies {
		body.Destroy()
	}
	game.Renderer.Destroy()
}
