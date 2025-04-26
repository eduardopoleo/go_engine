package renderer

import (
	"engine/vector"
	"log"
	"math"

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

	window, err := sdl.CreateWindow(name,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width, height,
		sdl.WINDOW_SHOWN,
	)

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

// Colors

const (
	WHITE uint32 = 0xFFFFFFFF
)

// Events type
const (
	QUIT           string = "QUIT"
	KEYBOARD       string = "KEYBOARD"
	KEYDOWN        string = "KEYDOWN"
	KEYUP          string = "KEYUP"
	MOUSE_UP_EVENT string = "MOUSE_UP_EVENT"
	MOUSEMOTION    string = "MOUSEMOTION"
)

// Keys
const (
	ESCAPE      string = "ESCAPE"
	LEFT_ARROW  string = "LEFT_ARROW"
	RIGHT_ARROW string = "RIGHT_ARROW"
	UP_ARROW    string = "UP_ARROW"
	DOWN_ARROW  string = "DOWN_ARROW"
	BUTTON_LEFT string = "BUTTON_LEFT"
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
	} else if mouseEvent, ok := sdlEvent.(*sdl.MouseButtonEvent); ok {
		if mouseEvent.Type == sdl.MOUSEBUTTONUP {
			event.Type = MOUSE_UP_EVENT
		}
	} else if _, ok := sdlEvent.(*sdl.MouseMotionEvent); ok {
		event.Type = MOUSEMOTION
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
		case sdl.BUTTON_LEFT:
			return BUTTON_LEFT
		}
	} else if mouseEvent, ok := event.OriginalEvent.(*sdl.MouseButtonEvent); ok {
		switch mouseEvent.Button {
		case sdl.BUTTON_LEFT:
			return BUTTON_LEFT
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

func (renderer *Renderer) DrawCircle(x int32, y int32, radius int32, rotation float64, color uint32) {
	gfx.CircleColor(renderer.SDLRenderer, x, y, radius, fromHex(color))
	gfx.LineColor(
		renderer.SDLRenderer,
		x,
		y,
		x+int32(math.Cos(rotation)*float64(radius)),
		y+int32(math.Sin(rotation)*float64(radius)),
		fromHex(color),
	)
}

func (renderer *Renderer) DrawFilledCircle(x int32, y int32, radius int32, color uint32) {
	gfx.FilledCircleColor(renderer.SDLRenderer, x, y, radius, fromHex(color))
}

func (renderer *Renderer) DrawLine(point1 vector.Vec2, point2 vector.Vec2, color uint32) {
	gfx.LineColor(
		renderer.SDLRenderer,
		int32(point1.X),
		int32(point1.Y),
		int32(point2.X),
		int32(point2.Y),
		fromHex(color),
	)
}

func (renderer *Renderer) GetWindowSize() (float64, float64) {
	width, height := renderer.SDLWindow.GetSize()

	return float64(width), float64(height)
}

func (renderer *Renderer) GetMouseCoordinates() (float64, float64) {
	mouseX, mouseY, _ := sdl.GetMouseState()

	return float64(mouseX), float64(mouseY)
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
