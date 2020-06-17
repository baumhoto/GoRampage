package main

import (
	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
	frameBuffer FrameBuffer
}

// NewRenderer creates a new instance of a Renderer
func NewRenderer(width int, height int) Renderer {
	fb := NewFrameBuffer(width, height, black)
	return Renderer{fb}
}

// draw renders the world into the window
func (r *Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / float64(world.worldmap.Height)

	// Draw map
	for y := 0; y < world.worldmap.Height; y++ {
		for x := 0; x < world.worldmap.Width; x++ {
			if world.worldmap.GetTile(x, y).isWall() {
				window.FillRect(x*int(scale), y*int(scale), 1*int(scale), 1*int(scale), draw.White)
			}
		}
	}

	// Draw player
	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	window.FillRect(int(rect.min.x), int(rect.min.y), int(rect.max.x-rect.min.x), int(rect.max.y-rect.min.y), draw.Blue)
}
