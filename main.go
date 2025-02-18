package main

import (
	"engine/game"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// Find a better place to put this behind!
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	defer sdl.Quit()

	window, err := sdl.CreateWindow("Test RenderGeometryRaw", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Destroy()

	game := game.Game{Running: true}

	for game.Running {
		game.Input()
		game.Update()
		game.Draw(renderer)

		sdl.Delay(16)
	}
}
