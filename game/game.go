package game

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running      bool
	WindowWidth  int
	WindowHeight int
}

func (game *Game) Input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			game.Running = false
		case *sdl.KeyboardEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				game.Running = false
			}
		}
	}
}

func (game *Game) Update() {
	// Todo in here
}

func (game *Game) Draw(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
	renderer.Clear()

	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	gfx.FilledCircleColor(renderer, 300, 300, 100, color)

	renderer.Present()
	sdl.Delay(33)
}
