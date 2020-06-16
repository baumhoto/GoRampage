package main

import (
	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
	window draw.Window
}

// draw renders the world into the window
func (r Renderer) draw(world World, window draw.Window) {
	_, height := window.Size()
	scale := float64(height) / world.size.y

	rect := world.player.rect()
	rect.min.Multiply(scale)
	rect.max.Multiply(scale)
	window.FillRect(int(world.player.position.x*scale),
		int(world.player.position.y*scale),
		int(rect.max.x-rect.min.y),
		int(rect.max.y-rect.min.y), draw.Blue)
}
