package consts

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
	CLEAR color.RGBA = color.RGBA{}
	BLACK            = color.RGBA{A: 255}
	WHITE            = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	GRAY             = color.RGBA{R: 192, G: 192, B: 192, A: 255}
	RED              = color.RGBA{R: 255, A: 255}
	GREEN            = color.RGBA{G: 255, A: 255}
	BLUE             = color.RGBA{B: 255, A: 255}
)
