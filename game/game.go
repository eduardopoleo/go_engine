package game

import (
	"engine/constants"
	"engine/entities"
	"engine/physics"
	"engine/renderer"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	Renderer            renderer.Renderer
	Particles           [3]entities.Particle
	PushForce           physics.Vec2
	TimeToPreviousFrame uint64
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer
	particle := entities.Particle{
		Position: physics.Vec2{X: 50, Y: 50},
		Radius:   5,
		Color:    0xFFFFFFFF,
		Mass:     2.0,
	}
	game.Particles[0] = particle
	game.PushForce = physics.Vec2{X: 0, Y: 0}
	game.TimeToPreviousFrame = sdl.GetTicks64()
	return game
}

func (game *Game) Input() {
	for event := game.Renderer.PollEvent(); event != nil; event = game.Renderer.PollEvent() {
		switch event.Type {
		case renderer.QUIT:
			println("Quit")
			game.Running = false
			game.Renderer.Destroy()
		case renderer.KEYDOWN:
			if event.Key() == renderer.ESCAPE {
				game.Running = false
				game.Renderer.Destroy()
			} else if event.Key() == renderer.LEFT_ARROW {
				game.PushForce.X = float32(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = float32(50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = float32(50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = float32(-50 * constants.PIXEL_PER_METER)
			}
		}
	}
}

func (game *Game) Update() {
	/*
		It forces the execution to be about 60 frames / second
	*/
	timeElapsed := sdl.GetTicks64() - game.TimeToPreviousFrame

	if constants.MILLISECONDS_PER_FRAME > timeElapsed {
		sdl.Delay(uint32(constants.MILLISECONDS_PER_FRAME - timeElapsed))
	}
	game.TimeToPreviousFrame = sdl.GetTicks64()

	particle := &game.Particles[0]

	particle.SumForces.Add(game.PushForce)
	/*
		Effectively we want the DT to always be 16 miliseconds per frame.
		The delay above will ensure this is the case when the rendering is too fast
		IF the rendering is too slow we still want to have a constant update per frame
		so we force the dt to be 16 miliseconds per frame (~ FPS 60)
	*/
	particle.Integrate(float32(constants.MILLISECONDS_PER_FRAME / 100))
}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	particle := &game.Particles[0]
	particle.Render(&game.Renderer)

	game.Renderer.Render()
	sdl.Delay(16)
}
