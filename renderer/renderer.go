package renderer

import (
	"log"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	SDLRenderer *sdl.Renderer
	SDLWindow   *sdl.Window
}

func NewRenderer(name string, width int32, height int32) Renderer {
	window, err := sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		log.Fatal(err)
	}

	return Renderer{SDLRenderer: renderer, SDLWindow: window}
}

func (renderer *Renderer) ClearScreen() {
	renderer.SDLRenderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
	renderer.SDLRenderer.Clear()
}

func (renderer *Renderer) Render() {
	renderer.SDLRenderer.Present()
}

func (renderer *Renderer) DrawCircle(x int32, y int32, radius int32, color uint32) {
	gfx.FilledCircleColor(renderer.SDLRenderer, x, y, radius, fromHex(color))
}

func fromHex(hexColor uint32) sdl.Color {
	return sdl.Color{
		R: uint8((hexColor >> 24) & 0xFF),
		G: uint8((hexColor >> 16) & 0xFF),
		B: uint8((hexColor >> 8) & 0xFF),
		A: uint8(hexColor & 0xFF),
	}
}
