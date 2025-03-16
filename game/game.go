package game

import (
	"engine/constants"
	"engine/entities"
	"engine/physics"
	"engine/renderer"
	"engine/vector"

	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	Running             bool
	Renderer            renderer.Renderer
	Particles           []entities.Particle
	SpringAnchor        entities.Particle
	SpringParticles     []entities.Particle
	PushForce           vector.Vec2
	TimeToPreviousFrame uint64
}

func NewGame(name string, width int32, height int32) Game {
	game := Game{Running: true}
	renderer := renderer.NewRenderer(name, width, height)
	game.Renderer = renderer
	game.PushForce = vector.Vec2{}
	game.TimeToPreviousFrame = sdl.GetTicks64()
	game.setupSpring(width)

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
				game.PushForce.X = float64(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.RIGHT_ARROW {
				game.PushForce.X = float64(50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.UP_ARROW {
				game.PushForce.Y = float64(-50 * constants.PIXEL_PER_METER)
			} else if event.Key() == renderer.DOWN_ARROW {
				game.PushForce.Y = float64(50 * constants.PIXEL_PER_METER)
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
				mouseX, mouseY := game.Renderer.GetMouseCoordinates()
				particle := entities.Particle{
					Position: vector.Vec2{X: mouseX, Y: mouseY},
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

	currentTime := sdl.GetTicks64()
	deltaTime := float64(currentTime-game.TimeToPreviousFrame) / 1000.0

	// Cap the maximum delta time to prevent instability
	if deltaTime > 0.016 {
		deltaTime = 0.016
	}

	game.TimeToPreviousFrame = sdl.GetTicks64()

	windowWidth, windowHeight := game.Renderer.GetWindowSize()

	// Update other particles
	for i := range game.Particles {
		particle := &game.Particles[i]
		weight := physics.NewWeightForce(particle.Mass)
		particle.SumForces = particle.SumForces.Add(weight)
		particle.SumForces = particle.SumForces.Add(game.PushForce)
	}

	for i := 0; i < len(game.SpringParticles); i++ {
		particle := &game.SpringParticles[i]
		weight := physics.NewWeightForce(particle.Mass)
		particle.SumForces = particle.SumForces.Add(weight)
		particle.SumForces = particle.SumForces.Add(game.PushForce)
	}

	if len(game.SpringParticles) > 0 {
		springForce := physics.NewSpringForce(&game.SpringParticles[0], &game.SpringAnchor)
		game.SpringParticles[0].SumForces = game.SpringParticles[0].SumForces.Add(springForce)
	}

	// Apply spring forces between particles
	for i := 1; i < len(game.SpringParticles); i++ {
		previousParticle := &game.SpringParticles[i-1]
		currentParticle := &game.SpringParticles[i]

		springForce := physics.NewSpringForce(currentParticle, previousParticle)
		currentParticle.SumForces = currentParticle.SumForces.Add(springForce)
		previousParticle.SumForces = previousParticle.SumForces.Add(springForce.Multiply(-1))
	}

	// Integrate last all particles not in between
	// use float64
	// use damping factor to increase stability
	for i := range game.Particles {
		particle := &game.Particles[i]
		particle.Integrate(deltaTime)
		bounce(particle, game, windowWidth, windowHeight, deltaTime)
	}

	for i := 0; i < len(game.SpringParticles); i++ {
		particle := &game.SpringParticles[i]
		particle.Integrate(deltaTime)
		bounce(particle, game, windowWidth, windowHeight, deltaTime)
	}

}

func (game *Game) Draw() {
	game.Renderer.ClearScreen()

	for i := range game.Particles {
		particle := &game.Particles[i]
		particle.Render(&game.Renderer)
	}

	for i := range game.SpringParticles {
		particle := &game.SpringParticles[i]

		particle.Render(&game.Renderer)
		var anchor entities.Particle
		if i == 0 {
			anchor = game.SpringAnchor
		} else {
			anchor = game.SpringParticles[i-1]
		}
		game.Renderer.DrawLine(
			anchor.Position,
			particle.Position,
			renderer.WHITE,
		)
	}

	game.SpringAnchor.Render(&game.Renderer)

	game.Renderer.Render()
}

// private

func (game *Game) setupSpring(windowWidth int32) {
	game.SpringAnchor = entities.Particle{
		Position: vector.Vec2{X: float64(windowWidth) / 2, Y: 0},
		Radius:   5,
		Color:    renderer.WHITE,
		Mass:     2.0,
	}

	var particle entities.Particle
	for i := 0; i < int(constants.SPRING_SIZE); i++ {
		particle = entities.Particle{
			Position: vector.Vec2{
				X: float64(windowWidth) / 2,
				Y: game.SpringAnchor.Position.Y + (constants.SPRING_REST_LENGTH * float64(i+1)),
			},
			Radius: 5,
			Color:  renderer.WHITE,
			Mass:   2.0,
		}
		game.SpringParticles = append(game.SpringParticles, particle)
	}
}

func bounce(particle *entities.Particle, game *Game, windowWidth float64, windowHeight float64, deltaTime float64) {
	if (particle.Position.X - float64(particle.Radius)) <= 0 {
		particle.Velocity.X = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.X
		particle.Position.X = float64(particle.Radius)
	} else if (particle.Position.X + float64(particle.Radius)) >= windowWidth {
		particle.Velocity.X = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.X
		particle.Position.X = windowWidth - float64(particle.Radius)
	} else if (particle.Position.Y - float64(particle.Radius)) <= 0 {
		particle.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.Y
		particle.Position.Y = float64(particle.Radius)
	} else if (particle.Position.Y + float64(particle.Radius)) >= windowHeight {
		particle.Velocity.Y = -constants.RESTITUTION_COEFFICIENT * particle.Velocity.Y
		particle.Position.Y = windowHeight - float64(particle.Radius)
	}
}
