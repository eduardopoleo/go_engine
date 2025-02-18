package main

import (
	"engine/game"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, renderer, error := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_BORDERLESS)

	if error != nil {
		panic("oh no!")
	}

	if window == nil {
		panic("ho no, no window")
	}

	color := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	gfx.LineColor(renderer, 300, 300, 500, 500, color)

	game := game.Game{Running: true}
	game.Input()
}
