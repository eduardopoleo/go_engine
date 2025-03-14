package constants

const (
	PIXEL_PER_METER         float32 = 50 // pix / meter
	FPS                     uint64  = 60
	MILLISECONDS_PER_FRAME  uint64  = 1000 / FPS
	GRAVITY                 float32 = 9.8 * PIXEL_PER_METER // kg⋅pix/s^2
	RESTITUTION_COEFFICIENT float32 = 0.7
)
