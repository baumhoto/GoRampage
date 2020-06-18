package main

import (
	"github.com/gonutz/prototype/draw"
)

// Renderer renders to the window
type Renderer struct {
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

	// Draw line of sight
	ray := Ray{world.player.position, world.player.direction}
	lineEnd := world.worldmap.hitTest(ray)
	window.DrawLine(int(world.player.position.x*scale), int(world.player.position.y*scale), int(lineEnd.x*scale), int(lineEnd.y*scale), draw.Green)
}
