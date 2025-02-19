package game

import (
	"engine/entities"
	"engine/physics"
	"engine/renderer"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running  bool
	Renderer renderer.Renderer
	Balls    [3]entities.Ball
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer
	ball := entities.Ball{
		Position: physics.Vec2{X: 50, Y: 50},
		Radius:   5,
		Color:    0xFFFFFFFF,
	}
	game.Balls[0] = ball
	return game
}

func (game *Game) Input() {
	for event := game.Renderer.PollEvent(); event != nil; event = game.Renderer.PollEvent() {
		switch event.Type {
		case renderer.QUIT:
			println("Quit")
			game.Running = false
		case renderer.KEYBOARD:
			if event.Key() == renderer.ESCAPE {
				game.Running = false
				game.Renderer.Destroy()
			}
		}
	}
}

func (game *Game) Update() {
	ball := &game.Balls[0]
	increase := physics.Vec2{X: 5, Y: 5}

	ball.Position.Add(increase)
	// fmt.Printf("ball position X: %d, Y: %d\n", ball.Position.X, ball.Position.Y)
}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	ball := &game.Balls[0]
	ball.Render(&game.Renderer)

	game.Renderer.Render()
	sdl.Delay(16)
}
