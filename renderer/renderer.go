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
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

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

func (renderer *Renderer) Destroy() {
	renderer.SDLWindow.Destroy()
	renderer.SDLRenderer.Destroy()
	sdl.Quit()
}

const (
	QUIT     string = "QUIT"
	KEYBOARD        = "KEYBOARD"
)

const (
	ESCAPE string = "ESCAPE"
)

type Event struct {
	Type          string
	OriginalEvent sdl.Event
}

func (renderer *Renderer) PollEvent() *Event {
	sdlEvent := sdl.PollEvent()

	// You can return nil if the function returns pointer
	if sdlEvent == nil {
		return nil
	}

	event := &Event{OriginalEvent: sdlEvent}

	switch sdlEvent.(type) {
	case *sdl.QuitEvent:
		event.Type = QUIT
	case *sdl.KeyboardEvent:
		event.Type = KEYBOARD
	}

	return event
}

func (event *Event) Key() string {
	if keyEvent, ok := event.OriginalEvent.(*sdl.KeyboardEvent); ok {
		switch keyEvent.Keysym.Sym {
		case sdl.K_ESCAPE:
			return ESCAPE
		}
	}
	return ""
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

// Private

func fromHex(hexColor uint32) sdl.Color {
	return sdl.Color{
		R: uint8((hexColor >> 24) & 0xFF),
		G: uint8((hexColor >> 16) & 0xFF),
		B: uint8((hexColor >> 8) & 0xFF),
		A: uint8(hexColor & 0xFF),
	}
}
