package common

import "image/color"

const (
	TIMESTEP         = 1.0 / 60.0
	MAXIMUM_TIMESTEP = 1.0 / 20.0
	WORLD_TIMESTEP   = 1.0 / 120.0
	SCREEN_WIDTH     = 320
	SCREEN_HEIGHT    = 240
	SCREEN_SCALE     = 2
)

var (
	Clear = color.RGBA{0, 0, 0, 0}
	Black = color.RGBA{0, 0, 0, 255}
	White = color.RGBA{255, 255, 255, 255}
	Gray  = color.RGBA{192, 192, 192, 255}
	Red   = color.RGBA{255, 0, 0, 255}
	Green = color.RGBA{0, 255, 0, 255}
	Blue  = color.RGBA{0, 0, 255, 255}
)
