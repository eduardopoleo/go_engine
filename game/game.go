package game

import (
	"engine/constants"
	"engine/entities"
	"engine/physics"
	"engine/renderer"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	Renderer            renderer.Renderer
	Particles           []entities.Particle
	PushForce           physics.Vec2
	TimeToPreviousFrame uint64
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer
	game.PushForce = physics.Vec2{}
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
				game.PushForce.Y = float32(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = float32(50 * constants.PIXEL_PER_METER)
			}
		case renderer.KEYUP:
			if event.Key() == renderer.LEFT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = 0
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = 0
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = 0
			}
		case renderer.MOUSE_UP_EVENT:
			if event.Key() == renderer.BUTTON_LEFT {
				fmt.Printf("Got here no?")
				mouseX, mouseY := game.Renderer.GetMouseCoordinates()
				particle := entities.Particle{
					Position: physics.Vec2{X: mouseX, Y: mouseY},
					Radius:   5,
					Color:    0xFFFFFFFF,
					Mass:     2.0,
				}
				game.Particles = append(game.Particles, particle)
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

	windowWidth, windowHeight := game.Renderer.GetWindowSize()

	for i := range game.Particles {
		particle := &game.Particles[i] // Get a reference to the particle

		gravity := physics.NewGravityForce(particle.Mass)
		particle.SumForces = particle.SumForces.Add(gravity)
		particle.SumForces = particle.SumForces.Add(game.PushForce)

		particle.Integrate((float32(constants.MILLISECONDS_PER_FRAME) / 1000))

		/*
			Inelastic collision can be simplified with a change in velocity does not need
			to be force based
		*/
		if (particle.Position.X - float32(particle.Radius)) <= 0 {
			particle.Velocity.X = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.X
			particle.Position.X = float32(particle.Radius)
		} else if (particle.Position.X + float32(particle.Radius)) >= windowWidth {
			particle.Velocity.X = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.X
			particle.Position.X = windowWidth - float32(particle.Radius)
		} else if (particle.Position.Y - float32(particle.Radius)) <= 0 {
			particle.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.Y
			particle.Position.Y = float32(particle.Radius)
		} else if (particle.Position.Y + float32(particle.Radius)) >= windowHeight {
			particle.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.Y
			particle.Position.Y = windowHeight - float32(particle.Radius)
		}
	}

}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	for i := range game.Particles {
		particle := &game.Particles[i]
		particle.Render(&game.Renderer)
	}

	game.Renderer.Render()
	sdl.Delay(16)
}
