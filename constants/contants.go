package constants

const (
	PIXEL_PER_METER         float64 = 50 // pix / meter
	FPS                     uint64  = 60
	MILLISECONDS_PER_FRAME  uint64  = 1000 / FPS
	GRAVITY                 float64 = 9.8 * PIXEL_PER_METER // kg⋅pix/s^2
	RESTITUTION_COEFFICIENT float64 = 0.9
	SPRING_REST_LENGTH      float64 = 15
	SPRING_CONSTANT         float64 = 300
	SPRING_SIZE             uint64  = 10
)
