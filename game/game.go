package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running bool
}

func (game *Game) Input() {
	for game.Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				game.Running = false
			case *sdl.KeyboardEvent:
				fmt.Println("hello there!")
				if e.Keysym.Sym == sdl.K_ESCAPE {
					game.Running = false
				}
			}
		}
	}
}
