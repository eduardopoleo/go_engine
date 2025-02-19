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

// Events type
const (
	QUIT     string = "QUIT"
	KEYBOARD string = "KEYBOARD"
	KEYDOWN  string = "KEYDOWN"
	KEYUP    string = "KEYUP"
)

// Keys
const (
	ESCAPE      string = "ESCAPE"
	LEFT_ARROW  string = "LEFT_ARROW"
	RIGHT_ARROW string = "RIGHT_ARROW"
	UP_ARROW    string = "UP_ARROW"
	DOWN_ARROW  string = "DOWN_ARROW"
)

type Event struct {
	Type          string
	OriginalEvent sdl.Event
}

func (renderer *Renderer) PollEvent() *Event {
	sdlEvent := sdl.PollEvent()
	if sdlEvent == nil {
		return nil // <- Prevent infinite loop when no more events
	}

	event := &Event{OriginalEvent: sdlEvent} // Use pointer directly

	if _, ok := sdlEvent.(*sdl.QuitEvent); ok {
		event.Type = QUIT
	} else if keyEvent, ok := sdlEvent.(*sdl.KeyboardEvent); ok {
		if keyEvent.Type == sdl.KEYDOWN {
			event.Type = KEYDOWN
		} else if keyEvent.Type == sdl.KEYUP {
			event.Type = KEYUP
		}
	}
	return event
}

func (event *Event) Key() string {
	if keyEvent, ok := event.OriginalEvent.(*sdl.KeyboardEvent); ok {
		switch keyEvent.Keysym.Sym {
		case sdl.K_ESCAPE:
			return ESCAPE
		case sdl.K_LEFT:
			return LEFT_ARROW
		case sdl.K_RIGHT:
			return RIGHT_ARROW
		case sdl.K_UP:
			return UP_ARROW
		case sdl.K_DOWN:
			return DOWN_ARROW
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
