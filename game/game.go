package game

import (
	"engine/renderer"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running  bool
	Renderer renderer.Renderer
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer

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
	// Todo in here
}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()
	game.Renderer.DrawCircle(300, 300, 100, 0xFFFFFFFF)
	game.Renderer.Render()
	sdl.Delay(33)
}
