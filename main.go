package main

import (
	"engine/game"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	game := game.NewGame("Awesome new game", 600, 800)

	for game.Running {
		game.Input()
		game.Update()
		game.Draw()

		sdl.Delay(16)
	}
}
